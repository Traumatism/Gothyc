package main

import (
	"fmt"
	"os"
	"strings"

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

func main() {
	fmt.Printf("%s", banner)

	parser := argparse.NewParser("Gothyc", "A Minecraft port scanner written in Go. 🐹")
	target := parser.String("t", "target", &argparse.Options{Required: true, Help: "Target CIDR"})
	port_range := parser.String("p", "ports", &argparse.Options{Required: true, Help: "Ports to scan"})
	threads := parser.Int("", "threads", &argparse.Options{Required: true, Help: "Threads ammount"})
	timeout := parser.Int("", "timeout", &argparse.Options{Required: true, Help: "Timeout in milliseconds"})
	retries := parser.Int("", "retries", &argparse.Options{Required: false, Help: "Number of times Gothyc will ping a target", Default: 3})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	hosts, ports := parse_target(*target), parse_port(*port_range)
	output_file := fmt.Sprintf("%s.gothyc.txt", strings.ReplaceAll(*target, "/", "_"))

	os.OpenFile(output_file, os.O_RDONLY|os.O_CREATE, 0755)

	gologger.Info().Msg("Output file set to `" + output_file + "`")
	gologger.Info().Msg(fmt.Sprintf("`%d * %d = %d` servers will be scanned", len(hosts), len(ports), len(hosts)*len(ports)))
	gologger.Info().Msg("Starting scan...")

	s := goccm.New(*threads)

	for _, host := range hosts {
		for _, port := range ports {
			s.Wait()

			go func(host string, port int, timeout int) {
				scan_port(host, port, timeout, output_file, *retries)
				s.Done()
			}(host, port, *timeout)

		}
	}

	s.WaitAllDone()

	gologger.Info().Msg("Scan finished")
}
