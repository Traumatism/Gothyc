compile:
	go build -o gothyc
	./gothyc -help

scan:
	go build -o gothyc
	./gothyc -ports 25560-25580 -target 164.132.200.0/24 -threads 100
