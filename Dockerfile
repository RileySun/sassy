FROM golang:1.24rc2-bookworm AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o sassy

EXPOSE 7070
EXPOSE 8080
EXPOSE 9090

CMD ["/build/sassy"]