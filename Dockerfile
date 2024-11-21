FROM golang:1.22-alpine AS gobuild

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build
ADD . /build

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ./nautilus ./cmd/nautilus/main.go

FROM alpine:3.20.3

ARG TARGETOS
ARG TARGETARCH
# Install required packages
RUN apk add --no-cache bash ca-certificates
# Manually update the local certificates
RUN update-ca-certificates
# Copy executable from build
COPY --from=gobuild /build/nautilus /nautilus
RUN chmod +x /nautilus
# Set entrypoint to run executable
ENTRYPOINT [ "/nautilus" ]
CMD [ "agent" ]