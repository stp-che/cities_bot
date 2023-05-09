FROM golang:1.20 AS deps

WORKDIR /go/src/citiesbot

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod tidy

FROM deps AS build

COPY . /go/src/citiesbot
WORKDIR /go/src/citiesbot

RUN mkdir bin
RUN go build -o ./bin/app ./cmd/app

FROM debian:buster-slim

RUN apt-get update && apt-get install file -y

COPY --from=build /go/src/citiesbot/bin/app /app/

RUN chmod +x /app/app
