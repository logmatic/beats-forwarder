package output

import cfg "github.com/logmatic/beats-forwarder/config"

// Interface that all outputs should be implemented
type Output  interface {
	Init(*cfg.Config)
	Connect() error
	WriteAndRetry([]byte) error
	Close()
}

// the vehicle ...
type OutputImplementer struct {
	Output
}

func Run(impl Output, config *cfg.Config) (Output, error) {

	remote := OutputImplementer{impl}
	remote.Init(config)
	err := remote.Connect()
	return remote, err

}
