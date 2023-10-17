# Match version of node with .nvmrc in root dir
FROM node:lts-hydrogen

WORKDIR /usr/src/app

COPY cmd/js-client/package.json cmd/js-client/package-lock.json ./

RUN npm ci

# Copy the rest of the application code
COPY cmd/js-client/. .
