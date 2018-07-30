FROM quay.io/vxlabs/dep as builder

RUN mkdir -p $GOPATH/src/github.com/vx-labs
WORKDIR $GOPATH/src/github.com/vx-labs/iot-mqtt-auth
RUN mkdir release
COPY Gopkg* ./
RUN dep ensure -vendor-only
COPY . ./
RUN go test ./... && \
    go build -buildmode=exe -a -o /bin/auth ./cmd/server

FROM alpine
ENTRYPOINT ["/usr/bin/server"]
RUN apk -U add ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=builder /bin/auth /usr/bin/server

