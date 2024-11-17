EXEC := charter.exe

CGO_ENABLED := 1

export CGO_ENABLED

include .env
export

.PHONY: tidy
tidy:
	go mod tidy -v


.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...


${EXEC}:
	go build -o=${EXEC}


.PHONY: run
run: ${EXEC}
	${EXEC}


.PHONY: clean
clean:
	del ${EXEC}