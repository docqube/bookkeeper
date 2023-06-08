FROM alpine:20230329

ENV DB_HOST="postgres"
ENV DB_NAME="backkeeper"
ENV DB_PORT="5432"
ENV DB_USER="postgres"
ENV DB_PASSWORD="secr3t"

ENV PORT="8080"
ENV GIN_MODE="release"

RUN apk add -U --no-cache bash ca-certificates openssl

EXPOSE 8080

COPY build/bookkeeper-api bookkeeper-api
COPY .sql/ .sql/
RUN chmod +x bookkeeper-api

CMD ["/bin/bash", "-c", "./bookkeeper-api --db-host=${DB_HOST} --db-name=${DB_NAME} --db-user=${DB_USER} --db-password=${DB_PASSWORD}"]
