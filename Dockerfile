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

# Install sqlite driver deps
RUN apk update && apk add --no-cache sqlite

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Bring backend binary
COPY --from=backend-build /app/server .

# Bring frontend build output
# COPY --from=frontend-build /app/frontend/dist ./frontend

# Bring db directory (including empty app.sql)
COPY db/ ./db

# Set permissions so non-root user can read/write the DB
RUN chown -R appuser:appgroup /app

# Switch from root → appuser
USER appuser

EXPOSE 8080
CMD ["./server"]
