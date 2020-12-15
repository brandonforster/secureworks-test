FROM golang:alpine

WORKDIR /go/src/app
COPY . .

# alpine doesn't include GCC or G++ and we need them for CGO
RUN apk --update upgrade && \
    apk add gcc && \
    apk add g++ && \
    apk add sqlite && \
    rm -rf /var/cache/apk/*

# import our schema into the DB
RUN sqlite3 resolver.db < internal/sqlite/migrations/00_init.sql

# set the app's port
ENV PORT 8080

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["resolver"]
