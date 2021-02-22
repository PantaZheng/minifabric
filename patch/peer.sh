#!/bin/sh

sed -i "s@dl-cdn.alpinelinux.org/@mirrors.aliyun.com/@g" \
    /etc/apk/repositories && apk update && apk upgrade && apk add coreutils

peer node start
