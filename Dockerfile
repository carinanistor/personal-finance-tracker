FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

ARG COMMIT=dev
ARG BUILT_AT=unknown

RUN CGO_ENABLED=0 go build \
      -ldflags="-s -w \
        -X atad-project/internal/version.Commit=${COMMIT} \
        -X atad-project/internal/version.BuiltAt=${BUILT_AT}" \
      -o /app ./cmd/app

FROM alpine:3.21

LABEL org.opencontainers.image.title="ATAD project" \
      org.opencontainers.image.authors="Surname Firstname <firstname.surname@student.upt.ro>" \
      org.opencontainers.image.source="https://github.com/user/repo"

RUN adduser -D -u 10001 app

COPY --from=build /app /usr/local/bin/app
USER app

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/app"]
