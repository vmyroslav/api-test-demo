FROM alpine:3.21.0

# Packages
RUN apk add --no-cache wget unzip curl

# Set default arguments
ARG HOVERFLY_VERSION="v1.10.6"
ARG HOVERFLY_ADMIN_PORT=8888
ARG HOVERFLY_PROXY_PORT=8500

ENV HOVERFLY_ADMIN_PORT=${HOVERFLY_ADMIN_PORT}
ENV HOVERFLY_PROXY_PORT=${HOVERFLY_PROXY_PORT}

# Download and install both hoverfly and hoverctl
RUN wget -q "https://github.com/SpectoLabs/hoverfly/releases/download/v${HOVERFLY_VERSION#v}/hoverfly_bundle_linux_amd64.zip" && \
    unzip hoverfly_bundle_linux_amd64.zip -d /tmp && \
    mv /tmp/hoverfly /usr/local/bin/ && \
    mv /tmp/hoverctl /usr/local/bin/ && \
    chmod +x /usr/local/bin/hoverfly && \
    chmod +x /usr/local/bin/hoverctl && \
    rm -rf hoverfly_bundle_linux_amd64.zip /tmp/*

# Create default hoverctl config with environment variables
RUN mkdir -p /root/.hoverfly && \
    echo "hoverfly.host: localhost" > /root/.hoverfly/config.yaml && \
    echo "hoverfly.admin.port: \"${HOVERFLY_ADMIN_PORT}\"" >> /root/.hoverfly/config.yaml && \
    echo "hoverfly.proxy.port: \"${HOVERFLY_PROXY_PORT}\"" >> /root/.hoverfly/config.yaml

EXPOSE ${HOVERFLY_PROXY_PORT} ${HOVERFLY_ADMIN_PORT}

ENTRYPOINT ["hoverfly", "-listen-on-host=0.0.0.0"]
