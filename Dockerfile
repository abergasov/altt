FROM golang:1.20 AS build
RUN echo "Based on commit: $GIT_HASH"

# All these steps will be cached
RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
# COPY the source code as the last step
COPY . .

RUN CGO_ENABLED=0 make build

# step 2 - create container to run
FROM debian:buster-slim
WORKDIR /app
COPY --from=build /app /app
COPY --from=build /app/bin/binary /app/
COPY --from=build /app/configs/* /app/configs/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
RUN chmod +x /app/binary
EXPOSE 8000/tcp
CMD /app/binary