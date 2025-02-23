FROM golang:1-alpine AS build

WORKDIR /app

COPY src .

RUN go build -o /app/dumbledore

FROM golang:1-alpine

WORKDIR /app

RUN adduser -D dumbledore

COPY --from=build /app/dumbledore /app/dumbledore

RUN chown dumbledore:dumbledore /app/dumbledore

USER dumbledore

ENV GIN_MODE=release

ENTRYPOINT ["./dumbledore"]