FROM alpine:3.19
RUN apk add --no-cache socat
USER 1000
EXPOSE 2000
ENTRYPOINT socat
CMD tcp-l:2000,fork exec:'/bin/cat'
