FROM golang:alpine AS builder

# Set necessary environment variables for our image
ENV GO111MODULE=on \
CGO_ENABLED=0 \
GOOS=linux

# Move to the working directory: /build
WORKDIR /build

# Copy the go.mod and go.sum files from the project and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Compile the code into a binary executable file called news
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

# Declare the service port
EXPOSE 8081

# The command to run
ENTRYPOINT ["/news", "conf/config.yaml"]
