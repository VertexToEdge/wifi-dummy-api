FROM golang:1.21 as builder
WORKDIR /src
COPY . .
RUN go build -o /bin/app ./main.go

FROM scratch
COPY --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]