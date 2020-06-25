FROM golang:1.14-alpine as builder

# Force Go to use the cgo based DNS resolver. This is required to ensure DNS
# queries required to connect to linked containers succeed.
ENV GODEBUG netdns=cgo

COPY . /naporta-api

# Install dependencies and build the binaries.
RUN apk add --no-cache --update alpine-sdk \
    gcc

RUN cd /naporta-api && \
        go install .

# Start a new, final image.
FROM alpine as final

RUN apk --no-cache add \
    bash

# Copy the binaries from the builder image.
COPY --from=builder /go/bin/naporta-api /bin/naporta-api
COPY --from=builder /naporta-api/naporta-api.conf /bin/naporta-api.conf

ENTRYPOINT ["./naporta-api"]
