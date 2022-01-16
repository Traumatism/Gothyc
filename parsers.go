package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func inc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func parse_target(target string) []string {
	var ips []string

	if strings.Contains(target, "/") {
		ip, ipnet, err := net.ParseCIDR(target)
		if err != nil {
			fmt.Println(err)
		}
		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
			ips = append(ips, ip.String())
		}
	} else {
		ips = append(ips, target)
	}

	return ips
}

func parse_port(port string) []int {
	parsed := []int{}

	for _, port_range := range strings.Split(port, ",") {
		if strings.Contains(port_range, "-") {
			port_range_split := strings.Split(port_range, "-")
			start, _ := strconv.Atoi(port_range_split[0])
			end, _ := strconv.Atoi(port_range_split[1])
			for i := start; i <= end; i++ {
				if validate_tcp_port(i) {
					parsed = append(parsed, i)
				}
			}
		} else {
			port_int, _ := strconv.Atoi(port_range)
			if validate_tcp_port(port_int) {
				parsed = append(parsed, port_int)
			}
		}
	}

	return parsed
}
