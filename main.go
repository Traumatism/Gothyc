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
    __
 __/  \__       Gothyc   A Minecraft port scanner written in Go. 🐹
/  \__/  \__
\__/  \__/  \   Version  0.3.0
   \__/  \__/   Author   @toastakerman

`

var scanned int = 0
var total int

func status() {
	for {
		gologger.Info().Msgf("%d/%d (%d%%)", scanned, total, uint64(float64(scanned)/float64(total)*100.0))

		if time.Sleep(time.Second * 1); scanned == total || total-scanned < 1 {
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

	estimation := time.Now().Add(time.Duration((total**timeout) / *threads) * time.Millisecond)

	gologger.Info().Msgf("Starting scan... Estimated time of completion: %s", estimation.Format("15:04:05"))

	go status()

	wg := sync.WaitGroup{}
	ch := make(chan struct{}, *threads)

	for _, host := range hosts {
		for _, port := range ports {
			ch <- struct{}{}
			target := fmt.Sprintf("%s:%d", host, port)
			wg.Add(1)

			go func(target string) {
				defer wg.Done()

				scanned++

				scan_port(
					target, *timeout, output, *retries, *output_fmt,
				)
				<-ch
			}(target)
		}
	}

	wg.Wait()

	gologger.Info().Msgf("Waiting for threads to finish...")
	wg.Wait()

	gologger.Info().Msg("Scan finished")
}
