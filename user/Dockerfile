FROM golang:latest
WORKDIR /root/micor_go
COPY / /root/micor_go
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o user
ENTRYPOINT ["./user"]
