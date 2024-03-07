###############################################################################
# Builder
###############################################################################
# This image has a debian 12 base
FROM golang:1.22.1-bookworm AS builder

# create a 'new' dir to make sure there are no file name collisions
WORKDIR /build

# separating out these files to keep deps cached
COPY go.mod go.sum ./

RUN go mod download

# Copy the rest of the files
COPY . .

RUN GOOS=linux go build

###############################################################################
# Actual image
###############################################################################
# The base of the builder image is also debian 12
FROM gcr.io/distroless/base-debian12

# Copy the built executable from the builder stage
COPY --from=builder /build/bits-and-bytes-bot /

CMD ["/bits-and-bytes-bot"]
