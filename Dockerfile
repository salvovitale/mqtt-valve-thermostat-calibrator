# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY . .

RUN go build -ldflags="-s -w" -mod vendor -o /calibrator cmd/calibrator/main.go

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /calibrator /calibrator

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/calibrator"]