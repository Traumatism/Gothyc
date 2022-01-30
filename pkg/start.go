package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/projectdiscovery/gologger"

	"github.com/traumatism/gothyc/pkg/net"
	"github.com/traumatism/gothyc/pkg/parse"
)

const banner = `
    __
 __/  \__       Gothyc   A Minecraft port scanner written in Go. 🐹
/  \__/  \__
\__/  \__/  \   Version  0.4.0
   \__/  \__/   Author   @toastakerman

`

func Start() {

	fmt.Printf("%s", banner)

	parser := argparse.NewParser("Gothyc", "A Minecraft port scanner written in Go. 🐹")

	target := parser.String("t", "target", &argparse.Options{Required: true, Help: "Target CIDR or file with CIDRs"})

	port_range := parser.String("p", "ports", &argparse.Options{Required: true, Help: "Ports to scan"})

	threads := parser.Int("c", "threads", &argparse.Options{Required: true, Help: "Threads ammount"})

	timeout := parser.Int("", "timeout", &argparse.Options{Required: true, Help: "Timeout in milliseconds"})

	retries := parser.Int("r", "retries", &argparse.Options{Required: false, Help: "Number of times Gothyc will ping a target", Default: 0})

	output_file := parser.String("o", "output", &argparse.Options{Required: false, Help: "Output file", Default: nil})

	output_fmt := parser.String("f", "format", &argparse.Options{Required: false, Help: "Output format (qubo/json/csv)", Default: "qubo"})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	hosts := parse.ParseTarget(*target)
	ports := parse.ParsePorts(*port_range)

	var output string

	if *output_file == "" {
		output = fmt.Sprintf("%s.gothyc.txt", strings.ReplaceAll(*target, "/", "_"))
	} else {
		output = *output_file
	}

	os.OpenFile(output, os.O_RDONLY|os.O_CREATE, 0755)

	gologger.Info().Msg("Output file set to '" + output + "'")

	scanner := &net.Scanner{
		Hosts:         hosts,
		Ports:         ports,
		Workers:       *threads,
		Timeout:       *timeout,
		Retries:       *retries,
		Output_file:   output,
		Output_format: *output_fmt,
	}

	scanner.Scan()

	gologger.Info().Msg("Scan finished")
}
