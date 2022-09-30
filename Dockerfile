FROM golang:1.19-alpine
WORKDIR /go/src/app
COPY go.* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go install

FROM alpine:latest
COPY entrypoint.sh /
COPY --from=0 /go/bin/firestore-wire /bin/firestore-wire
ENTRYPOINT ["/entrypoint.sh"]
