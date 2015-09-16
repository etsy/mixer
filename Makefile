.PHONY: build

build:
	go build -o ./bin/mixer mixer.go
	go build -o ./bin/match match/match.go
	go build -o ./bin/cleanup cleanup/cleanup.go
	go build -o ./bin/staff staff/staff.go

clean:
	rm ./bin/*

sql:
	mysqldump -u root --opt mixer -d --single-transaction | sed 's/ AUTO_INCREMENT=[0-9]*\b//' > db/mixer.sql

newdep:
	godep save -r ./... 

# run like: make updep import=github.com/codahale/metrics
updep:
	go get -u -v -d $(import)
	godep update $(import)
	godep save -r ./...
