package parse

import (
	"bufio"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/projectdiscovery/gologger"
	"github.com/traumatism/gothyc/pkg/validate"
)

func inc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		if ip[i]++; ip[i] > 0 {
			break
		}
	}
}

func ParseTarget(target string) []string {
	var ips []string

	if _, err := os.Stat(target); err == nil {
		file, err := os.Open(target)

		if err != nil {
			gologger.Fatal().Msg("Failed to open target file")
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			ips = append(ips, ParseTarget(scanner.Text())...)
		}

		return ips
	}

	if strings.Contains(target, "/") {
		ip, ipnet, err := net.ParseCIDR(target)
		if err != nil {
			gologger.Fatal().Msg("Error parsing CIDR: " + err.Error())
		}
		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
			ips = append(ips, ip.String())
		}
	} else {
		ips = append(ips, target)
	}

	return ips
}

func ParsePorts(port string) []int {
	ports := []int{}

	for _, port_range := range strings.Split(port, ",") {
		if strings.Contains(port_range, "-") {
			port_range_split := strings.Split(port_range, "-")

			start, err := strconv.Atoi(port_range_split[0])

			if err != nil {
				gologger.Fatal().Msg("Error parsing port range: " + err.Error())
			}

			end, err := strconv.Atoi(port_range_split[1])

			if err != nil {
				gologger.Fatal().Msg("Error parsing port range: " + err.Error())
			}

			for i := start; i <= end; i++ {
				if validate.ValidatePort(i) {
					ports = append(ports, i)
				}
			}

		} else {
			port_int, err := strconv.Atoi(port_range)

			if err != nil {
				gologger.Fatal().Msg("Error parsing port range: " + err.Error())
			}

			if validate.ValidatePort(port_int) && err == nil {
				ports = append(ports, port_int)
			}
		}
	}

	return ports
}
