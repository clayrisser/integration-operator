FROM node:12-alpine AS builder

RUN apk add --no-cache \
  make

WORKDIR /tmp/app

COPY package*.json Makefile *.mk /tmp/app/
RUN mkdir src && make install

FROM node:12-alpine

RUN apk add --no-cache \
  make

COPY docker/latest/entrypoint.sh /usr/local/sbin/entrypoint

RUN chmod +x /usr/local/sbin/entrypoint

COPY --from=builder /tmp/app /opt/app

WORKDIR /opt/app

COPY . /opt/app/
RUN make +build

ENV DEBUG_OPERATOR=false

ENTRYPOINT [ "/usr/local/sbin/entrypoint" ]
