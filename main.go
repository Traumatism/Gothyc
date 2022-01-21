package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/akamensky/argparse"
	"github.com/projectdiscovery/gologger"
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
	var last int

	for {
		last = scanned
		gologger.Info().Msgf("%d/%d (%d%%)", scanned, total, uint64(float64(scanned)/float64(total)*100.0))

		if time.Sleep(time.Second * 5); scanned == total || last == scanned {
			break
		}
	}
}

func main() {
	fmt.Printf("%s", banner)

	parser := argparse.NewParser("Gothyc", "A Minecraft port scanner written in Go. 🐹")

	target := parser.String("t", "target", &argparse.Options{Required: true, Help: "Target CIDR or file with CIDRs"})

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

	gologger.Info().Msg("Starting scan...")

	go status()

	s := make(chan struct{}, *threads)

	var wg sync.WaitGroup

	for _, host := range hosts {
		for _, port := range ports {
			s <- struct{}{}
			wg.Add(1)

			go func(host string, port int) {
				defer func() { <- s }()
				scanned++
				scan_port(host, port, *timeout, output, *retries, *output_fmt)
				wg.Done()
			}(host, port)
		}
	}

	gologger.Info().Msgf("Waiting for threads to finish...")
	wg.Wait()

	gologger.Info().Msg("Scan finished")
}
