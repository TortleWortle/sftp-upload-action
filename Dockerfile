FROM golang:1.17

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -o sftp_upload -v ./src/...

ENTRYPOINT ["./sftp_upload"]