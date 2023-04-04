VERSION 0.6

FROM golang:1.17-bullseye

build:
  WORKDIR /app
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN go build -o ./build/frp-manager ./main.go

  SAVE ARTIFACT build /build AS LOCAL ./build
