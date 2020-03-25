# server-watchdog-go

Golang library for interacting with the [Server monitoring](https://github.com/randlabs/server-watchdog) tool.

# Installation

To install the library, just import it using the following sentence:
```golang
import "github.com/randlabs/server-watchdog-go"
```

# Usage sample

```golang
import swc "github.com/randlabs/server-watchdog-go"

//...

swdClient, err := swc.Create(swc.ClientOptions{
    Host: "127.0.0.1",
    Port: 3004,
    ApiKey: "set-some-key",
    DefaultChannel: "default",
})
	
err = swdClient.Info("This is a sample information message", "")
if err == nil {
    log.Printf("Success\n")
} else {
    log.Fatalf("Error sending message [%v]\n", err.Error())
}
```

# Documentation

To create the client instance use:

#### `swdClient, err := swc.Create(options)`

##### `options`

###### `options.Host`

Specifies the server host address.

###### `options.Port`

Specifies the server port.

###### `options.UseSsl`

Indicates if the connection to server must use a secure channel.

###### `options.ApiKey`

Sets the key used to access the server. This is intended to be secret. The key must match the string configured in the server configuration.

###### `options.DefaultChannel`

Sets the default channel to use when a channel is not specified on notification methods. Read the server documentation for details about channels.

###### `options.TimeoutMs` (optional)

Establishes the maximum time to use when sending messages to the server in millisecond units. A default value of ten (10) seconds is used if this option is not specified.

## Methods

###### `swdClient.Error(message, channel)`

Sends an error message to the server using the specified or default channel (if channel is empty).

Returns nil or an error.

###### `swdClient.Warn(message, channel)`

Sends a warning message to the server using the specified or default channel (if channel is empty).

Returns nil or an error.

###### `swdClient.Info(message, channel)`

Sends an information message to the server using the specified or default channel (if channel is empty).

Returns nil or an error.

###### `swdClient.ProcessWatch(pid, name, severity[, channel)`

Informs the server to monitor the specified process. `severity` can be `error`, `warn` or `info`. If empty, `error` is used.

If the process is killed or exits with an exit code different than zero, the server will send the proper notification to the specified or default channel (if channel is empty).

Returns nil or an error.

###### `swdClient.ProcessUnwatch(pid, channel)`

Informs the server to stop monitoring the specified process.

Returns nil or an error.

## License

See [LICENSE](LICENSE) file.