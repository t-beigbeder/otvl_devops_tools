FROM alpine:3.19
RUN apk add --no-cache socat
USER 1000
ENTRYPOINT /bin/sh
EXPOSE 2000
CMD echo socat tcp-l:2000,fork exec:'/bin/cat'
