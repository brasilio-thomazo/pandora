FROM alpine:latest AS build

WORKDIR /build
RUN apk add --no-cache ca-certificates git curl bash build-base \
    && curl -sL https://dl.google.com/go/go1.25.3.linux-amd64.tar.gz | tar -C /usr/local -xzf - \
    && export PATH=$PATH:/usr/local/go/bin

ENV PATH=$PATH:/usr/local/go/bin
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o cerberus -trimpath -ldflags "-s -w" .

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=build /build/cerberus /usr/local/bin/cerberus
CMD ["cerberus"]
EXPOSE 8080
