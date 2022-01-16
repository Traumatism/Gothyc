package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/projectdiscovery/gologger"
)

type FullResponse struct {
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
	}

	Version struct {
		Name string `json:"name"`
	}

	Description string
}

type Response struct {
	Players struct {
		Online int `json:"online"`
		Max    int `json:"max"`
	} `json:"players"`

	Version struct {
		Name string `json:"name"`
	} `json:"version"`
}

type ReponseMOTD struct {
	Description struct {
		Text string `json:"text"`
	}
}

func scan_port(ip string, port int, timeout int, output_file string) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, time.Duration(timeout)*time.Millisecond)

	if err != nil {
		return
	}

	if _, err = conn.Write([]byte("\x07\x00/\x01_\x00\x01\x01\x01\x00")); err != nil {
		return
	}

	lenght, err := readUnsignedVarInt(conn)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)

	if _, err = io.CopyN(buf, conn, int64(lenght)); err != nil {
		return
	}

	packet_id, err := readUnsignedVarInt(buf)

	if err != nil || uint32(packet_id) != uint32(0x00) {
		return
	}

	raw_data, err := readString(buf)

	if err != nil {
		return
	}

	defer conn.Close()

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

	re := regexp.MustCompile(`§[a-fl-ork0-9]|\n`)
	motd := re.ReplaceAllString(raw_motd, "")

	output_str := fmt.Sprintf(
		"(%s)(%d/%d)(%s)(%s)\n",
		target,
		data.Players.Online, data.Players.Max, data.Version.Name, motd,
	)

	fmt.Print(output_str)

	f, err := os.OpenFile(output_file, os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		gologger.Fatal().Msg(err.Error())
		return
	}

	defer f.Close()

	if _, err = f.WriteString(output_str); err != nil {
		gologger.Fatal().Msg(err.Error())
		return
	}

}
