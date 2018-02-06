FROM vxlabs/glide as builder

RUN mkdir -p $GOPATH/src/github.com/vx-labs
WORKDIR $GOPATH/src/github.com/vx-labs/iot-mqtt-auth
RUN mkdir release
COPY glide* ./
RUN glide install
COPY . ./
RUN go test $(glide nv) && \
    go build -buildmode=exe -a -o /bin/auth ./cmd/server

FROM alpine
ENTRYPOINT ["/usr/bin/server"]
RUN apk -U add ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=builder /bin/auth /usr/bin/server

