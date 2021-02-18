FROM golang:1.15.7-alpine3.13 as build

## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app
## We copy everything in the root directory
## into our /app directory
ADD . /app
## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /app

RUN go mod download
## we run go build to compile the binary
## executable of our Go program
RUN go build -o /mural-device

FROM alpine:3.11.3
COPY --from=build /mural-device /mural-device
# ## Our start command which kicks off
# ## our newly created binary executable
ENTRYPOINT [ "/mural-device" ]
