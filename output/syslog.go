package output

import (
	cfg "github.com/logmatic/beats-forwarder/config"
	"log/syslog"
)

type SyslogClient struct {
	writer *syslog.Writer
}

func (c *SyslogClient) Init(config *cfg.Config) {

	network := config.Output.Syslog.Network

	if (network != nil && (*network == "tcp" || *network == "udp")) {
		c.writer, _ = syslog.Dial(*network, *config.Output.Syslog.Raddr, syslog.LOG_INFO, *config.Output.Syslog.Tag)
	} else {
		c.writer, _ = syslog.New(syslog.LOG_INFO, *config.Output.Syslog.Tag)
	}

}

func (c *SyslogClient) WriteAndRetry(payload []byte) (error) {
	_, err := c.writer.Write(payload)
	return err
}

func (c *SyslogClient) Connect() (error) {
	return nil
}

func (c *SyslogClient) Close() {
	c.writer.Close()
}
