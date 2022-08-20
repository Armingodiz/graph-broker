FROM golang:1.16-alpine AS go-build

ENV GO111MODULE=on

WORKDIR /app

COPY ./receiver/go.mod .

RUN go mod download

COPY ./receiver .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o receiver .

FROM alpine:3.9.5 as dns
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=go-build /app .
ENTRYPOINT ["./receiver"]