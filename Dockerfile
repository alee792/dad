FROM golang as builder

ADD . /
RUN go build ./cmd/dad

FROM scratch
WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dad /
COPY --from=builder /bin /

ENTRYPOINT ["./dad"]
