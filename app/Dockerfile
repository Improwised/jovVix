# Base node 12 image
FROM node:18-alpine

WORKDIR /app

COPY ./package*.json ./

# Install dependencies; remove existing node_modules, install dependencies at a time
RUN npm install
ARG S6_OVERLAY_VERSION="3.1.6.2"
ARG S6_OVERLAY_ARCH="x86_64"
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz /tmp
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-${S6_OVERLAY_ARCH}.tar.xz /tmp
RUN tar -C / -Jxpf /tmp/s6-overlay-noarch.tar.xz && \
  tar -C / -Jxpf /tmp/s6-overlay-${S6_OVERLAY_ARCH}.tar.xz && \
  apk add --no-cache nginx curl && \
  rm -rf /tmp/*
RUN apk add --no-cache nginx supervisor wget curl
RUN mkdir -p /run/nginx
COPY ./docker/nginx.conf /etc/nginx/nginx.conf
COPY ./s6-overlay /etc/s6-overlay
COPY .env.example .env
COPY . ./
RUN npm run build
EXPOSE 5000
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["export $(xargs < .env) && /init"]
