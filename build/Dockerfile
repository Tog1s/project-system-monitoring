FROM golang:1.21 as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /go/src

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /go/src
RUN go build -o /server cmd/server/*


FROM scratch
COPY --from=builder /server /server
COPY ./configs/config.yaml /configs/config.yaml

ENTRYPOINT ["/server"]

