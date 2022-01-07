FROM golang:1.16-alpine as build

WORKDIR /app

COPY . ./

RUN go mod download

RUN cd cmd && go build -o ../dist/poi


# Finally clean container with only the service included to reduce size
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/dist service

CMD [ "service/poi" ]