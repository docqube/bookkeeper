FROM alpine:20230329

ENV BOOKKEEPER_DATABASE_HOST="postgres"
ENV BOOKKEEPER_DATABASE_PORT="5432"
ENV BOOKKEEPER_DATABASE_USER="postgres"
ENV BOOKKEEPER_DATABASE_PASSWORD="secr3t"
ENV BOOKKEEPER_DATABASE_DATABASE="bookkeeper"

ENV BOOKKEEPER_PORT="8080"
ENV GIN_MODE="release"

RUN apk add -U --no-cache bash ca-certificates openssl

EXPOSE 8080

COPY build/bookkeeper-api bookkeeper-api
COPY .sql/ .sql/
RUN chmod +x bookkeeper-api

CMD ["/bin/bash", "-c", "./bookkeeper-api"]
