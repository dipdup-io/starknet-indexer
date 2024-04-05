# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.22.1-alpine as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

RUN apk --no-cache add ca-certificates
RUN apk add --update git musl-dev gcc build-base

RUN mkdir -p $GOPATH/src/github.com/dipdup-io/starknet-indexer/

COPY ./go.* $GOPATH/src/github.com/dipdup-io/starknet-indexer/
WORKDIR $GOPATH/src/github.com/dipdup-io/starknet-indexer
RUN go mod download

COPY cmd/indexer cmd/indexer
COPY internal internal
COPY pkg pkg

WORKDIR $GOPATH/src/github.com/dipdup-io/starknet-indexer/cmd/indexer/
RUN go build -a -o /go/bin/indexer .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/starknet-indexer/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/indexer /go/bin/indexer
COPY ./build/dipdup.yml ./
COPY ./build/interfaces ./interfaces
COPY ./build/bridged_tokens ./
COPY ./cmd/indexer/graphql ./graphql
COPY ./cmd/indexer/views ./views

ENTRYPOINT ["/go/bin/indexer", "-c", "dipdup.yml"]