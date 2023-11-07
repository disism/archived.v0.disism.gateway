FROM ubuntu:latest

COPY ./app /app/

WORKDIR /app

RUN chmod +x ./app

EXPOSE 8032

CMD ["./app", "run"]