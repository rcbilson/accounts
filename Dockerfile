FROM golang:1.21 AS build-server
WORKDIR /src
COPY backend/go.mod backend/go.sum .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    go mod download && go mod verify
COPY backend .
# sqlite requires cgo
ARG CGO_ENABLED=1
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    go build -o /bin/server ./cmd/server

FROM node:19-bullseye AS build-frontend
WORKDIR /src
COPY frontend/package.json frontend/yarn.lock .
RUN yarn config set network-timeout 300000
RUN yarn install
COPY frontend .
RUN yarnpkg run build

FROM golang:1.21
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
