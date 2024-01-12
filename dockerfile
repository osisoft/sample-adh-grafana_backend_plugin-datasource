# Build Stage 1
# Compile frontend component
#
FROM node:alpine3.19 AS appbuild1

WORKDIR /usr/src/app
COPY . .
RUN npm install
RUN npm run build
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

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
FROM grafana/grafana:10.0.3

ENV GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=aveva-adh-datasource
WORKDIR /var/lib/grafana/plugins/aveva-adh-datasource
COPY --from=appbuild1 /usr/src/app/dist ./dist
COPY --from=appbuild2 /usr/src/app/dist ./dist
COPY package.json .