FROM alpine:latest
# FROM scratch
LABEL developer="aserdipanda@gmail.com"
# EXPOSE 8080
COPY ./app ./
COPY ./.env ./

RUN chmod -R 777 /app
CMD ["/app"]