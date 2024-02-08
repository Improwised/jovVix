FROM node:18-alpine

COPY ./package.json ./package-lock.json ./

RUN npm install -g json-server

WORKDIR /data

EXPOSE 5000

ENTRYPOINT ["json-server"]

CMD ["--help"]
