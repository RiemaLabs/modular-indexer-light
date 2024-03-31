FROM golang:1.22-alpine as builder-go
WORKDIR /usr/src/light-indexer
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o light-indexer .


FROM node:lts-alpine as builder-node
WORKDIR /usr/src/dashboard
# It should be changed to git clone in the future
COPY ./modular-indexer-light-dashboard .
RUN  yarn install
RUN yarn build:prod


FROM nginx
WORKDIR /deploy
COPY --from=builder-go /usr/src/light-indexer/light-indexer .
COPY ./config.json .
COPY --from=builder-node  /usr/src/dashboard/build  ./html
COPY ./deploy/light-indexer.conf /etc/nginx/conf.d/default.conf

CMD ["sh", "-c", "nginx && ./light-indexer"]