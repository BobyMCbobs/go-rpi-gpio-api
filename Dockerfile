FROM golang:1.13.4-alpine3.10 AS api
WORKDIR /app
COPY . .
RUN addgroup -g 997 gpio
RUN adduser -G gpio -D -H -h /app user
ARG AppBuildVersion="0.0.0"
ARG AppBuildHash="???"
ARG AppBuildDate="???"
ARG AppBuildMode="development"
ARG GOARCH=arm
RUN GOARCH="$GOARCH" CGO_ENABLED=0 GOOS=linux go build \
  -a \
  -installsuffix cgo \
  -ldflags "-extldflags '-static' -s -w \
    -X gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common.AppBuildVersion=$AppBuildVersion \
    -X gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common.AppBuildHash=$AppBuildHash \
    -X gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common.AppBuildDate=$AppBuildDate \
    -X gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common.AppBuildMode=$AppBuildMode" \
  -o gpio-api \
  src/main.go

FROM scratch
WORKDIR /app
ENV PATH=/app
COPY --from=api /app/gpio-api .
COPY --from=api /etc/passwd /etc/passwd
COPY --from=api /etc/group /etc/group
EXPOSE 8080
USER user
ENTRYPOINT ["/app/gpio-api"]
