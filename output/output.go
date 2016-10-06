package output

import cfg "github.com/logmatic/beats-forwarder/config"

// Interface that all outputs should be implemented
type Output  interface {
	Init(*cfg.Config) error
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
	err := remote.Init(config)
	if err != nil {
		return nil, err
	}
	err = remote.Connect()
	return remote, err

}
