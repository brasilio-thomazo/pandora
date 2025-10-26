FROM alpine:latest AS build

WORKDIR /usr/src
RUN apk add --no-cache ca-certificates build-base \
    icu-dev linux-headers readline-dev zlib-dev openssl-dev \
    ossp-uuid-dev libxml2-dev libxslt-dev curl git util-linux-dev \
    llvm-dev clang-dev curl-dev lz4-dev zstd-dev \
    flex bison bash liburing-dev perl \
    && git clone --depth 1 https://github.com/postgres/postgres.git \
    && cd postgres \
    && ./configure \
    --prefix=/usr/local \
    --with-icu \
    --with-uuid=ossp \
    --with-ssl=openssl \
    --with-libxml \
    --with-libxslt \
    --with-zlib \
    --with-liburing \
    --with-libcurl \
    --with-llvm \
    --with-lz4 \
    --with-zstd \
    && make -j$(nproc) \
    && make install

FROM alpine:latest
ENV PGDATA=/var/lib/postgresql/data
ENV POSTGRES_USER=postgres
ENV POSTGRES_DB=postgres
ENV POSTGRES_SLOT_NAME=replication_slot
ENV POSTGRES_REPLICATION_USER=replicator
COPY --from=build /usr/local /usr/local
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN apk add --no-cache ca-certificates bash icu-libs icu-data-full readline zlib openssl ossp-uuid-libs \
    libxml2 libxslt llvm17-libs clang17-libs lz4-libs liburing zstd-libs libcurl gosu \
    && chmod +x /usr/local/bin/entrypoint.sh \
    && mkdir -p /var/lib/postgresql /var/run/postgresql /etc/postgresql /docker-entrypoint-initdb.d \
    && addgroup postgres \
    && adduser -D -G postgres postgres \
    && chown -R postgres:postgres /var/lib/postgresql \
    && chown -R postgres:postgres /var/run/postgresql \
    && chown -R postgres:postgres /etc/postgresql \
    && sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" /usr/local/share/postgresql/postgresql.conf.sample \
    && sed -i "s/#unix_socket_directories = '\/tmp'/unix_socket_directories = '\/var\/run\/postgresql,\/tmp'/" /usr/local/share/postgresql/postgresql.conf.sample \
    && sed -i "s/#log_destination = 'stderr'/log_destination = 'stderr'/" /usr/local/share/postgresql/postgresql.conf.sample \
    && sed -i "s/#logging_collector = off/logging_collector = off/" /usr/local/share/postgresql/postgresql.conf.sample \
    && sed -i "s/#log_min_messages = warning/log_min_messages = info/" /usr/local/share/postgresql/postgresql.conf.sample \
    && cp /usr/local/share/postgresql/postgresql.conf.sample /etc/postgresql/postgresql.conf \
    && chown postgres:postgres /etc/postgresql/postgresql.conf

ENTRYPOINT [ "entrypoint.sh" ]
CMD [ "postgres" ]
