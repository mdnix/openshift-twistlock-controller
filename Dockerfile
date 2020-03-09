FROM golang:1.13.6 AS build

WORKDIR /tmp/controller-build
ADD . /tmp/controller-build/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:latest
COPY --from=build /tmp/controller-build/twistlock-controller /opt/controller/
COPY config.yaml /opt/controller/
RUN apk upgrade --available
WORKDIR /opt/controller/
CMD ["./twistlock-controller"]
