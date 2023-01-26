FROM --platform=$BUILDPLATFORM golang:1.19 as build
WORKDIR /src
COPY go.mod go.sum .
RUN go mod download
ARG TARGETOS TARGETARCH
ARG VERSION=development
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-X main.version=$VERSION" -o /out/aurora

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=build /out/aurora /app/aurora
WORKDIR /app
ENV PATH="/app"
CMD ["aurora"]
