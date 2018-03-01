FROM golang

LABEL maintainer="Witaha Honcharenko (vgoncharenko@magento.com)"

RUN set -xe \
	&& apt-get update && apt-get install -y \
	openssh-server

RUN set -xe \
    && printf "123123q\n123123q" | passwd root \
    && sed -i "s/#PermitRootLogin prohibit-password/PermitRootLogin yes/" /etc/ssh/sshd_config \
    && sed -i "s/#Port 22/Port 22/" /etc/ssh/sshd_config

RUN set -xe \
    && git clone https://github.com/vgoncharenko/google-page-speed.git \
    && go build google-page-speed/lib/googlepagespeed.go \
    && mv googlepagespeed /usr/bin/googlepagespeed \
    && chmod +x /usr/bin/googlepagespeed

WORKDIR /root
