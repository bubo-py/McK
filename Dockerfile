FROM golang:1.19-alpine as builder
COPY go.mod go.sum /go/src/github.com/bubo-py/McK/
WORKDIR /go/src/github.com/bubo-py/McK/
RUN go mod download
COPY . /go/src/github.com/bubo-py/McK/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/McK /go/src/github.com/bubo-py/McK/cmd/main.go

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/bubo-py/McK/build/McK /usr/bin/McK

ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip

EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/McK"]
