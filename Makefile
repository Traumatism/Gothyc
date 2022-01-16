compile:
	go build -o gothyc.out
	./gothyc.out -help

scan:
	go build -o gothyc.out
	./gothyc.out --ports 20000-20005,25000-25005,25560-25580,25665 \
				 --target 149.202.86.0/24 \
				 --threads 200 \
				 --timeout 5000

install:
	go build -v
	go install -v
