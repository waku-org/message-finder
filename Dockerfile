FROM golang:1.21 as builder

RUN mkdir -p /go/src/github.com/waku-org/test-waku-query/go

WORKDIR /go/src/github.com/waku-org/test-waku-query/go

ADD . .

RUN make

# Copy the binary to the second image
FROM debian:12.5-slim

LABEL source="https://github.com/waku-org/test-waku-query/go"
LABEL description="Storenode query tool"
LABEL commit="unknown"

COPY --from=builder /go/src/github.com/waku-org/test-waku-query/go/build/query /usr/local/bin/query

ENTRYPOINT ["/usr/local/bin/query"]
CMD ["-help"]
