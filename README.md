# Currencier API
## Setup
1. Get an API token from [fxratesapi](https://fxratesapi.com) (it is free, only your github account needen, or email)
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
Set CURRENCIER_TESTS_HOST variable in your .env file, so the file should contain somethink like this:
```
 FXRATEAPI_API_TOKEN=fxr_live_sd6f575da657g6df786s56765587
 CURRENCIER_TESTS_HOST="http://currencier:8080"
```
http://currencier:8080 for testing with .test compose project.  

And now run the tests!
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
## But why not cbr.ru?
I think cbr.ru is very outdated. It's actually difficult to get it to work correctly with the Go client. Plus, there are only a few currencies available! What a mess.
[fxratesapi](https://fxratesapi.com) is free and works very well.