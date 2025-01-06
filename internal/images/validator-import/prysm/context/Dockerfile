ARG PRYSM_VERSION
ARG NETWORK

FROM ${PRYSM_VERSION} as prysm-validator

FROM debian:buster-slim
COPY --from=prysm-validator /app/cmd/validator/validator validator

VOLUME [ "/keystore" ]

COPY validator-init.sh .

RUN chmod +x validator-init.sh

CMD ["/bin/sh", "validator-init.sh"]