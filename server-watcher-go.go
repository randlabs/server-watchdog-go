/*
Golang library for interacting with the Server monitoring tool.

Source code and other details for the project are available at GitHub:

	https://github.com/randlabs/server-watchdog-go

More usage please see README.md and tests.
*/

package server_watchdog_go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

//------------------------------------------------------------------------------

type ClientOptions struct {
	Host           string `json:"host"`
	Port           uint16 `json:"port"`
	UseSsl         bool `json:"useSsl,omitempty"`
	ApiKey         string `json:"apiKey"`
	DefaultChannel string `json:"defaultChannel"`
	TimeoutMs      uint32 `json:"timeout,omitempty"`
}

type ServerWatcherClient struct {
	url            string
	apiKey         string
	defaultChannel string
	timeoutMs      uint32
}

type processWatchJSON struct {
	Pid      int `json:"pid"`
	Name     string `json:"name"`
	Channel  string `json:"channel"`
	Severity string `json:"severity"`
}

type processUnwatchJSON struct {
	Pid     int `json:"pid"`
	Channel string `json:"channel"`
}

type notifyJSON struct {
	Message  string `json:"message"`
	Channel  string `json:"channel"`
	Severity string `json:"severity"`
}

//------------------------------------------------------------------------------

// Create : Creates a new Server Watcher client module
func Create(options ClientOptions) (*ServerWatcherClient, error) {
	c := &ServerWatcherClient{}

	//check host
	if len(options.Host) == 0 {
		return nil, errors.New("invalid host")
	}

	//check port
	if options.Port < 0 || options.Port > 65535 {
		return nil, errors.New("invalid port")
	}

	c.url = "http"
	if options.UseSsl {
		c.url += "s"
	}
	c.url += "://" + options.Host + ":" + strconv.Itoa(int( options.Port))

	//check api key
	if len(options.ApiKey) == 0 {
		return nil, errors.New("invalid API key")
	}
	c.apiKey = options.ApiKey

	//check default channel
	if len(options.DefaultChannel) == 0 {
		return nil, errors.New("invalid default channel")
	}
	c.defaultChannel = options.DefaultChannel

	//timeout
	if options.TimeoutMs == 0 {
		c.timeoutMs = 30000
	} else {
		c.timeoutMs = options.TimeoutMs
	}

	return c, nil
}

// GetDefaultChannel : Retrieves the default channel to use
func (s *ServerWatcherClient) GetDefaultChannel() string {
	return s.defaultChannel
}

// Error : Sends an error message to the server
func (s *ServerWatcherClient) Error(message string, channel string) error {
	return s.notify(message, channel, "error")
}

// Warn : Sends a warning message to the server
func (s *ServerWatcherClient) Warn(message string, channel string) error {
	return s.notify(message, channel, "warn")
}

// Info : Sends an information message to the server
func (s *ServerWatcherClient) Info(message string, channel string) error {
	return s.notify(message, channel, "info")
}

// ProcessWatch : Instructs the server to monitor the specified process
func (s *ServerWatcherClient) ProcessWatch(pid int, name, severity, channel string) error {
	var payload processWatchJSON
	var payloadBytes []byte
	var err error

	payload.Pid = pid
	if pid == 0 {
		payload.Pid = os.Getpid()
	} else if pid < 1 {
		return errors.New("invalid process id")
	}

	if len(name) == 0 {
		payload.Name, err = os.Executable()
		if err != nil {
			return err
		}
	} else {
		payload.Name = name
	}

	payload.Severity, err = s.validateSeverity(severity)
	if err != nil {
		return err
	}

	payload.Channel = s.validateChannel(channel)

	payloadBytes, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = s.sendRequest("process/watch", payloadBytes,  false)
	return err
}

// ProcessUnwatch : Instructs the server to stop monitoring the specified process
func (s *ServerWatcherClient) ProcessUnwatch(pid int, channel string) error {
	var payload processUnwatchJSON
	var payloadBytes []byte
	var err error

	payload.Pid = pid
	if pid == 0 {
		payload.Pid = os.Getpid()
	} else if pid < 1 {
		return errors.New("invalid process id")
	}

	payload.Channel = s.validateChannel(channel)

	payloadBytes, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = s.sendRequest("process/unwatch", payloadBytes,  false)
	return err
}

func (s *ServerWatcherClient) notify(message string, channel string, severity string) error {
	var payload notifyJSON
	var payloadBytes []byte
	var err error

	if len(message) == 0 {
		return nil
	}
	payload.Message = message

	payload.Channel = s.validateChannel(channel)

	payload.Severity = severity

	payloadBytes, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = s.sendRequest("notify", payloadBytes,  false)
	return err
}

func (s *ServerWatcherClient) sendRequest(api string, payloadBytes []byte, hasResponse bool) ([]byte, error) {
	var req *http.Request
	var resp *http.Response
	var bodyBytes []byte
	var err error

	client := &http.Client{}

	req, err = http.NewRequest("POST", s.url + "/" + api, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", s.apiKey)

	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(s.timeoutMs) * time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err = client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != 200 {
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return []byte{}, err
		}

		msg := string(bodyBytes)
		if len(msg) == 0 {
			msg = "Unsuccessful response from node."
		}
		return []byte{}, fmt.Errorf("%v [Status: %v]", msg, resp.StatusCode)
	}

	if !hasResponse {
		_ = resp.Body.Close()
		return []byte{}, nil
	}

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	return bodyBytes, nil
}

func (s *ServerWatcherClient) validateChannel(channel string) string {
	if len(channel) == 0 {
		return s.defaultChannel
	}
	return channel
}

func  (s *ServerWatcherClient) validateSeverity(severity string) (string, error) {
	if len(severity) == 0 {
		return "error", nil
	}
	if severity != "error" && severity != "warn" && severity != "info" {
		return "", errors.New("invalid severity")
	}
	return severity, nil
}
