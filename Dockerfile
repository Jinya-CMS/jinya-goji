ARG CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX

FROM $CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX/library/alpine:latest

COPY jinya-cms-goji /app/jinya-cms-goji

RUN apk add vips

WORKDIR /app

CMD ["/app/jinya-cms-goji"]
