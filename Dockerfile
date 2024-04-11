FROM golang:1.22.1



## ---------- ARGS
ARG TARGET_API
ARG API_PORT



## ---------- ENVS
ENV TARGET_API=${TARGET_API}
ENV API_PORT=${API_PORT}



## ---------- BUILD
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build main.go && \
    chmod +x main



## ---------- MAIN
# ENTRYPOINT [ "./entrypoint.sh", "${TARGET_API}" ]
ENTRYPOINT [ "./main" ]
HEALTHCHECK \
    --interval=30s \
    --timeout=30s \
    --start-period=5s \
    --retries=3 \
    CMD [ "bash", "-c", "curl -s -w '%{http_code}' -o /dev/null localhost:${API_PORT}" ]
