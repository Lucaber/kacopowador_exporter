FROM golang:latest as build
ADD . /go/src/gitlab.micronited.de/lusu/kacopowador_exporter
WORKDIR /go/src/gitlab.micronited.de/lusu/kacopowador_exporter
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux
RUN go get
RUN go build
FROM scratch
COPY --from=build /go/src/gitlab.micronited.de/lusu/kacopowador_exporter/kacopowador_exporter /kacopowador_exporter
CMD ["/kacopowador_exporter"]