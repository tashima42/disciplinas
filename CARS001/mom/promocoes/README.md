# Promocoes

## RabbitMQ

```sh
docker run \
  -d \
  -p 5672:5672 \
  --env RABBITMQ_DEFAULT_USER=user \
  --env RABBITMQ_DEFAULT_PASS=password \
  --name rabbit \
  rabbitmq:4.2
```

## Build

```sh
make build
```

## Create keys

```sh
./dist/promocoes crypto -n gateway -p gateway
./dist/promocoes crypto -n ranking -p ranking
./dist/promocoes crypto -n promocao -p promocao
```

## Run services

Run each one in a separate terminal

```sh
./dist/promocoes gateway
./dist/promocoes ranking
./dist/promocoes promocao
./dist/promocoes notificacao
./dist/promocoes consumer --categores=category1,category2
```

## Interact with gateway

The gateway has a terminal user interface with a menu and multiple options
