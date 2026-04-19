default: testacc

# Run acceptance tests
.PHONY: testacc
testacc: TF_ACC=1
testacc:
	gotestsum --no-color=false -ftestname -- -race -count=1 ./...

# Run unit tests only (no acceptance tests)
.PHONY: test
test:
	gotestsum --no-color=false -ftestname -- -race ./...
