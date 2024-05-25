FROM alpine:3.19
RUN apk add --no-cache socat
USER 1000
EXPOSE 2000
ENTRYPOINT /bin/sh
CMD -c socat tcp-l:2000,fork exec:'/bin/cat'
