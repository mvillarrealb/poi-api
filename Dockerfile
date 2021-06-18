FROM golang:1.15 as builder
RUN mkdir -p /poi-api/api
WORKDIR /poi-api
ADD api ./api
COPY go.mod go.sum main.go ./
RUN go build -ldflags "-linkmode external -extldflags -static" -o main .

FROM scratch
ENV DATABASE_HOST="127.0.0.1" \
DATABASE_PORT="5432" \
DATABASE_USERNAME="postgres" \
DATABASE_PASSWORD="password"

COPY --from=builder /poi-api/main ./main
CMD ["./main"]