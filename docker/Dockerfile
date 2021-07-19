FROM alpine:3.14
RUN apk add -U --no-cache ca-certificates

WORKDIR /data
ENV HOME=/data
COPY ./bin/starsd /usr/bin/starsd
COPY ./docker/entry-point.sh ./entry-point.sh
EXPOSE 26657

CMD ["starsd", "start", "--pruning", "nothing", "--log_format", "json"]
