FROM golang:1.23.0-alpine3.20 as deps
RUN apk add --update --no-cache ca-certificates git
ENV GOPATH=/go
WORKDIR /deps
COPY go.mod /deps
COPY go.sum /deps
RUN go mod download

FROM golang:1.23.0-alpine3.20 as builder
RUN apk add --update --no-cache ca-certificates git
COPY --from=deps /go /go
ENV GOPATH=/go
COPY . /alertfy
WORKDIR /alertfy
RUN mkdir bin -p && go build -o bin/main ./cmd/alertfy

FROM alpine:3.20
LABEL maintainer="Murtaza Udaipurwala <murtaza@murtazau.xyz>"
COPY --from=builder /alertfy/bin/main /amify/main
WORKDIR /alertfy
RUN adduser --disabled-password --no-create-home alertfy
USER alertfy
EXPOSE 5748
HEALTHCHECK --interval=10s CMD wget -q --spider http://localhost:5748/health || exit 1
ENTRYPOINT [ "./main" ]
