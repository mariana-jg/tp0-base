package common

import (
	"fmt"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/protocol"
)

func parseBetLine(line []string, agency int) (*protocol.Bet, error) {

	if len(line) != 5 {
		return nil, fmt.Errorf("incomplete data")
	}

	document, err := strconv.ParseUint(line[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse document: %w", err)
	}

	number, err := strconv.ParseUint(line[4], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("parse number: %w", err)
	}

	return protocol.NewBet(
		uint8(agency),
		line[0],
		line[1],
		document,
		line[3],
		uint16(number),
	), nil
}
