FROM golang:1.17 AS build

RUN useradd -u 10001 benthos

WORKDIR /build/
COPY . /build/

RUN apt-get update
RUN apt-get install ca-certificates -y
COPY ./certs/digicert.crt /usr/local/share/ca-certificates/digicert.crt
COPY ./certs/rapidssl.crt /usr/local/share/ca-certificates/rapidssl.crt
RUN update-ca-certificates

ARG ACCESS_TOKEN
ENV ACCESS_TOKEN=$ACCESS_TOKEN
RUN go env -w GOPRIVATE=github.com/ETP-Cloud
RUN git config --global url."https://ETP-Cloud:${ACCESS_TOKEN}@github.com".insteadOf "https://github.com"
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM busybox AS package

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /build/devdevman .
COPY ./public /public

USER benthos

EXPOSE 5000

ENTRYPOINT ["/devdevman"]

CMD [ ]
