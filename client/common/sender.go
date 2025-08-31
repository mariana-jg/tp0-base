package common

import "github.com/7574-sistemas-distribuidos/docker-compose-init/protocol"

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
