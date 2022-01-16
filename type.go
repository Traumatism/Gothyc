package main

import (
	"encoding/binary"
	"io"
)

type byteReaderWrap struct {
	reader io.Reader
}

func (w *byteReaderWrap) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	_, err := w.reader.Read(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], err
}

func read_varint(r io.Reader) (uint32, error) {
	v, err := binary.ReadUvarint(&byteReaderWrap{r})
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

func read_string(r io.Reader) (string, error) {

	l, err := read_varint(r)

	if err != nil {
		return "", err
	}
	buf := make([]byte, l)
	n, err := r.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}
