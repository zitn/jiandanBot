FROM golang as gobuilder
WORKDIR /test
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o main .
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

FROM scratch
COPY --from=gobuilder /test .
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gobuilder /etc_passwd /etc/passwd

COPY config.json .
USER nobody
CMD ["/main"]
