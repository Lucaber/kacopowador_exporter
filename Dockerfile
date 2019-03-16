FROM golang:latest as build
ADD . /go/src/gitlab.micronited.de/lusu/kacopowador_exporter
WORKDIR /go/src/gitlab.micronited.de/lusu/kacopowador_exporter
ENV GO111MODULE=on
RUN go get
RUN go build -o kacopowador_exporter
FROM scratch
COPY --from=build /go/src/gitlab.micronited.de/lusu/kacopowador_exporter/kacopowador_exporter /
CMD ["/kacopowador_exporter"]