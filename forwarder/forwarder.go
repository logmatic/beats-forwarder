package forwarder

import (
	"fmt"
	"os"
	"os/signal"
	"encoding/json"
	"bytes"
	"time"

	"github.com/elastic/go-lumber/server"
	"github.com/Sirupsen/logrus"
	"github.com/logmatic/beats-forwarder/output"
	cfg "github.com/logmatic/beats-forwarder/config"

	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
	"errors"
)

var Registry map[string]output.Output

func init() {

	Registry = make(map[string]output.Output)
	Registry["syslog"] = &output.SyslogClient{}
	Registry["logmatic"] = &output.LogmaticClient{}
	Registry["udp_tcp"] = &output.SocketClient{}

}

func Run(config *cfg.Config) error {

	var out output.Output
	if config.Output.Type != nil {
		out = Registry[*config.Output.Type]
	} else {
		return errors.New("Config error: output.type not set")
	}

	if out == nil {
		return errors.New("Config error: output.type '" + *config.Output.Type + "' is unknown")
	}

	// start the remote connection
	logrus.Infof("Register '%s' as the remote output", *config.Output.Type)
	remote, err := output.Run(out, config)
	if err != nil {
		return err
	}


	// todo (gpolaert) factorize
	input_port := 5044
	input_host := "0.0.0.0"
	input_lj_v1 := false
	input_lj_v2 := true
	input_keepalive := 3
	input_timeout := 30

	if config.Input.Host != nil {
		input_host = *config.Input.Host
	}

	if config.Input.Port != nil {
		input_port = *config.Input.Port
	}

	if config.Input.LJ.V1 != nil {
		input_lj_v1 = *config.Input.LJ.V1
	}

	if config.Input.LJ.V2 != nil {
		input_lj_v2 = *config.Input.LJ.V2
	}

	if config.Input.Keepalive != nil {
		input_keepalive = *config.Input.Keepalive
	}

	if config.Input.Timeout != nil {
		input_timeout = *config.Input.Timeout
	}



	// start the listener
	local, err := server.ListenAndServe(
		fmt.Sprintf("%s:%d", input_host, input_port),
		server.V1(input_lj_v1),
		server.V2(input_lj_v2),
		server.Keepalive(time.Duration(input_keepalive) * time.Second),
		server.Timeout(time.Duration(input_timeout) * time.Second),
		server.TLS(getTLSConfig(config)))

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

		logrus.Debugf("Received batch of %v events", len(batch.Events))
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

func getTLSConfig(config *cfg.Config) *tls.Config {

	if (*config.Input.TlsConfig.Enable == true) {

		tlsConfig := &tls.Config{}
		logrus.Infof("Setting an encrypted communication for the input")

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
		return tlsConfig
	}

	return nil

}


