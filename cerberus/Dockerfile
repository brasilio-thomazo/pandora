FROM alpine:latest AS build


RUN apk add --no-cache ca-certificates build-base \
    icu-dev linux-headers zlib-dev openssl-dev curl \
    && curl -sL https://dl.google.com/go/go1.25.3.linux-amd64.tar.gz | tar -C /usr/local -xzf - \
    && export PATH=$PATH:/usr/local/go/bin

ENV PATH=$PATH:/usr/local/go/bin
WORKDIR /usr/src/pandora/cerberus
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o cerberus -trimpath -ldflags "-s -w" .

FROM alpine:latest

RUN apk add --no-cache ca-certificates bash icu-libs icu-data-full zlib openssl curl
COPY --from=build /usr/src/pandora/cerberus/cerberus /usr/local/bin/cerberus

ENTRYPOINT [ "cerberus" ]

EXPOSE 8080