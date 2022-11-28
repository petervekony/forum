# syntax=docker/dockerfile:1
FROM golang:1.18-alpine as builder
LABEL description="Gritface forum project @grit:lab"
LABEL creators="oluwatosin, tvntvn, dicarste, adrian1206 and petrus_ambrosius"
RUN apk update && apk upgrade
RUN apk add --no-cache sqlite
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o gritface

# FROM golang:1.18-alpine
# WORKDIR /app
# COPY --from=builder ./gritface ./
# COPY log localhost.crt localhost.csr localhost.key readme.md ./
# COPY server/public_html ./
EXPOSE 443
CMD [ "./gritface" ]