package common

import (
	"encoding/csv"
	"io"
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
		//log.Info("action: exit | result: success")
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
			return nil
		}
		time.Sleep(time.Duration(i*500) * time.Millisecond)
	}
	log.Criticalf(
		"action: connect | result: fail | client_id: %v | error: %v",
		c.config.ID,
		err,
	)
	return err
}
func (c *Client) mustStop() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

func (c *Client) MakeBet(path string) bool {
	agency, _ := strconv.Atoi(c.config.ID)
	if c.mustStop() {
		return false
	}
	if err := c.createClientSocket(); err != nil {
		log.Errorf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return false
	}
	defer c.conn.Close()

	file, err := os.Open(path)
	if err != nil {
		log.Errorf("action: open_bets_file | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return false
	}
	defer file.Close()

	reader := csv.NewReader(file)
	batch := make([]*protocol.Bet, 0, c.config.BatchSize)
	allSucceeded := true

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		bet, err := parseBetLine(line, agency)
		if err != nil {
			log.Errorf("action: parse_bet_line | result: fail | error: %s", err)
			allSucceeded = false
			continue
		}
		batch = append(batch, bet)

		if len(batch) == c.config.BatchSize {
			ok, err := c.sendBatch(batch)
			if err == ErrServerShutdown {
				_ = c.conn.Close()
				return false
			}
			if err != nil || !ok {
				log.Infof("action: apuesta_enviada | result: fail | batch_size: %v", len(batch))
				allSucceeded = false
				break // cortar el loop si falló envío/ACK
			}
			log.Infof("action: apuesta_enviada | result: success | batch_size: %v", len(batch))
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		ok, err := c.sendBatch(batch)
		if err == ErrServerShutdown {
			_ = c.conn.Close()
			return false
		}
		if err != nil || !ok {
			log.Infof("action: apuesta_enviada | result: fail | batch_size: %v", len(batch))
			allSucceeded = false
		} else {
			log.Infof("action: apuesta_enviada | result: success | batch_size: %v", len(batch))
		}
	}

	winners, err := c.sendDoneAndReadWinners(agency)
	if err == ErrServerShutdown {
		_ = c.conn.Close()
		return false
	}
	if err != nil {
		log.Infof("action: consulta_ganadores | result: fail | cant_ganadores: 0")
	} else {
		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners))
	}
	return allSucceeded
}
