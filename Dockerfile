# Building the react ui

FROM node:lts AS react-ui-build

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR /app
RUN git clone https://github.com/Blitz-Cloud/ettiUi.git /app
RUN npm i -g pnpm && pnpm i
RUN pnpm run build




# Building the binary of the App
FROM golang:1.23.3-alpine AS build

# Install SQLite development libraries for Cgo in the build stage
# 'sqlite-dev' is the package name for Alpine
RUN apk add build-base
RUN apk add --no-cache sqlite-dev
# setting the workdir
WORKDIR /go/src/ettiHelper

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

ARG TARGETARCH
# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=1 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o app .


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest AS release

WORKDIR /app

# Create the `public` dir and copy all the assets into it
RUN mkdir ./views
COPY ./views ./views
COPY ./static ./static
COPY ./ettiContent.db ./ettiContent.db
COPY --from=react-ui-build /app/build ./build

#  should be replaced here as well
COPY --from=build /go/src/ettiHelper/app .

# Add packages
RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/app

# Exposes port 3000 because our program listens on that port
EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "/app/app"]