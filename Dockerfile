# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY . .

RUN go mod download

COPY *.go ./

RUN go build -o /calibrator cmd/calibrator/main.go

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /calibrator /calibrator

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/calibrator"]