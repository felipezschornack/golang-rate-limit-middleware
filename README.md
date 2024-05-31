# Rate Limiter

## Objetivo

Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

## Configurações

As seguintes variáveis de ambiente podem ser configuradas:

- REDIS_HOSTNAME = redis
- REDIS_PORT = 6379
- ACCESS_TOKEN_RATE_LIMIT_IN_SECONDS = 10
- ACCESS_TOKEN_BLOCKING_WINDOW_IN_SECONDS = 10
- IP_RATE_LIMIT_IN_SECONDS = 10
- IP_BLOCKING_WINDOW_IN_SECONDS = 120
- WEB_SERVER_PORT = 8080

## Executando a aplicação utilizando Docker
Com o Docker instalado em sua estação de trabalho (https://www.docker.com/), execute o comando:
```
docker compose up -d --build rate-limiter redis
```

Isso instanciará a aplicação e o Redis.

A aplicação pode ser acessada no endereço:

http://localhost:8080/

Pode ser testada também via Curl:
```
curl --location 'http://localhost:8080/' \
    --header 'API_KEY: 999999'
```

Se a requisição tiver sucesso, o response será:
```
Status: 200 OK
Body: Request not blocked by rate limit!
```

Em caso de bloqueio, o response será:
```
Status: 429 Too Many Requests
Body: you have reached the maximum number of requests or actions allowed within a certain time frame
```

## Testes

Os testes podem ser feitos utilizando a ferramenta [Apache Benchmark](https://httpd.apache.org/docs/2.4/programs/ab.html).

Para tanto, execute o comando:
```
docker compose up -d ab
```

A seguir, execute o comando abaixo e descubra o "CONTAINER ID" da imagem "rate-limit-ab"
```
docker ps
```

De posse do id, execute o comando abaixo para entrar no container:
```
docker exec -it [container-id] bash
```

Dentro do container, o comando abaixo pode ser executado para testar a aplicação:

```
ab -v 1 -c 5 -n 10000 -H "API_KEY: b6674a9" http://rate-limiter:8080/
```
ou 
```
ab -v 1 -c 5 -n 10000 -H http://rate-limiter:8080/
```

onde:

- -v = Set verbosity level - 4 and above prints information on headers, 3 and above prints response codes (404, 200, etc.), 2 and above prints warnings and info.

- -c = Number of multiple requests to perform at a time.

- -n = Number of requests to perform for the benchmarking session.

- -H = Append extra headers to the request. The argument is typically in the form of a valid header line, containing a colon-separated field-value pair (i.e., "Accept-Encoding: zip/zop;8bit").