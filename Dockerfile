FROM golang:1.21 as build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/wifi_api ./main.go

FROM scratch
COPY --from=build /bin/wifi_api /bin/app
CMD ["/bin/app"]