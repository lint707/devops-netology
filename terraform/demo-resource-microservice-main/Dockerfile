FROM golang:1.17.1-alpine3.14 AS builder

WORKDIR /usr/local/go/src/

COPY ./app/ /usr/local/go/src/

RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -o gobegin cmd/gobegin/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /usr/local/go/src/gobegin /app/
# uncomment this string if you are use http calls
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app/gobegin"]