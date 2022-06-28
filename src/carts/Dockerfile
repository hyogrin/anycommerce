FROM golang:1.18 as build-deps
RUN mkdir -p /go/src/github.com/anycommerce/carts
WORKDIR /go/src/github.com/anycommerce/carts
COPY . /go/src/github.com/anycommerce/carts/
RUN go env -w GOPROXY=https://goproxy.io,direct && go get ./...
RUN make build

FROM alpine
COPY --from=build-deps /go/src/github.com/anycommerce/carts/carts /bin/
EXPOSE 80
WORKDIR /bin
CMD [ "/bin/carts" ]