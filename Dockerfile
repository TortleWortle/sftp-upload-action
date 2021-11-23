FROM golang:1.17

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -o sftp_upload_action -v ./src/...

RUN mv sftp_upload_action /usr/bin

ENTRYPOINT ["sftp_upload_action"]