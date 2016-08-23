install:
	go install ./cmd/...

build:
	for dir in $(shell ls cmd); do \
		cd cmd/$$dir && go build -v -o bin/$$dir && cd ../..; \
	done

run: build
	cmd/profilesvc/bin/profilesvc
