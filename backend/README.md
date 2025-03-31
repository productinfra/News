# news Backend Code

## Docker Deployment Commands
**dockerfile**
```dockerfile
FROM golang:alpine AS builder

# Set necessary environment variables for our image
ENV GO111MODULE=on \
CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64

# Move to the working directory: /build
WORKDIR /build

# Copy the go.mod and go.sum files from the project and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Compile the code into a binary executable file news
RUN go build -o news .

###################
# Now create a smaller image
###################
FROM debian:stretch-slim
#FROM scratch

COPY ./wait-for.sh /
COPY ./templates /templates
COPY ./static /static
COPY ./conf /conf

# Copy the executable file from the builder image to the current directory
COPY --from=builder /build/news /

#RUN set -eux \
#
#    && apt-get update \
#    && apt-get install -y --no-install-recommends netcat \
#    && chmod 755 wait-for.sh

# Declare the server port
EXPOSE 8081

# Command to run the service
ENTRYPOINT ["/news", "conf/config.yaml"]

```

> docker build -t news .


> docker run -d -p 8081:8081 --name news news


> docker logs -f news