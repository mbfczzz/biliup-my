# ============ 前端构建 ============
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# ============ Go 构建 ============
FROM golang:1.21-alpine AS backend
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -o biliup ./cmd/server/

# ============ 运行 ============
FROM alpine:3.19
RUN apk add --no-cache ffmpeg ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=backend /app/biliup .
COPY --from=frontend /app/frontend/dist ./frontend/dist

RUN mkdir -p /app/downloads

EXPOSE 19159

ENTRYPOINT ["./biliup"]
CMD ["--port", "19159", "--db-host", "mysql", "--db-port", "3306", "--db-user", "biliup", "--db-pass", "biliup123", "--db-name", "biliup"]
