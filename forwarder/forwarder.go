package forwarder

import (
	"fmt"
	"os"
	"os/signal"
	"encoding/json"
	"bytes"
	"time"

	"github.com/elastic/go-lumber/server"
	"github.com/logmatic/beats-forwarder/output"
	cfg "github.com/logmatic/beats-forwarder/config"
)


var Registry map[string]output.Output

func init() {

	registry := make(map[string]output.Output)
	registry["syslog"] = &output.SyslogClient{}
	registry["logmatic"] = &output.LogmaticClient{}
	registry["udp_tcp"] = &output.LogmaticClient{}
}

func Run(config *cfg.Config) error {



	outputType := *config.Output.Type

	// start the remote connection
	remote, err := output.Run(Registry[outputType], config)
	if err != nil {
		return err
	}

	// start the listener
	local, err := server.ListenAndServe(
		fmt.Sprintf("%s:%d", *config.Input.Host, *config.Input.Port),
		server.V1(*config.Input.LJ.V1),
		server.V2(*config.Input.LJ.V2),
		server.Keepalive(time.Duration(*config.Input.Keepalive) * time.Second),
		server.Timeout(time.Duration(*config.Input.Timeout) * time.Second))

	if err != nil {
		return err
	}

	// wait until received a stop signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	go func() {
		<-sig

		_ = local.Close()
		os.Exit(0)
	}()


	// main loop
	for batch := range local.ReceiveChan() {

		fmt.Printf("Received batch of %v events\n", len(batch.Events))
		payload := new(bytes.Buffer)

		for _, beat := range batch.Events {
			payload.Reset()
			json.NewEncoder(payload).Encode(beat)
			remote.WriteAndRetry(payload.Bytes())
		}
		batch.ACK()

	}

	return nil

}



