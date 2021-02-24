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

FROM alpine
COPY --from=build /mural-device /mural-device
# In order to persist info from container we use this dir and bind to a host volume
# docker run -v <Path to host dir>:/containerFiles -p <host port to use>:42069 -it my-go-app
RUN mkdir /containerFiles
# we need to add a directory which links to our fs folder to contain mural software info
# ## Our start command which kicks off
# ## our newly created binary executable
ENTRYPOINT [ "/mural-device" ]
