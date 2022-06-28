FROM golang:1.18 as build-deps
RUN mkdir -p /go/src/github.com/anycommerce/product
WORKDIR /go/src/github.com/anycommerce/product
COPY . /go/src/github.com/anycommerce/product/
RUN go env -w GOPROXY=https://goproxy.io,direct && go get ./...
RUN make build

FROM alpine
COPY --from=build-deps /go/src/github.com/anycommerce/product/product /bin/
COPY --from=build-deps /go/src/github.com/anycommerce/product/data/*.* /bin/data/
EXPOSE 80
WORKDIR /bin
CMD [ "/bin/product" ] 