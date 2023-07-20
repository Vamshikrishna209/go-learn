# Build stage

FROM golang:1.20.5-alpine3.18 AS builder
WORKDIR /app
COPY . .
# Copy all files to the working directory of container image (.) from local machine or build context
RUN go build -o main main.go


# Run Stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
EXPOSE 8080

CMD [ "/app/main" ]