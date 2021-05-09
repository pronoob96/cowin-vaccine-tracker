FROM golang:alpine

WORKDIR /build

COPY . .

RUN apk add --update make
RUN make build

CMD ["./main"]