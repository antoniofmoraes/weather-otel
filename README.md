# Labs Go Expert FullCycle
# Weather API - Deploy Cloud Run

## Configuração

Copiar o arquivo `.env.exemplo` com o nome `.env`

## Rodando localmente

- Faça o build da imagem `docker build . -t weather_app`
- Rode a imagem `docker run -p 8080 weather_app`

## Google Cloud Run

O projeto foi publicado usando Google Cloud Run.
https://weather-app-361381993255.us-central1.run.app/clima/{CEP}