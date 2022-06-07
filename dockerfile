# Build Stage 1
# This build created a staging docker image 
#
FROM node:alpine3.10 AS appbuild
RUN apk add g++ make python
WORKDIR /usr/src/app
COPY . .
RUN npm ci
RUN npm run build

# base grafana image
FROM grafana/grafana

ENV GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=aveva-sds-datasource
WORKDIR /var/lib/grafana/plugins/aveva-data-hub-sample 
COPY --from=appbuild /usr/src/app/dist ./dist
COPY package.json .
