FROM golang:1.18 as build-stage

WORKDIR /src
COPY . .

WORKDIR /src/go
RUN make test
RUN make build
RUN ./resource_control_sample --help

FROM quay.io/prometheus/busybox-linux-amd64:glibc AS bin
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

COPY --from=build-stage /src/go/resource_control_sample /
COPY go/resources/ /resources/
RUN chmod +x /resource_control_sample

USER nobody
# free port see https://github.com/prometheus/prometheus/wiki/Default-port-allocations
EXPOSE 9905

ENTRYPOINT [ "/resource_control_sample" ]
