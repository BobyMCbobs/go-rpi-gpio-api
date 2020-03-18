FROM arm32v7/golang:1.13.4-alpine3.10 AS api
WORKDIR /app
COPY . .
RUN addgroup -g 997 gpio
RUN adduser -G gpio -D -H -h /app user
RUN GOARCH=arm CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -o gpio-api src/main.go

FROM scratch
WORKDIR /app
ENV PATH=/app
COPY --from=api /app/gpio-api .
COPY --from=api /etc/passwd /etc/passwd
COPY --from=api /etc/group /etc/group
EXPOSE 8080
USER user
ENTRYPOINT ["/app/gpio-api"]
