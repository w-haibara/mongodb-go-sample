FROM golang:alpine AS mongo-go-builder

USER root
WORKDIR /mongo-go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add make

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./ ./
RUN rm mongo-go; make test && make

FROM scratch

WORKDIR /mongo-go

COPY --from=mongo-go-builder /mongo-go/mongo-go .
CMD ["./mongo-go"]
