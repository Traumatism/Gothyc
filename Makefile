install:
	git fetch && git pull
	go build -v && go install -v

upload:
	git add . && git commit -m "update" && git push

update:
	git fetch && git pull
