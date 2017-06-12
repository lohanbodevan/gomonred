# GOMONRED
Bootstrap project to run HTTP REST API on GO lang with MongoDB and Redis

## API Endpoints
GET /cars
```
[
  {
    "name": "Golf",
    "brand": "VW"
  },
  {
    "name": "Polo",
    "brand": "VW"
  },
  {
    "name": "Uno",
    "brand": "Fiat"
  }
]
```

POST /cars
```
{
  "name": "Lancer Evolution",
  "brand": "Mitsubishi"
}
```

## Requirements
* Docker Compose

## RUN
```
docker-compose up
```

http://localhost:8080/cars
