package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

type Response struct {
	Players struct {
		Online int `json:"online"`
		Max    int `json:"max"`
	} `json:"players"`

	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
}

func scan_port(ip string, port int, timeout int) {
	conn, err := net.DialTimeout("tcp", ip+":"+fmt.Sprintf("%d", port), time.Duration(timeout)*time.Millisecond)

	if err != nil {
		return
	}

	_, err = conn.Write([]byte("\x07\x00/\x01_\x00\x01\x01\x01\x00")) // SLP packet + request status

	if err != nil {
		return
	}

	l, err := readUnsignedVarInt(conn)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)

	if _, err = io.CopyN(buf, conn, int64(l)); err != nil {
		return
	}

	p, err := readUnsignedVarInt(buf)

	if err != nil || uint32(p) != uint32(0x00) {
		return
	}

	d, err := readString(buf)

	if err != nil {
		return
	}

	defer conn.Close()

	data := &Response{}

	if err = json.Unmarshal([]byte(d), data); err != nil {
		return
	}

	var motd string

	type ReponseMOTD struct {
		Description struct {
			Text string `json:"text"`
		}
	}

	results := &ReponseMOTD{}

	if err = json.Unmarshal([]byte(d), results); err != nil {
		var result map[string]interface{}
		json.Unmarshal([]byte(d), &result)
		motd = fmt.Sprintf("%s", result["description"])
	} else {
		motd = results.Description.Text
	}

	fmt.Printf("(%s:%d)(%d/%d)(%s)(%s)\n", ip, port, data.Players.Online, data.Players.Max, data.Version.Name, motd)
}
