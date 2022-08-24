FROM golang:1.19-alpine3.16 as builder

WORKDIR /goman/
RUN apk --no-cache add build-base git
ADD . .

ENV GOBIN=/goman/bin/
ENV GOOS=linux
ENV GOARCH=amd64

RUN go mod download
RUN go install -a -trimpath -ldflags='-w -s' ./cmd/...


FROM alpine:3.16
WORKDIR /apps/
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /goman/bin/* /apps/
ENV PATH="/apps/:${PATH}"
