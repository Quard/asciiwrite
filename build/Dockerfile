FROM golang:latest as build

WORKDIR /app
COPY . .

ENV GO111MODULE on
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -v -o /asciiwrite cmd/asciiwrite/main.go
RUN chmod +x /asciiwrite



# FROM gcr.io/distroless/base
FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /asciiwrite /asciiwrite
COPY configs/credentials/asciiwrite-firebase.json /asciiwrite-firebase.json
ENV GOOGLE_APPLICATION_CREDENTIALS=/asciiwrite-firebase.json
ARG AUTHTOKEN
CMD /asciiwrite run -t $AUTHTOKEN