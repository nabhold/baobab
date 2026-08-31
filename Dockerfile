FROM golang:1.27.0-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/controlplane ./cmd/controlplane
FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/controlplane /controlplane
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/controlplane"]
