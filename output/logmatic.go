package output

import (
	cfg "github.com/logmatic/beats-forwarder/config"
)

type LogmaticClient struct {
	socket         *SocketClient
	logmaticAPIKey string
}


func (c *LogmaticClient) Init(config *cfg.Config) {


	//todo (gpolaert) handle config errors
	c.logmaticAPIKey = *config.Output.Logmatic.Key + " "

	socket := &SocketClient{network: *config.Output.Logmatic.Network, raddr: *config.Output.Logmatic.Raddr}
	socket.config = SocketConfig{maxBackoff: 30, maxRetries: 10}

	c.socket = socket
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
