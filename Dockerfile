ARG TOOLCHAIN_VERSION=latest
# Builder image with toolchain - can be overridden with --build-arg
FROM docker.io/library/golang:${TOOLCHAIN_VERSION} AS builder

# Set the working directory and copy the source code
WORKDIR /cosi
COPY . .

# Build the binary
RUN go build -o /cosi/bin/sample-cosi-driver ./cmd/sample-cosi-driver

# Runtime image
FROM docker.io/rockylinux/rockylinux:9-ubi-micro

COPY --from=builder /cosi/bin/sample-cosi-driver /usr/local/bin/sample-cosi-driver

# Set the working directory
WORKDIR /cosi

# Create a non-root user
RUN echo "cosi:*:1001:cosi" >> /etc/group && \
    echo "cosi:*:1001:1001::/cosi:/bin/false" >> /etc/passwd

# Set permissions on the binary
RUN chown 1001:1001 /usr/local/bin/sample-cosi-driver && \
    chmod 0755 /usr/local/bin/sample-cosi-driver && \
    mkdir -p /var/lib/cosi && \
    chown -R 1001:1001 /var/lib/cosi

# Run as non-root
USER cosi

# Expose the default port - port 80 for S3, 443 for S3 with TLS
EXPOSE 80
EXPOSE 443

# set volume mount point for cosi socket
VOLUME [ "/var/lib/cosi" ]

# Disable healthcheck
HEALTHCHECK NONE

# Metadata params
LABEL maintainer="Mateusz Urbanek <mateusz.urbanek.98@gmail.com>"

ENTRYPOINT [ "/usr/local/bin/sample-cosi-driver" ]
CMD []
