FROM alpine
RUN apk update && apk add tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Taipei /etc/localtime
RUN echo "Asia/Taipei" > /etc/timezone
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY ./slackbot-go /app/slackbot-go
COPY ./plugins /app/plugins
RUN chmod 777 slackbot-go
CMD ["./slackbot-go"]

EXPOSE 8001