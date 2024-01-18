FROM golang:1.21-alpine as build-server
WORKDIR /src
COPY backend/go.mod backend/go.sum .
RUN go mod download
COPY backend .
#RUN go build -o /bin/import ./cmd/import
#RUN go build -o /bin/learn ./cmd/learn
#RUN go build -o /bin/query ./cmd/query
#RUN go build -o /bin/update ./cmd/update
RUN go build -o /bin/server ./cmd/server

FROM node:19-bullseye as build-frontend
WORKDIR /src
COPY frontend/package.json frontend/yarn.lock .
RUN yarn install
COPY frontend .
RUN yarnpkg run build

FROM alpine:latest
#RUN apk update
#RUN apk upgrade
#RUN apk add --no-cache sqlite
COPY --from=build-frontend /src/build /app/frontend
COPY --from=build-server /bin /app/bin
EXPOSE 9090
ENV ACCOUNTS_DBFILE /app/data/xact.db
ENV ACCOUNTSERVER_FRONTENDPATH /app/frontend
ENV ACCOUNTSERVER_PORT 9090
CMD ["/app/bin/server"]
