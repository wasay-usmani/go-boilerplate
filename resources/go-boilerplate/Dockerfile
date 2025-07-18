ARG GO_VERSION=1.24.2

# builder
FROM golang:${GO_VERSION}-alpine as builder

RUN apk add --no-cache git bash sed build-base

RUN mkdir -p /build

WORKDIR /build

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS="-mod=vendor"

RUN go build -v -a -o go-boilerplate \
    -ldflags "-w -extldflags '-static'" -a -tags netgo \
    /build/cmd/go-boilerplate/main.go

# actual container
FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /build/go-boilerplate .

EXPOSE 8080/tcp

ENTRYPOINT [ "./go-boilerplate" ]