package validate

func ValidatePort(port int) bool {
	return port >= 1 && port <= 65535
}
