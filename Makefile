install:
	git fetch && git pull
	cd cmd/gothyc && go build -v && go install -v && cd ../..

upload:
	git add . && git commit -m "update" && git push

update:
	git fetch && git pull
