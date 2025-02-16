FROM golang:1.23

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod tidy

COPY . .

COPY .env .

RUN go build -o ecommerce-products

RUN chmod +x ecommerce-products

EXPOSE 9001

CMD [ "./ecommerce-products" ]