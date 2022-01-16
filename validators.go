package main

func validate_tcp_port(port int) bool {
	return port >= 1 && port <= 65535
}
