FROM alpine

ENV PROXYLIST=proxylist.txt \
    URL=https://m.vk.com \
    APIURL=http://localhost:3000/api/v1/addproxy \
    TREDS=50

COPY . /usr/local/bin

RUN chmod +x /usr/local/bin/*

EXPOSE 3000

ENTRYPOINT ["docker-entrypoint.sh"]
