FROM golang:1.19-alpine3.16 as build
WORKDIR /code

# Create a layer that contains the dependencies which most likely won't change often.
# This will allow changes in the code to not bust the dependencies cache, which
# will speed up the build.
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN  GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" src/main.go

# Use the scratch container which is empty as our final environment. This will reduce
# the amount of binaries, libraries, etc within the image, thus increasing security by
# decreasing the attach surface.
FROM scratch
USER 10001:10001
EXPOSE 8080
COPY --from=build /code/main /main
ENTRYPOINT [ "/main" ]
