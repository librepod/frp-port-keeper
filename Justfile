@start:
  watchexec -r -e go -- go run .

@test:
  go test -v ./...

@build:
  go build -o bin/frp-port-keeper ./main.go
