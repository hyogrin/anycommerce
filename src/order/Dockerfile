FROM golang:1.18 as build-deps
RUN mkdir -p /go/src/github.com/anycommerce/order
WORKDIR /go/src/github.com/anycommerce/order
COPY . /go/src/github.com/anycommerce/order/
RUN go env -w GOPROXY=https://goproxy.io,direct && go get ./...
RUN make build

FROM alpine
COPY --from=build-deps /go/src/github.com/anycommerce/order/order /bin/
COPY --from=build-deps /go/src/github.com/anycommerce/order/data/*.* /bin/data/
EXPOSE 80
WORKDIR /bin
CMD [ "/bin/order" ] 