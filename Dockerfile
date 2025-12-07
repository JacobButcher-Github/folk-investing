# ---------- FRONTEND BUILD ----------
# FROM node:20 AS frontend-build
# WORKDIR /app/frontend
# COPY frontend/ .
# RUN npm install && npm run build

# ---------- BACKEND BUILD ----------
FROM golang:1.22 AS backend-build
WORKDIR /app
COPY backend/ .
RUN go mod tidy && CGO_ENABLED=1 go build -o server ./main.go

# ---------- FINAL IMAGE ----------
FROM alpine:3.19

RUN apk update && apk add --no-cache sqlite

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Backend binary
COPY --from=backend-build /app/server .

# frontend build
# COPY --from=frontend-build /app/frontend/dist ./frontend

COPY db/ ./db

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080
CMD ["./server"]
