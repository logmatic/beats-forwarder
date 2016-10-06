package output

import (
	cfg "github.com/logmatic/beats-forwarder/config"
	"errors"
)

type LogmaticClient struct {
	socket         *SocketClient
	logmaticAPIKey string
}

func (c *LogmaticClient) Init(config *cfg.Config) error {

	if (config.Output.Logmatic.Key == nil || *config.Output.Logmatic.Key == "" ) {
		return errors.New("No Logmatic API Key provided.")
	}
	c.logmaticAPIKey = *config.Output.Logmatic.Key + " "

	if (config.Output.Logmatic.Network == nil || config.Output.Logmatic.Raddr == nil) {
		return errors.New("Logmatic configuration missing.")
	}

	socket := &SocketClient{network: *config.Output.Logmatic.Network, raddr: *config.Output.Logmatic.Raddr}
	socket.config = &SocketConfig{maxBackoff: 30, maxRetries: 10}

	c.socket = socket

	return nil
}

func (c *LogmaticClient) WriteAndRetry(payload []byte) (error) {
	return c.socket.WriteAndRetry(append([]byte(c.logmaticAPIKey), payload...))
}

func (c *LogmaticClient) Connect() (error) {
	return c.socket.Connect()
}

func (c *LogmaticClient) Close() {
	c.socket.Close()
}
