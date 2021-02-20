FROM golang:alpine AS mongo-builder

USER root
WORKDIR /mongodb-go-sample
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add make

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./ ./
RUN rm mongodb-go-sample; make test && make

FROM scratch

WORKDIR /mongodv-go-sample

COPY --from=mongo-builder /mongodb-go-sample/mongodb-go-sample .
CMD ["./mongodb-go-sample"]
