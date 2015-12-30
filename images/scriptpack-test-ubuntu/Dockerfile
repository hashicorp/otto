FROM ubuntu:14.04
MAINTAINER Mitchell Hashimoto <mitchell.hashimoto@gmail.com>

ENV BATS_VERSION 0.4.0
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -q \
	&& apt-get install -y -q --no-install-recommends bash make curl ca-certificates \
	&& curl -o "/tmp/v${BATS_VERSION}.tar.gz" -L \
		"https://github.com/sstephenson/bats/archive/v${BATS_VERSION}.tar.gz" \
	&& tar -x -z -f "/tmp/v${BATS_VERSION}.tar.gz" -C /tmp/ \
	&& bash "/tmp/bats-${BATS_VERSION}/install.sh" /usr/local \
	&& rm -rf /tmp/*
