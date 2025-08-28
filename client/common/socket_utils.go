package common

import (
	"errors"
	"net"
)

// Made to avoid short writes.
func mustWriteAll(conn net.Conn, data []byte) error {
	total := len(data)
	sent := 0
	for sent < total {
		pending := data[sent:]
		n, err := conn.Write(pending)
		if err != nil {
			return err
		}
		sent += n
	}
	return nil
}

// Made to avoid short reads.
func mustReadAll(conn net.Conn, expected_length int) ([]byte, error) {
	data := make([]byte, expected_length)
	read := 0
	for read < expected_length {
		n, err := conn.Read(data[read:])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, errors.New("connection closed")
		}
		read += n
	}
	return data, nil
}
