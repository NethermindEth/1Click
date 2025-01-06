FROM debian:stable-20221004-slim

WORKDIR /

ARG LH_VERSION
ARG NETWORK

RUN apt-get update && apt-get install -y curl dos2unix bash

RUN version=$(echo $LH_VERSION | cut -d':' -f 2) \
&& curl -L https://github.com/sigp/lighthouse/releases/download/$version/lighthouse-$version-aarch64-unknown-linux-gnu.tar.gz --output lh.gz

RUN tar -xvf lh.gz

COPY validator-init.sh .

RUN dos2unix validator-init.sh

RUN chmod +x validator-init.sh

ENV VAL_NETWORK=$NETWORK

CMD ["/bin/sh", "validator-init.sh"]
