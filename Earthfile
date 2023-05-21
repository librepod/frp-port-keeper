VERSION 0.6

FROM golang:1.18-bullseye

build:
  WORKDIR /app
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/frp-port-keeper ./main.go

  SAVE ARTIFACT build /build AS LOCAL ./build
