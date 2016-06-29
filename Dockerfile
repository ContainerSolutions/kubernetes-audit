FROM golang:1.6

COPY . /go/src/audit
WORKDIR /go/src/audit
RUN make build

CMD ["/go/src/audit/bin/audit"]

