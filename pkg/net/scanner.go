package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/traumatism/gothyc/pkg/utils"
)

type Scanner struct {
	Hosts         []string
	Ports         []int
	Timeout       int
	Output_format string
	Output_file   string
	Retries       int
	Workers       int
	Active        int
}

func (s *Scanner) Ping(target string) (string, error) {
	conn, err := net.DialTimeout("tcp", target, time.Duration(s.Timeout)*time.Millisecond)

	if err != nil {
		return "", err
	}

	// lazy pkt generation
	if _, err := conn.Write([]byte("\x07\x00/\x01_\x00\x01\x01\x01\x00")); err != nil {
		return "", err
	}

	total_lenght, err := ReadVarint(conn)

	if err != nil {
		return "", err
	}

	buf_total := bytes.NewBuffer(nil)

	if _, err = io.CopyN(buf_total, conn, int64(total_lenght)); err != nil {
		return "", err
	}

	packet_id, err := ReadVarint(buf_total)

	if err != nil || uint32(packet_id) != uint32(0x00) {
		return "", err
	}

	lenght, err := ReadVarint(buf_total)

	if err != nil {
		return "", err
	}

	buf_data := make([]byte, lenght)

	if err != nil {
		return "", err
	}

	max, err := buf_total.Read(buf_data)

	if err != nil {
		return "", err
	}

	defer conn.Close()

	return string(buf_data[:max]), nil
}

func (s *Scanner) ScanTarget(host string, port int) {
	target := fmt.Sprintf("%s:%d", host, port)
	raw_data, err := s.Ping(target)

	if err != nil {
		return
	}
	data := &utils.Response{}

	if err = json.Unmarshal([]byte(raw_data), data); err != nil {
		return
	}

	var raw_motd string

	results := &utils.ReponseMOTD{}

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

	output_result := utils.OutputResult{
		Target:      target,
		Version:     data.Version.Name,
		Players:     fmt.Sprintf("%d/%d", data.Players.Online, data.Players.Max),
		Description: motd,
	}

	output_str := utils.FormatQubo(output_result)

	fmt.Printf("%s\n", output_str)

	for {
		f, err := os.OpenFile(s.Output_file, os.O_APPEND|os.O_WRONLY, 0600)

		if err != nil {
			continue
		}

		if s.Output_format == "csv" {
			output_str = utils.FormatCSV(output_result)

		} else if s.Output_format == "json" {
			output_str = utils.FormatJSON(output_result)
		}

		f.WriteString(fmt.Sprintf("%s\n", output_str))
		break
	}
}

func (s *Scanner) Scan() {
	var wg sync.WaitGroup

	for _, host := range s.Hosts {
		for _, port := range s.Ports {
			for {
				if s.Active <= s.Workers {
					wg.Add(1)

					s.Active++

					go func(host string, port int) {
						s.ScanTarget(host, port)
						s.Active--
						defer wg.Done()
					}(host, port)
					break

				}
			}
		}
	}

	wg.Wait()
}
