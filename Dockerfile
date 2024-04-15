FROM golang as builder
WORKDIR /src/service
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=`git tag --sort=-version:refname | head -n 1`" -o rms-health -a -installsuffix cgo rms-health.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata && apk add --update docker openrc
RUN mkdir /app
WORKDIR /app
COPY --from=builder /src/service/rms-health .
COPY --from=builder /src/service/configs/rms-health.json /etc/rms/
CMD ["./rms-health"]