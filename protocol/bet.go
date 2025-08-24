package protocol

import (
	"bytes"
	"encoding/binary"
)

type Bet struct {
	Agency   uint8
	Name     string
	Lastname string
	// 8 bytes
	Document uint64
	// "YYYY-MM-DD"
	Birthdate string
	// 2 bytes
	Number uint16
}

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

	// campo de valor fijo que ocupa exactamente 1 byte -> writebyte
	err := payload.WriteByte(bp.Agency)
	if err != nil {
		return nil, err
	}

	// campo de valor variable, mando longitud+data
	nameBytes := []byte(bp.Name)
	err = binary.Write(&payload, binary.BigEndian, uint16(len(nameBytes)))
	if err != nil {
		return nil, err
	}
	payload.WriteString(bp.Name)

	// campo de valor variable, mando longitud+data
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
