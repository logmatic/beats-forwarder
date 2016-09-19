package output

import (
	cfg "github.com/logmatic/beats-forwarder/config"
)

type LogmaticClient struct {
	socket         *Socket
	logmaticAPIKey string
}

func NewLogmaticClient() *LogmaticClient {
	return &LogmaticClient{}
}

func (c *LogmaticClient) Init(config *cfg.Config) {

	socketConfig := SocketConfig{maxBackoff: 30, maxRetries: 10}

	//todo (gpolaert) handle config errors
	c.logmaticAPIKey = *config.Output.Logmatic.Key + " "
	c.socket = NewSocket(*config.Output.Logmatic.Network, *config.Output.Logmatic.Raddr, socketConfig)
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
