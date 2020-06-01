FROM golang as builder
WORKDIR /build
COPY . .
RUN go mod download && go build

### Multi-stage build stuff, ignore for now
# FROM scratch
# COPY --from builder /build/hackaton-2020 /bin/hackaton-2020

CMD ["/build/hackaton-2020"]
