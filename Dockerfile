FROM golang:1.22-alpine AS gobuild

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build
ADD . /build

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w -extldflags "-static"' -o ./nautilus ./cmd/nautilus/main.go

RUN chmod +x ./nautilus

FROM alpine:3.20.3

ARG TARGETOS
ARG TARGETARCH
# Install required packages and manually update the local certificates
RUN apk add --no-cache bash ca-certificates && update-ca-certificates
# Copy executable from build
COPY --from=gobuild /build/nautilus /nautilus
# Expose port 8080 by default
EXPOSE 8080
# Set entrypoint to run executable
ENTRYPOINT [ "/nautilus" ]
CMD [ "agent" ]