FROM golang:1.22-alpine as base
RUN apk add --update --no-cache make git

ENV CGO_ENABLED=0

WORKDIR /root/

# TBD: this could be moved to the base image, but here is fully in control of developers
COPY Makefile .
RUN make tools-ci

COPY . .

### LINTING
FROM base as lint

CMD ["make", "lint"]

### TESTING
FROM base as test

CMD ["make", "test"]

### BUILDING
FROM base as builder

RUN make build

### PRODUCTION
FROM alpine:3
WORKDIR /root/


COPY --from=builder /root/bin/ .
EXPOSE 8080
CMD [ "./api-example" ]
