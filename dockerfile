FROM --platform=linux/amd64 golang:1.23.1-alpine AS builder


# Creates an app directory to hold your app’s source code
WORKDIR /build
 
# Copies everything from your root directory into /app
COPY . .
 
# Installs Go dependencies
RUN go get
 
# Builds your app with optional configuration
RUN go build -o /messageapi
 
# Tells Docker which network port your container listens on
EXPOSE 9092
 
# Specifies the executable command that runs when the container starts
# CMD [ "/messageapi" ]
ENTRYPOINT [ "/messageapi" ]