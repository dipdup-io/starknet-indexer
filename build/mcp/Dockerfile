# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.24.2-alpine AS builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

RUN apk --no-cache add ca-certificates
RUN apk add --update git musl-dev gcc build-base

RUN mkdir -p $GOPATH/src/github.com/dipdup-io/starknet-indexer/

COPY ./go.* $GOPATH/src/github.com/dipdup-io/starknet-indexer/
WORKDIR $GOPATH/src/github.com/dipdup-io/starknet-indexer
RUN go mod download

COPY cmd/mcp cmd/mcp
COPY internal internal
COPY pkg pkg

WORKDIR $GOPATH/src/github.com/dipdup-io/starknet-indexer/cmd/mcp/
RUN go build -a -o /go/bin/mcp .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/mcp/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/mcp /go/bin/mcp
COPY ./build/dipdup.yml ./

ENTRYPOINT ["/go/bin/mcp", "-c", "dipdup.yml"]