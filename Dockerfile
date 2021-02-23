FROM alpine:latest

#LABEL maintainer="litong01@us.ibm.com"
LABEL maintainer="pantazheng@vip.qq.com"

ENV PYTHONUNBUFFERED=1

RUN sed -i "s@dl-cdn.alpinelinux.org/@mirrors.aliyun.com/@g" /etc/apk/repositories && \
    apk update && apk upgrade && \
    apk add --no-cache bash ansible docker-cli openssl xxd dos2unix coreutils && \
    if [ ! -e /usr/bin/python ]; then ln -sf python3 /usr/bin/python ; fi && \
    mkdir -p /usr/lib/python3.8/site-packages/Crypto/Random/Fortuna && \
    ansible-galaxy collection install community.general

COPY . /home
COPY plugins /usr/lib/python3.8/site-packages/ansible/plugins
COPY pypatch /usr/lib/python3.8/site-packages/Crypto/Random/Fortuna
RUN rm -rf /var/cache/apk/* && rm -rf /tmp/* && apk update && \
    dos2unix -q /home/main.sh /home/scripts/mainfuncs.sh \
    /usr/lib/python3.8/site-packages/ansible/plugins/callback/minifab.py && \
    apk del dos2unix && rm -rf /var/cache/apk/* && rm -rf /tmp/*

ENV PATH $PATH:/home/bin
WORKDIR /home