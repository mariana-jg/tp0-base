package common

import (
	"encoding/binary"
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

func (c *Client) sendBatch(batch []*protocol.Bet) bool {
	message := protocol.BatchToBytes(batch)

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

func (c *Client) sendDoneAndReadWinners(agency int) ([]uint64, bool) {
	message := protocol.DoneToBytes(uint8(agency))

	if err := mustWriteAll(c.conn, message); err != nil {
		log.Errorf("action: done_enviado | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return nil, false
	}

	ack, err := mustReadAll(c.conn, 1)
	if err == nil && len(ack) == 1 && ack[0] == 1 {
		log.Infof("action: ack_from_server | result: success | agency: %v", agency)
	}

	countB, err := mustReadAll(c.conn, 2)
	if err != nil {
		log.Errorf("action: read_winners_count | result: fail | err: %v", err)
		return nil, false
	}
	count := binary.BigEndian.Uint16(countB)

	winners := make([]uint64, 0, count)
	for i := 0; i < int(count); i++ {
		dniB, err := mustReadAll(c.conn, 8)
		if err != nil {
			log.Errorf("action: read_winner_dni | result: fail | err: %v", err)
			return nil, false
		}
		winners = append(winners, binary.BigEndian.Uint64(dniB))
	}

	return winners, true
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
		if err != nil {
			log.Errorf("action: read_bet | result: fail | client_id: %v | error: %v", c.config.ID, err)
			allSucceeded = false
			break
		}
		if len(line) != 5 {
			log.Errorf("action: read_bet | result: fail | client_id: %v | error: Incomplete data", c.config.ID)
			continue
		}

		document, err := strconv.ParseUint(line[2], 10, 64)
		if err != nil {
			log.Errorf("action: parse_document | result: fail | value: %q | error: %v", line[2], err)
			continue
		}
		number, err := strconv.ParseUint(line[4], 10, 16)
		if err != nil {
			log.Errorf("action: parse_number | result: fail | value: %q | error: %v", line[4], err)
			continue
		}

		bet := protocol.NewBet(
			uint8(agency),
			line[0],
			line[1],
			document,
			line[3],
			uint16(number),
		)
		batch = append(batch, bet)

		if len(batch) == c.config.BatchSize {
			if c.sendBatch(batch) {
				log.Infof("action: apuesta_enviada | result: success | batch_size: %v", len(batch))
			} else {
				log.Infof("action: apuesta_enviada | result: fail | batch_size: %v", len(batch))
				allSucceeded = false
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if c.sendBatch(batch) {
			log.Infof("action: apuesta_enviada | result: success | batch_size: %v", len(batch))
		} else {
			log.Infof("action: apuesta_enviada | result: fail | batch_size: %v", len(batch))
			allSucceeded = false
		}
	}

	winners, ok := c.sendDoneAndReadWinners(agency)
	if ok {
		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners))
	} else {
		log.Infof("action: consulta_ganadores | result: fail | cant_ganadores: 0")
	}

	return allSucceeded
}
