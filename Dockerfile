# Step 1 - build app
FROM golang:alpine AS build

# Step 1.1 install build-base
RUN apk add build-base
WORKDIR /src
COPY . .

# Step 1.2 install sqlite3
RUN apk add git
RUN go get github.com/mattn/go-sqlite3

# Step 1.3 build application
ENV CGO_ENABLED=1
RUN go build  -ldflags "-linkmode external -extldflags -static" -o challenge

# Step 2 - prepare image
FROM alpine
WORKDIR /app
COPY --from=build /src /app

# Step 2.2 expose ports
EXPOSE 80
EXPOSE 443
ENTRYPOINT ["/app/challenge"]
