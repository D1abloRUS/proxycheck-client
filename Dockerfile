FROM alpine

ENV PROXYLIST=proxylist.txt \
    URL=https://m.vk.com \
    APIURL=http://localhost:3000 \
    TREDS=50

COPY . /usr/local/bin

RUN apk --no-cache add --update \
      openssl \
      ca-certificates \
    && wget http://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.mmdb.gz \
    && gunzip GeoLite2-Country.mmdb.gz \
    && rm -f GeoLite2-Country.mmdb.gz \
    && chmod +x /usr/local/bin/*

ENTRYPOINT ["docker-entrypoint.sh"]
