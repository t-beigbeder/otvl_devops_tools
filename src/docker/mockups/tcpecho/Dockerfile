FROM alpine:3.19
RUN apk add --no-cache netcat-openbsd
ENTRYPOINT bin/bash
CMD [ "ncat", "-l", "2000", "--keep-open", "--exec", "/bin/cat"]
