package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	// flag para salir del loop principal
	done chan struct{}
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		done:   make(chan struct{}),
	}
	// canal donde se envia cuando el proceso recibe un SIGTERM
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	// rutina que se desbloquea cuando llega un SIGTERM, espera esa senial (corre en paralelo)
	go func() {
		<-signalChan
		close(client.done)
		log.Info("action: exit | result: success")
	}()
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

func (c *Client) debeTerminar() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

func (c *Client) esperaShutdown(d time.Duration) bool {
	// Chequeo entre esperas
	select {
	//Canal se cierra por un SIGTERM, interrumpo
	case <-c.done:
		return true
	// si no llego senial, espero un tiempo y despues continuo
	// se saca el sleep porque sino no se enteraria del shutdown hasta que se despierte
	case <-time.After(d):
		return false
	}
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Chequeo antes de conectarme
		if c.debeTerminar() {
			return
		}

		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		// TODO: Modify the send to avoid short-write
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Chequeo entre espera
		if c.esperaShutdown(c.config.LoopPeriod) {
			return
		}

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
