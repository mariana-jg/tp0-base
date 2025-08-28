package common

import (
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/protocol"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

const MAX_TRIES = 5

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BatchSize     int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	done   chan struct{}
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
	var err error

	for i := 1; i <= MAX_TRIES; i++ {
		conn, err := net.Dial("tcp", c.config.ServerAddress)
		if err == nil {
			c.conn = conn
			log.Infof(
				"action: connect | result: success | attempt: %d/%d | client_id: %v",
				i,
				MAX_TRIES,
				c.config.ID,
			)
			return err
		}
		time.Sleep(time.Duration(i*500) * time.Millisecond)
	}
	log.Criticalf(
		"action: connect | result: fail | client_id: %v | error: %v",
		c.config.ID,
		err,
	)
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

func (c *Client) sendBatch(batch []*protocol.Bet) bool {
	message := protocol.BatchToBytes(batch)

	if err := c.createClientSocket(); err != nil {
		log.Errorf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return false
	}

	defer c.conn.Close()

	if err := mustWriteAll(c.conn, message); err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return false
	}

	ack, err := mustReadAll(c.conn, 1)
	if err != nil {
		log.Errorf("action: read_ack | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return false
	}

	return len(ack) == 1 && ack[0] == 1
}

func (c *Client) MakeBet(path string) bool {
	agency, _ := strconv.Atoi(c.config.ID)
	bets, err := c.ReadBetsFromFile(path, agency)
	if err != nil {
		return false
	}
	batches := c.CreateBatch(bets)

	if c.mustStop() {
		return false
	}

	allSucceeded := true

	for _, batch := range batches {
		if c.sendBatch(batch) {
			log.Infof("action: apuesta_enviada | result: success | batch_size: %v", len(batch))
		} else {
			log.Infof("action: apuesta_enviada | result: fail | batch_size: %v", len(batch))
			allSucceeded = false
		}
	}

	log.Infof("action: exit | result: success")
	return allSucceeded
}
