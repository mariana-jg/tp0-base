package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/protocol"
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
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
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

func (c *Client) mustStop() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

func (c *Client) MakeBet(bet *protocol.Bet) bool {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	if c.mustStop() {
		return false
	}

	message, err := bet.ToBytes()
	if err != nil {
		log.Errorf("action: apuesta_serializada | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return false
	}

	c.createClientSocket()

	defer c.conn.Close()

	err = avoidShortWrites(c.conn, message)
	if err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return false
	}
	ack, err := avoidShortReads(c.conn, 1)
	if err != nil {
		log.Errorf("action: read_ack | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return false
	}

	if ack[0] == 1 && len(ack) == 1 {
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			bet.Document,
			bet.Number,
		)
		return true
	} else {
		log.Infof("action: apuesta_enviada | result: fail | dni: %v | numero: %v",
			bet.Document,
			bet.Number,
		)
		return false
	}
}
