.PHONY: test fuzz lint

# go vet and the tests with the race detector.
test:
	go vet ./... && go test -race ./...

# Fuzz Open for a while (FUZZTIME=5m for longer).
FUZZTIME ?= 1m
fuzz:
	go test -run '^$$' -fuzz FuzzOpen -fuzztime $(FUZZTIME) .

lint:
	shellcheck scripts/*.sh && scripts/commit-type_test.sh && go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
