package protocol

import (
	"bytes"
	"encoding/binary"
)

type Bet struct {
	Agency    uint8
	Name      string
	Lastname  string
	Document  uint64
	Birthdate string
	Number    uint16
}

const (
	TYPE_BET_BATCH = 1
)

func NewBet(agency uint8, name string, lastname string, document uint64, birthdate string, number uint16) *Bet {
	return &Bet{
		Agency:    agency,
		Name:      name,
		Lastname:  lastname,
		Document:  document,
		Birthdate: birthdate,
		Number:    number,
	}
}

func (bp Bet) ToBytes() ([]byte, error) {

	var payload bytes.Buffer

	err := payload.WriteByte(bp.Agency)
	if err != nil {
		return nil, err
	}

	nameBytes := []byte(bp.Name)
	err = binary.Write(&payload, binary.BigEndian, uint16(len(nameBytes)))
	if err != nil {
		return nil, err
	}
	payload.WriteString(bp.Name)

	lastnameBytes := []byte(bp.Lastname)
	err = binary.Write(&payload, binary.BigEndian, uint16(len(lastnameBytes)))
	if err != nil {
		return nil, err
	}
	payload.WriteString(bp.Lastname)

	err = binary.Write(&payload, binary.BigEndian, bp.Document)
	if err != nil {
		return nil, err
	}

	_, err = payload.Write([]byte(bp.Birthdate))
	if err != nil {
		return nil, err
	}

	err = binary.Write(&payload, binary.BigEndian, bp.Number)
	if err != nil {
		return nil, err
	}

	data := payload.Bytes()
	var frame bytes.Buffer

	err = binary.Write(&frame, binary.BigEndian, uint16(len(data)))
	if err != nil {
		return nil, err
	}

	frame.Write(data)

	return frame.Bytes(), nil
}

func BatchToBytes(batch []*Bet) []byte {
	//[0]: 0x01 TYPE_BET_BATCH
	data := []byte{TYPE_BET_BATCH, byte(len(batch))}

	for _, bet := range batch {
		bet_bytes, _ := bet.ToBytes()
		data = append(data, bet_bytes...)
	}

	return data
}
