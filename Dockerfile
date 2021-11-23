FROM golang:1.17 AS builder

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -o sftp_upload_action -v ./src/...

FROM ubuntu:latest

COPY --from=builder /go/src/app/sftp_upload_action ./

RUN mv sftp_upload_action /usr/bin

ENTRYPOINT ["sftp_upload_action"]
