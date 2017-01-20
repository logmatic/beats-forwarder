package output

import (
	cfg "github.com/logmatic/beats-forwarder/config"
	"github.com/hashicorp/go-syslog"
)

type SyslogClient struct {
	writer gsyslog.Syslogger
}

func (c *SyslogClient) Init(config *cfg.Config) error {

	network := config.Output.Syslog.Network

	if (network != nil && (*network == "tcp" || *network == "udp")) {
		c.writer, _ = gsyslog.DialLogger(*network, *config.Output.Syslog.Raddr, gsyslog.LOG_INFO, "LOCAL0", *config.Output.Syslog.Tag)
	} else {
		c.writer, _ = gsyslog.NewLogger(gsyslog.LOG_INFO, "LOCAL0", *config.Output.Syslog.Tag)
	}
	//todo (@gpolaert) handles errors
	return nil

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
