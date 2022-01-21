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

func ping(target string, timeout int) (string, error) {
	conn, err := net.DialTimeout("tcp", target, time.Duration(timeout)*time.Millisecond)

	if err != nil {
		return "", err
	}

	// lazy pkt generation
	if _, err := conn.Write([]byte("\x07\x00/\x01_\x00\x01\x01\x01\x00")); err != nil {
		return "", err
	}

	total_lenght, err := read_varint(conn)

	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)

	if _, err = io.CopyN(buf, conn, int64(total_lenght)); err != nil {
		return "", err
	}

	packet_id, err := read_varint(buf)

	if err != nil || uint32(packet_id) != uint32(0x00) {
		return "", err
	}

	lenght, err := read_varint(buf)

	if err != nil {
		return "", err
	}

	buf_2 := make([]byte, lenght)

	if err != nil {
		return "", err
	}

	max, err := buf.Read(buf_2)

	if err != nil {
		return "", err
	}

	defer conn.Close()

	return string(buf_2[:max]), nil
}

func scan_port(ip string, port int, timeout int, output_file string, retries int, format string) {
	target := fmt.Sprintf("%s:%d", ip, port)

	var (
		raw_data string
		err      error
	)

	for i := 0; i <= retries; i++ {
		raw_data, err = ping(target, timeout)
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

	output_result := OutputResult{
		target:      target,
		version:     data.Version.Name,
		players:     fmt.Sprintf("%d/%d", data.Players.Online, data.Players.Max),
		description: motd,
	}

	output_str := format_qubo(output_result)

	fmt.Printf("%s\n", output_str)

	for {
		f, err := os.OpenFile(output_file, os.O_APPEND|os.O_WRONLY, 0600)

		if err != nil {
			continue
		}
		if format == "csv" {
			output_str = format_csv(output_result)
		} else if format == "json" {
			output_str = format_json(output_result)
		} else if format == "qubo" {
			output_str = format_qubo(output_result)
		}

		f.WriteString(fmt.Sprintf("%s\n", output_str))
		break
	}
}
