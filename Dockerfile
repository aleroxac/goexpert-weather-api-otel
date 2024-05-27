FROM golang:1.22.3
# FROM golang:1.22.3-alpine3.20



## ---------- ARGS
ARG TARGET_API
ARG API_PORT



## ---------- ENVS
ENV TARGET_API=${TARGET_API}
ENV API_PORT=${API_PORT}



## ---------- BUILD
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build ${TARGET_API}/cmd/app/main.go && \
    chmod +x main



## ---------- MAIN
ENTRYPOINT [ "./main" ]
HEALTHCHECK \
    --interval=10s \
    --timeout=5s \
    --start-period=5s \
    --retries=3 \
    CMD  [ $(curl -s -o /dev/null -w "%{http_code}" http://localhost:${API_PORT}/status) -eq 200 ] || exit 1
