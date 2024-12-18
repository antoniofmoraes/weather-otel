# Labs Go Expert FullCycle
## Weather API - Integração open telemetry e zapkin

### Configuração

Copiar o arquivo `.env.exemplo` com o nome `.env`

### Rodando localmente

- Rode o docker-compose: `docker compose up --build`

### Zapkin

O zapkin estará rodando em http://localhost:9411

### Exemplo de requisição

```bash
curl --location 'http://localhost:8080/weather' \
      --header 'Content-Type: application/json' \
      --data '{
          "cep": "81900550"
      }'
```