# Build Stage 1
# Compile frontend component
#
FROM node:alpine3.10 AS appbuild1
RUN apk add g++ make python
WORKDIR /usr/src/app
COPY . .
RUN npm ci
RUN npm run build

# Build Stage 2
# Compile backend component
#
FROM golang:alpine as appbuild2

RUN go install github.com/magefile/mage@latest

WORKDIR /usr/src/app
COPY . .
RUN go get -u github.com/grafana/grafana-plugin-sdk-go
RUN go mod tidy

RUN mage -v

# base grafana image
FROM grafana/grafana:8.3.0-beta2

ENV GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=aveva-sds-datasource
WORKDIR /var/lib/grafana/plugins/aveva-data-hub-sample 
COPY --from=appbuild1 /usr/src/app/dist ./dist
COPY --from=appbuild2 /usr/src/app/dist ./dist
COPY package.json .
