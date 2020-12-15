FROM golang:alpine

WORKDIR /go/src/app
COPY . .

# alpine doesn't include GCC or G++ and we need them for CGO
RUN apk --update upgrade && \
    apk add gcc && \
    apk add g++ && \
    apk add sqlite && \
    rm -rf /var/cache/apk/*

RUN sqlite3 resolver.db < internal/sqlite/migrations/00_init.sql

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["resolver"]
