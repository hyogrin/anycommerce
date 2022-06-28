FROM golang:1.18 as build-deps
RUN mkdir -p /go/src/github.com/anycommerce/user
WORKDIR /go/src/github.com/anycommerce/user
COPY . /go/src/github.com/anycommerce/user/
RUN go env -w GOPROXY=https://goproxy.io,direct && go get ./...
RUN make build

FROM alpine
COPY --from=build-deps /go/src/github.com/anycommerce/user/user /bin/
COPY --from=build-deps /go/src/github.com/anycommerce/user/data/*.* /bin/data/
EXPOSE 80
WORKDIR /bin
CMD [ "/bin/user" ] 