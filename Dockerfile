FROM golang:1.16-alpine AS compiler

WORKDIR /app/src

COPY . .

RUN go build -o sertifikadogrula

FROM alpine

COPY --from=compiler /app/src/sertifikadogrula /sertifikadogrula

ENTRYPOINT ["/sertifikadogrula"]
