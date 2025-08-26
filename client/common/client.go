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

func (c *Client) CreateBetsFromCSV(pathBets string, agencia int) ([][]*protocol.Bet, error) {
	file, err := os.Open(pathBets)
	if err != nil {
		log.Errorf("action: open_file | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var allBets []*protocol.Bet

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("action: read_bet | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return nil, err
		}

		if len(line) != 5 {
			log.Errorf("action: read_bet | result: fail | client_id: %v | error: Insufficient data on line", c.config.ID)
			continue
		}
		numeroApostado, _ := strconv.ParseUint(line[4], 10, 64)
		documento, _ := strconv.ParseUint(line[4], 10, 16)
		bet := protocol.NewBet(
			uint8(agencia),
			line[0],
			line[1],
			documento,
			line[3],
			uint16(numeroApostado),
		)

		allBets = append(allBets, bet)
	}

	var betBatches [][]*protocol.Bet
	for i := 0; i < len(allBets); i += c.config.BatchSize {
		end := i + c.config.BatchSize
		if end > len(allBets) {
			end = len(allBets)
		}
		betBatches = append(betBatches, allBets[i:end])
	}

	return betBatches, nil
}

func (c *Client) MakeBet(path string) bool {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed

	agency, _ := strconv.Atoi(c.config.ID)
	batches, err := c.CreateBetsFromCSV(path, agency)
	if err != nil {
		return false
	}

	if c.mustStop() {
		return false
	}

	for _, batch := range batches {
		message := protocol.BatchToBytes(batch)
		if err != nil {
			log.Errorf("action: batch_serializado | result: fail | client_id: %v | error: %v", c.config.ID, err)
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
			log.Infof("action: apuesta_enviada | result: success | batch_size: %v",
				len(batch),
			)
			return true
		} else {
			log.Infof("action: apuesta_enviada | result: fail | batch_size: %v",
				len(batch),
			)
			return false
		}
	}
	log.Infof("action: exit | result: success")
	c.conn.Close()
	return true
}
