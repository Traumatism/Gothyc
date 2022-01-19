package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/projectdiscovery/gologger"
	"github.com/zenthangplus/goccm"
)

const banner = `
              ___________
_______ ________  /___  /______  ________
__  __ '/  __ \  __/_  __ \_  / / /  ___/ by @toastakerman
_  /_/ // /_/ / /_ _  / / /  /_/ // /__
_\__, / \____/\__/ /_/ /_/_\__, / \___/   A Minecraft port scanner
/____/                    /____/          written in Go. 🐹

`

var scanned int = 0
var total int

func status() {
	for {
		gologger.Info().Msgf("%d/%d (%f%%)", scanned, total, math.Round(float64(scanned)/float64(total)*100))
		time.Sleep(time.Second * 20)
	}
}

func main() {
	fmt.Printf("%s", banner)

	parser := argparse.NewParser("Gothyc", "A Minecraft port scanner written in Go. 🐹")

	target := parser.String("t", "target", &argparse.Options{Required: true, Help: "Target CIDR"})

	port_range := parser.String("p", "ports", &argparse.Options{Required: true, Help: "Ports to scan"})

	threads := parser.Int("c", "threads", &argparse.Options{Required: true, Help: "Threads ammount"})
	timeout := parser.Int("", "timeout", &argparse.Options{Required: true, Help: "Timeout in milliseconds"})
	retries := parser.Int("r", "retries", &argparse.Options{Required: false, Help: "Number of times Gothyc will ping a target", Default: 3})
	output_file := parser.String("o", "output", &argparse.Options{Required: false, Help: "Output file", Default: nil})
	output_fmt := parser.String("f", "format", &argparse.Options{Required: false, Help: "Output format (qubo/json/csv)", Default: "qubo"})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	hosts := parse_target(*target)
	ports := parse_port(*port_range)

	var output string

	if *output_file == "" {
		output = fmt.Sprintf("%s.gothyc.txt", strings.ReplaceAll(*target, "/", "_"))
	} else {
		output = *output_file
	}

	os.OpenFile(output, os.O_RDONLY|os.O_CREATE, 0755)

	gologger.Info().Msg("Output file set to '" + output + "'")

	total = len(hosts) * len(ports)

	gologger.Info().Msg(fmt.Sprintf("'%d * %d = %d' servers will be scanned", len(hosts), len(ports), total))
	gologger.Info().Msg("Starting scan...")

	go status()
	s := goccm.New(*threads)

	for _, host := range hosts {
		for _, port := range ports {
			s.Wait()

			go func(host string, port int) {
				scan_port(host, port, *timeout, output, *retries, *output_fmt)
				scanned++
				s.Done()
			}(host, port)

		}
	}

	s.WaitAllDone()

	gologger.Info().Msg("Scan finished")
}
