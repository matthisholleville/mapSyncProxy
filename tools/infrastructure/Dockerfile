FROM haproxytech/haproxy-alpine:2.8

COPY ./haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg

COPY ./dataplaneapi.yaml /etc/haproxy/dataplaneapi.yaml

RUN mkdir -p /usr/local/etc/haproxy/maps/ && touch /usr/local/etc/haproxy/maps/rate-limits.map

COPY ./rate-limits.txt /usr/local/etc/haproxy/maps/rate-limits.map