# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13

# Add Maintainer Info
LABEL maintainer="Kataru Onosato <konosato@idcf.jp>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Set up for multibyte-code
RUN localedef -v -c -i ja_JP -f UTF-8 ja_JP.UTF-8; echo "";
ENV LANG=ja_JP.UTF-8

# Set up for time-zone
RUN rm -f /etc/localtime
RUN ln -fs /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# Copy go mod and sum files
COPY go.mod go.sum install-realize.sh sqlboiler.toml ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
#RUN go build -o main .

# Expose port 1323 to the outside world
#EXPOSE 1323

# Command to run the executable
#CMD ["./main.go"]
