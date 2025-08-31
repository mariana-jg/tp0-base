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

	// 1) agency (uint8)
	if err := payload.WriteByte(bp.Agency); err != nil {
		return nil, err
	}

	// 2) name_len (uint16) + name (bytes)
	nameBytes := []byte(bp.Name)
	if err := binary.Write(&payload, binary.BigEndian, uint16(len(nameBytes))); err != nil {
		return nil, err
	}
	if _, err := payload.Write(nameBytes); err != nil {
		return nil, err
	}

	// 3) last_name_len (uint16) + last_name (bytes)
	lastnameBytes := []byte(bp.Lastname)
	if err := binary.Write(&payload, binary.BigEndian, uint16(len(lastnameBytes))); err != nil {
		return nil, err
	}
	if _, err := payload.Write(lastnameBytes); err != nil {
		return nil, err
	}

	// 4) document (uint64)
	if err := binary.Write(&payload, binary.BigEndian, bp.Document); err != nil {
		return nil, err
	}

	// 5) birth_len (uint16) + birth (bytes)
	birthBytes := []byte(bp.Birthdate) // ej: "YYYY-MM-DD" o el formato que uses
	if err := binary.Write(&payload, binary.BigEndian, uint16(len(birthBytes))); err != nil {
		return nil, err
	}
	if _, err := payload.Write(birthBytes); err != nil {
		return nil, err
	}

	// 6) number (uint16)
	if err := binary.Write(&payload, binary.BigEndian, bp.Number); err != nil {
		return nil, err
	}

	// Frame: [frame_len (uint16)] + payload
	data := payload.Bytes()
	var frame bytes.Buffer
	if err := binary.Write(&frame, binary.BigEndian, uint16(len(data))); err != nil {
		return nil, err
	}
	if _, err := frame.Write(data); err != nil {
		return nil, err
	}

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
