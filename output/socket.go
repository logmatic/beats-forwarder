package output

import (
	"fmt"
	"os"
	"time"
	"math"
	"sync"
	"errors"
	"net"
	"crypto/tls"
)

type Connection interface {
	Write([] byte) (int, error)
	Close() (error)
}

type Socket struct {
	conn    Connection
	config  SocketConfig

	network string
	raddr   string
}

type SocketConfig struct {
	maxRetries int
	maxBackoff int
}

func NewSocket(network string, raddr string, config SocketConfig) *Socket {

	socket := Socket{network: network, raddr: raddr}
	socket.config = config

	return &socket

}

func (socket *Socket) Connect() (error) {

	fmt.Fprintf(os.Stderr, "Connection to %s (%s)\n", socket.raddr, socket.network)

	var conn Connection
	var err error
	if (socket.network == "tls" || socket.network == "ssl") {
		conn, err = tls.Dial("tcp", socket.raddr, &tls.Config{})

	} else {
		conn, err = net.Dial(socket.network, socket.raddr)

	}
	if (err != nil) {
		return err
	}
	socket.conn = conn;
	return nil
}

func (socket *Socket) reconnect() (error) {
	socket.Close()
	return socket.Connect()
}

func (socket *Socket) WriteAndRetry(payload []byte) (error) {

	for i := 0; i < socket.config.maxRetries; i++ {

		// backoff mechanism
		backoffTime := int(math.Min(math.Pow(float64(i), 2), float64(socket.config.maxBackoff)));
		if (backoffTime > 0) {
			fmt.Fprintf(os.Stderr, "Making a new attempt in %d seconds (%d/%d)\n", backoffTime, i, socket.config.maxRetries);
		}

		time.Sleep(time.Duration(backoffTime) * time.Second)

		mutex := sync.Mutex{}
		mutex.Lock()

		if (socket.conn == nil) {

			// reconnect
			err := socket.reconnect();
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connect, error: %v\n", err)
				socket.Close()
				continue
			}
		}
		mutex.Unlock();


		err := socket.writeOnce(payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write, error: %v\n", err)
			continue
		}

		return nil

	}

	return errors.New(fmt.Sprintf("Failed to connect to %s (%s)", socket.raddr, socket.network))
}

func (socket *Socket) writeOnce(payload []byte) (error) {

	_, err := socket.conn.Write(payload)
	if err != nil {
		socket.Close()
		return err

	}
	return nil

}

func (socket *Socket) Close() {
	socket.conn.Close()
	socket.conn = nil
}
