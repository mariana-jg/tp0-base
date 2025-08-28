package common

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/protocol"
)

func (c *Client) ReadBetsFromFile(pathBets string, agencia int) ([]*protocol.Bet, error) {
	file, err := os.Open(pathBets)
	if err != nil {
		log.Errorf("action: open_bets_file | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var bets []*protocol.Bet

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
			log.Errorf("action: read_bet | result: fail | client_id: %v | error: Incomplete data", c.config.ID)
			continue
		}

		document, _ := strconv.ParseUint(line[2], 10, 16)
		name := line[0]
		lastname := line[1]
		birthdate := line[3]
		number, _ := strconv.ParseUint(line[4], 10, 64)

		bet := protocol.NewBet(
			uint8(agencia),
			name,
			lastname,
			document,
			birthdate,
			uint16(number),
		)

		bets = append(bets, bet)
	}

	return bets, nil

}

func (c *Client) CreateBatch(bets []*protocol.Bet) [][]*protocol.Bet {
	var betBatches [][]*protocol.Bet

	for i := 0; i < len(bets); i += c.config.BatchSize {
		end := i + c.config.BatchSize
		if end > len(bets) {
			end = len(bets)
		}
		betBatches = append(betBatches, bets[i:end])
	}

	return betBatches
}
