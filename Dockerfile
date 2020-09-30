FROM alpine:latest

COPY ./bin /usr/local/bin/

EXPOSE 8000

CMD ["gh-fetch"]