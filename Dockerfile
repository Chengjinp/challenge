FROM golang:1.14-alpine AS build
RUN apk add build-base
WORKDIR /src
COPY . .


RUN apk add git
RUN go get github.com/mattn/go-sqlite3

RUN CGO_ENABLED=1 go build  -ldflags "-linkmode external -extldflags -static" -o /bin/challenge

FROM scratch
COPY --from=build /bin/challenge /bin/challenge


EXPOSE 8080
ENTRYPOINT ["/bin/challenge"]
