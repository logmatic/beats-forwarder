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
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
)

var Registry map[string]output.Output

func init() {

	Registry = make(map[string]output.Output)
	Registry["syslog"] = &output.SyslogClient{}
	Registry["logmatic"] = &output.LogmaticClient{}
	Registry["udp_tcp"] = &output.SocketClient{}

}

func Run(config *cfg.Config) error {

	outputType := *config.Output.Type

	// start the remote connection
	fmt.Printf("Register '%s' as the remote output\n", outputType)
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
		server.Timeout(time.Duration(*config.Input.Timeout) * time.Second),
		server.TLS(digestTLSConfig))

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

func digestTLSConfig(config *cfg.Config) *tls.Config {

	tlsConfig := &tls.Config{}
	if (config.Bool("input.tls.enable", 0) == true) {

		// load client cert
		cert, err := tls.LoadX509KeyPair(*config.Input.TlsConfig.CertPath, *config.Input.TlsConfig.KeyPath)
		if err != nil {
			return tlsConfig
		}

		// load CA
		caCert, err := ioutil.ReadFile(*config.Input.TlsConfig.CaPath)
		if err != nil {
			return tlsConfig
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// get config
		tlsConfig.Certificates = []tls.Certificate{cert}
		tlsConfig.RootCAs = caCertPool

		tlsConfig.BuildNameToCertificate()
	}

	return Input
}


