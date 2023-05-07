# Build stage
FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/main .
COPY --from=build /app/env ./env/
COPY --from=build /app/db ./db/   
EXPOSE 8085
CMD ["./main"]