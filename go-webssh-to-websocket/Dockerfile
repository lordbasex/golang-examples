##################################
# STEP 1 build executable binary
##################################
FROM golang:1.22.2-alpine3.19 AS builder
LABEL maintainer="Federico Pereira <lord.basex@gmail.com>"
ADD . /usr/src/webssh
WORKDIR /usr/src/webssh
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go/bin/webssh main.go

##################################
# STEP 2 build a small image
##################################
FROM alpine:3.19 as webssh
LABEL maintainer="Federico Pereira <lord.basex@gmail.com>"
RUN apk update && apk add --update-cache dumb-init tzdata
COPY --from=builder /go/bin/webssh /go/bin/webssh
COPY --from=builder /usr/src/webssh/public /go/bin/public
EXPOSE 8280
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
WORKDIR /go/bin
CMD ["/go/bin/webssh"]
