# kl-assignment

## run local
make run

## run docker
docker compose up -d or make run

## run migrate
docker compose up -d migrate or make migrate

## API Document

### account
Example curl
curl --location 'http://localhost:8080/api/account/5a5ae9dae1a111ef80ca0242ac180002' \
--header 'Content-Type: application/json'

###
Example curl
curl --location 'http://localhost:8080/api/user/000093b8e1a211ef95a30242ac180002' \
--header 'Content-Type: application/json'