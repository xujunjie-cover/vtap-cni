FROM debian:stable-slim

COPY vtap-cni /usr/src/vtap-cni
COPY vtap-cni-daemon /usr/src/vtap-cni-daemon


WORKDIR /

ENTRYPOINT [ "/usr/src/vtap-cni-daemon" ]
