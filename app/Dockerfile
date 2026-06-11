FROM improwised/docker-nodejs-base:22.11.0-latest
WORKDIR /app
COPY ./package*.json ./
RUN npm install
RUN apk add --no-cache nginx supervisor wget curl
RUN mkdir -p /run/nginx
COPY ./docker/nginx.conf /etc/nginx/nginx.conf
COPY ./s6-overlay /etc/s6-overlay
COPY . ./
RUN npm run build
EXPOSE 5000
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["/init"]
