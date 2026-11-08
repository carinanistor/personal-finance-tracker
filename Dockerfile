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

LABEL org.opencontainers.image.title="Personal Finance Tracker" \
      org.opencontainers.image.authors="Nistor Carina <carina.nistor@student.upt.ro>" \
      org.opencontainers.image.source="https://github.com/carinanistor/personal-finance-tracker"

RUN adduser -D -u 10001 app

COPY --from=build /app /usr/local/bin/app
USER app

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/app"]
