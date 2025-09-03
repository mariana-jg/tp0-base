package common

import (
	"encoding/binary"
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/protocol"
)

const SERVER_SHUTDOWN = 255
const ACK = 1

func (c *Client) sendBatch(batch []*protocol.Bet) (bool, error) {
	message := protocol.BatchToBytes(batch)

	if err := mustWriteAll(c.conn, message); err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return false, err
	}

	ack, err := mustReadAll(c.conn, 1)
	if err != nil {
		log.Errorf("action: read_ack_batch | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return false, err
	}

	if ack[0] == SERVER_SHUTDOWN {
		log.Infof("action: server_shutdown | result: detected | client_id: %v", c.config.ID)
		return false, ErrServerShutdown
	}
	if ack[0] == ACK {
		return true, nil
	}
	return false, fmt.Errorf("unexpected ACK code: %d", ack[0])
}

// sender.go
func (c *Client) sendDoneAndReadWinners(agency int) ([]uint64, error) {
	message := protocol.DoneToBytes(uint8(agency))

	if err := mustWriteAll(c.conn, message); err != nil {
		log.Errorf("action: done_enviado | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return nil, err
	}

	ack, err := mustReadAll(c.conn, 1)
	if err == nil && len(ack) == 1 && ack[0] == ACK {
		log.Infof("action: read_ack_done | result: success | agency: %v", agency)
	} else if err != nil {
		return nil, err
	}

	if ack[0] == SERVER_SHUTDOWN {
		log.Infof("action: server_shutdown | result: detected | client_id: %v", c.config.ID)
		return nil, ErrServerShutdown
	}

	countB, err := mustReadAll(c.conn, 2)
	if err != nil {
		log.Errorf("action: read_winners_count | result: fail | err: %v", err)
		return nil, err
	}
	count := binary.BigEndian.Uint16(countB)

	winners := make([]uint64, 0, count)
	for i := 0; i < int(count); i++ {
		dniB, err := mustReadAll(c.conn, 8)
		if err != nil {
			log.Errorf("action: read_winner_dni | result: fail | err: %v", err)
			return nil, err
		}
		winners = append(winners, binary.BigEndian.Uint64(dniB))
	}
	return winners, nil
}
