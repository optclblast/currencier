# Currencier API
## Setup
1. Get an API token from [fxratesapi](https://fxratesapi.com) (it is free)
2. Install docker on your machine 
3. Set API token from [your fxrateapi account](https://fxratesapi.com/app/tokens) into .env file
```
FXRATEAPI_API_TOKEN=<your-token-here>
```
4. Build and run!
```bash
sudo docker compose build && sudo docker compose up -d
```
Or
``` bash
make up
```
## Tests
``` bash
sudo docker compose -f docker-compose.test.yaml up --build --abort-on-container-exit
sudo docker compose -f docker-compose.test.yaml down --volumes
```
Or 
``` bash
make test
```
## Usage 
**Endpoints:**  
[GET] /currency  
Query params:
* date - **required**. format 31.01.2001 (dd.mm.yyyy)
* val - **required**. currency [ISO 4217](https://en.wikipedia.org/wiki/ISO_4217#Active_codes_(list_one)) char code
* val_to - **optional**. Default: RUB

Example:   
``` bash
curl --request GET \
  --url 'http://localhost:8080/currency?val=EUR&date=03.05.2024'
```

Response:
```json
{
  "date": "2024-01-05T00:00:00Z",
  "currency_id": "EUR",
  "currency_id_to": "RUB",
  "value": 99.493813
}
```