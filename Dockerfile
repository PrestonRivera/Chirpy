FROM golang:1.24.1 AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o Chirpy .

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache postgresql-client go
RUN go install github.com/pressly/goose/v3/cmd/goose@latest


COPY --from=build /app/Chirpy .
COPY --from=build /app/static ./static

ENV PLATFORM="dev"
ENV JWTSecret="aS+bzfD9A9LJtDOKjJfuOgueZSwPy6zhyZtaqzfO5Shf/dpYqlPvyYCESTcOTpSDvubaLfxeo8S9/hLrigVX4w=="
ENV POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
ENV DB_URL="postgres://postgres:postgres@db:5432/chirpy?sslmode=disable"
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=chirpy
ENV DB_HOST=db

EXPOSE 8080

CMD ["./Chirpy"]