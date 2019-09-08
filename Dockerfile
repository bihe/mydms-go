## fronted build-phase
## --------------------------------------------------------------------------
FROM node:lts-alpine AS FRONTEND-BUILD
WORKDIR /frontend-build
COPY ./frontend.angular .
RUN npm install -g @angular/cli@latest && npm install && npm run build -- --prod --base-href /ui/
## --------------------------------------------------------------------------

## backend build-phase
## --------------------------------------------------------------------------
FROM golang:alpine AS BACKEND-BUILD

ARG buildtime_variable_version=2.0.0
ARG buildtime_variable_timestamp=YYYYMMDD
ARG buildtime_varialbe_commit=b75038e5e9924b67db7bbf3b1147a8e3512b2acb

ENV VERSION=${buildtime_variable_version}
ENV BUILD=${buildtime_variable_timestamp}
ENV COMMIT=${buildtime_varialbe_commit}

WORKDIR /backend-build
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.Version=${VERSION}-${COMMIT} -X main.Build=${BUILD}" -tags prod -o mydms.api
#COPY --from=FRONTEND-BUILD /frontend-build/dist  ./ui
## --------------------------------------------------------------------------

## runtime
## --------------------------------------------------------------------------
FROM alpine:latest
LABEL author="henrik@binggl.net"
WORKDIR /opt/mydms
RUN mkdir -p /opt/mydms/uploads && mkdir -p /opt/mydms/ui && mkdir -p /opt/mydms/etc && mkdir -p /opt/mydms/logs
COPY --from=BACKEND-BUILD /backend-build/mydms.api /opt/mydms
COPY --from=FRONTEND-BUILD /frontend-build/dist  /opt/mydms/ui
RUN ls -l /opt/mydms
RUN ls -l /opt/mydms/etc
RUN ls -l /opt/mydms/ui

EXPOSE 3000

CMD ["/opt/mydms/mydms.api","--c=/opt/mydms/etc/application.json","--port=3000", "--hostname=0.0.0.0"]
