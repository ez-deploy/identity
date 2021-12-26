FROM golang:1.17 as builder

WORKDIR /go/src/github.com/ez-deploy/identity
COPY . .

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.io && \
    CGO_ENABLED=0 go build -tags netgo -o identity ./main.go

FROM busybox

WORKDIR /

COPY --from=builder /go/src/github.com/ez-deploy/identity/identity /identity

EXPOSE 80
ENTRYPOINT [ "/identity" ]