# Base node 12 image
FROM node:18-alpine

# Selects env
ARG MODE='example'

WORKDIR /usr/src/app

COPY ./package*.json ./

# Install dependencies; remove existing node_modules, install dependencies at a time
RUN npm install

# Add extra tools below 3 lines for production deployment no need for locally
# RUN apk add --no-cache nginx supervisor wget curl

# # Nginx puts pid files
# RUN mkdir -p /run/nginx

# # Copy nginx config file
# COPY ./docker/nginx.conf /etc/nginx/nginx.conf

# Copy everything
COPY . ./

# Rename env file based on Mode
RUN mv .env.${MODE} .env

# Build project
RUN npm run build > build.log

# Run on port 3000
EXPOSE 3000

# Start nuxt app use if ngingx and supervisor no required
# CMD sleep 5000
CMD npm run start

# Start supervisior a
# CMD /usr/bin/supervisord -c ./docker/supervisiord.conf
