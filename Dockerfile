FROM golang:1.14-alpine as deps

ADD go.mod /app/go.mod
WORKDIR /app

RUN ["go", "mod", "download", "-x"]

FROM deps as build

ADD . /app
RUN go build -o "/bin/deploy-lander" ./cmd/server && \
    chmod a+x /bin/deploy-lander

FROM alpine

COPY --from=build /bin/deploy-lander /deploy-lander

RUN addgroup -g 9999 -S user && \
    adduser -u 9999 -G user -S -H user

USER user
EXPOSE 8080
CMD ["/deploy-lander"]
