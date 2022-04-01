FROM gcr.io/coder-dogfood/master/coder-dev-ubuntu

COPY main.go /tmp/main.go

COPY funnyd.service /etc/systemd/system/funnyd.service

USER root

RUN set -euo pipefail && \
    go build -o /tmp/systemd-funnyd /tmp/main.go && \
    chmod +x /tmp/systemd-funnyd && \
    systemctl enable funnyd

USER coder
