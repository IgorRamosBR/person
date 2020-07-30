FROM golang as builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . .

RUN go mod download
RUN go build cmd/main.go

FROM scratch
COPY --from=builder /app/ /app/

ENTRYPOINT ["/app/main"]