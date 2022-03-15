FROM alpine:latest AS RUN
WORKDIR /GoPracticum/app
COPY ./shortener .
CMD ["./shortener"]
