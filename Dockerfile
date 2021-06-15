FROM golang:1.13 AS build
RUN useradd -u 10001 benthos
WORKDIR /build/
COPY . /build/
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build .


FROM busybox AS package
WORKDIR /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /zoneinfo.zip
COPY --from=build /build/sheets-plugin .

ENV ZONEINFO=/zoneinfo.zip
EXPOSE 4195

ENTRYPOINT ["/sheets-plugin"]
CMD ["-c", "/benthos.yaml"]