FROM golang:1.14 as builder

WORKDIR /go/src/kitt
COPY . .

RUN make

FROM debian:buster-slim

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends git; \
	rm -rf /var/lib/apt/lists/*;

COPY --from=builder /go/src/kitt/bin/kitt /usr/bin
ENTRYPOINT ["/usr/bin/kitt"]
