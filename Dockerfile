# stage de build
FROM golang:1.22.0 AS build

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o rinha cmd/main.go

# stage imagem final
FROM scratch

WORKDIR /app

COPY --from=build /app/rinha ./

EXPOSE 8000

CMD [ "./rinha" ]