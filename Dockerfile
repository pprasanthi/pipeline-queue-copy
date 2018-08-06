FROM golang:1.10 as Builder

ENV PROJECT /go/src/gitlab.com/fenrirunbound/pipeline-queue
ADD . $PROJECT
WORKDIR $PROJECT

RUN go get -u github.com/golang/dep/cmd/dep
RUN make clean
RUN make vendor
RUN make test

