FROM golang:1.22-bullseye AS builder

WORKDIR /confik

COPY go.mod go.sum ./

RUN --mount=type=cache,target="/root/.cache/go-build" go mod download

COPY cmd ./cmd
COPY pkg ./pkg

ENV CGO_ENABLED=0

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o ./confik ./cmd/confik

RUN chmod +x ./confik

FROM scratch

WORKDIR /

COPY --from=builder /confik/confik /confik

ENTRYPOINT ["/confik"]