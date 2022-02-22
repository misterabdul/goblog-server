FROM alpine:3.15.0

WORKDIR /app/

COPY ./build/output/. /app/

ENV APP_NAME=goblog \
    APP_HOST=localhost \
    APP_PORT=80 \
    APP_MODE=release \
    COOKIE_DOMAIN=.localhost \
    COOKIE_SECURE=false \
    DB_PROTOCOL=mongodb \
    DB_HOST=localhost \
    DB_PORT=27017 \
    DB_USER=goblog \
    DB_PASS=12345678 \
    DB_NAME=goblog \
    AUTH_SECRET=secret \
    AUTH_ACCESS_DURATION=60 \
    AUTH_REFRESH_DURATION=14

ENTRYPOINT [ "/app/server" ]
