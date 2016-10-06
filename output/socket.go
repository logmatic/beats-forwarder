package output

import (
	"fmt"
	"time"
	"math"
	"sync"
	"net"

	"crypto/tls"
	"io/ioutil"
	"crypto/x509"

	"github.com/Sirupsen/logrus"
	cfg "github.com/logmatic/beats-forwarder/config"
)

type Connection interface {
	Write([] byte) (int, error)
	Close() (error)
}

type SocketClient struct {
	conn      Connection
	config    *SocketConfig
	tlsConfig *tls.Config

	network   string
	raddr     string
}

type SocketConfig struct {
	maxRetries int
	maxBackoff int
}

func (c *SocketClient) Init(config *cfg.Config) error {

	c.network = *config.Output.UDPTCP.Network
	c.raddr = *config.Output.UDPTCP.Raddr
	c.config = &SocketConfig{maxBackoff: 30, maxRetries: 10}

	if (*config.Output.UDPTCP.TlsConfig.Enable == true) {

		c.network = "tcp"
		// load client cert
		cert, err := tls.LoadX509KeyPair(*config.Output.UDPTCP.TlsConfig.CertPath, *config.Output.UDPTCP.TlsConfig.KeyPath)
		if err != nil {
			return err
		}

		// load CA
		caCert, err := ioutil.ReadFile(*config.Output.UDPTCP.TlsConfig.CaPath)
		if err != nil {
			return err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// get config
		c.tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		c.tlsConfig.BuildNameToCertificate()
	} else {

		c.tlsConfig = &tls.Config{}
	}

	return nil
}

func (socket *SocketClient) Connect() (error) {

	logrus.Infof("Connection to %s (%s)\n", socket.raddr, socket.network)

	var conn Connection
	var err error
	if (socket.network == "tls" || socket.network == "ssl") {
		conn, err = tls.Dial("tcp", socket.raddr, socket.tlsConfig)

	} else {
		conn, err = net.Dial(socket.network, socket.raddr)

	}
	if (err != nil) {
		return err
	}
	socket.conn = conn;
	return nil
}

func (socket *SocketClient) reconnect() (error) {
	socket.Close()
	return socket.Connect()
}

func (socket *SocketClient) WriteAndRetry(payload []byte) (error) {

	for i := 0; i < socket.config.maxRetries; i++ {

		// backoff mechanism
		backoffTime := int(math.Min(math.Pow(float64(i), 2), float64(socket.config.maxBackoff)));
		if (backoffTime > 0) {
			logrus.Warnf("Making a new attempt in %d seconds (%d/%d)", backoffTime, i, socket.config.maxRetries);
		}

		time.Sleep(time.Duration(backoffTime) * time.Second)

		mutex := sync.Mutex{}
		mutex.Lock()

		if (socket.conn == nil) {

			// reconnect
			err := socket.reconnect();
			if err != nil {
				logrus.Errorf("Unable to connect, error: %v", err)
				socket.Close()
				continue
			}
		}
		mutex.Unlock();


		err := socket.writeOnce(payload)
		if err != nil {
			logrus.Errorf("Unable to write, error: %v", err)
			continue
		}

		return nil

	}

	return fmt.Errorf("Failed to connect to %s (%s)", socket.raddr, socket.network)
}

func (socket *SocketClient) writeOnce(payload []byte) (error) {

	_, err := socket.conn.Write(payload)
	if err != nil {
		socket.Close()
		return err

	}
	return nil

}

func (socket *SocketClient) Close() {
	socket.conn.Close()
	socket.conn = nil
}
