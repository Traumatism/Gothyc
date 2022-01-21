package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"time"
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
func ping(conn net.Conn) (string, error) {
	if _, err := conn.Write([]byte("\x07\x00/\x01_\x00\x01\x01\x01\x00")); err != nil {
		return "", err
	}

	lenght, err := read_varint(conn)

	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)

	if _, err = io.CopyN(buf, conn, int64(lenght)); err != nil {
		return "", err
	}

	packet_id, err := read_varint(buf)

	if err != nil || uint32(packet_id) != uint32(0x00) {
		return "", err
	}

	raw_data, err := read_string(buf)

	if err != nil {
		return "", err
	}

	defer conn.Close()

	return raw_data, nil
}

func scan_port(ip string, port int, timeout int, output_file string, retries int, format string) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, time.Duration(timeout)*time.Millisecond)

	if err != nil {
		return
	}

	var raw_data string

	for i := 0; i <= retries; i++ {
		raw_data, err = ping(conn)
		if err != nil {
			if i == retries {
				return
			}
			continue
		}
		break
	}

	data := &Response{}

	if err = json.Unmarshal([]byte(raw_data), data); err != nil {
		return
	}

	if data.Version.Name == "TCPShield.com" {
		return
	}

	var raw_motd string

	results := &ReponseMOTD{}

	if err = json.Unmarshal([]byte(raw_data), results); err != nil {
		var result map[string]interface{}
		json.Unmarshal([]byte(raw_data), &result)
		raw_motd = fmt.Sprintf("%s", result["description"])
	} else {
		raw_motd = results.Description.Text
	}

	var motd string

	motd = regexp.MustCompile(`§[a-fl-ork0-9]|\n`).ReplaceAllString(raw_motd, "")
	motd = regexp.MustCompile(`\ +|\t`).ReplaceAllString(motd, " ")

	t := OutputResult{
		target:      target,
		version:     data.Version.Name,
		players:     fmt.Sprintf("%d/%d", data.Players.Online, data.Players.Max),
		description: motd,
	}

	output_str := format_qubo(t)

	fmt.Printf("%s\n", output_str)

	for {
		f, err := os.OpenFile(output_file, os.O_APPEND|os.O_WRONLY, 0600)

		if err != nil {
			continue
		}
		if format == "csv" {
			output_str = format_csv(t)
		} else if format == "json" {
			output_str = format_json(t)
		} else if format == "qubo" {
			output_str = format_qubo(t)
		}

		f.WriteString(fmt.Sprintf("%s\n", output_str))
		break
	}
}
