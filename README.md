# Poi Service

[![Build Status](https://github.com/OneOfTheseDays83/poiService/workflows/test%20and%20build/badge.svg)](https://github.com/OneOfTheseDays83/poiService/actions?workflow=test%20and%20build)
[![Coverage Status](https://github.com/OneOfTheseDays83/poiService/badge.svg?branch=master)](https://github.com/OneOfTheseDays83/poiService?branch=master)

This service allows users to store POIs at a given location. Detailed functionality:
* Create a new POI in the service.
* Retrieve a specific POI from the service.
* Update an existing POI in the service.
* Delete a POI from the service.
* List all POIs available in the service.
   * It is possible to find all POIs within a given radius.
   * It is possible to get all POIs

The service can only be used with a valid JWT. This must be retrieved by the client and provided with each request.

## Architecture
The application is set up as a microservice with a REST API (c.f. [Architecture.md](./doc/Architecture.md))

## Build
If you didn't change anything you don't need to build the application and can rather jump to "Start the service".

### Get dependencies
Download the dependencies first. This will download all the needed go modules and
the docker container (mongodb) that is used for data storage. In a real system a cloud managed
database service should be used.
```shell
make download-deps 
```

### Build
```shell
make build
```

## Use the poi service
### Start dependencies
In order to run the poi service oauth server and the mongodb is needed.
This needs only to be done once during system lifecycle or after stopped them.
```shell
make start-environment
```

### Start the poi service
```shell
make start 
```

### Requests
#### Get credentials
```shell
curl -s -k -X POST -H "Content-Type: application/x-www-form-urlencoded" -d grant_type=client_credentials -u 'my-client:secret' http://localhost:4444/oauth2/token
```

Store the retrieved JWT (e.g. export TOKEN=ey...)"

#### Create Poi
Replace the bearer token by the one you got from the enrollment status response!
```shell
curl -v -X POST http://localhost:8000/v1/pois -H "Authorization: Bearer "$TOKEN --data  '{"name" : "Dresden", "longitude" : 13.737262, "latitude" : 51.050407}'
curl -v -X POST http://localhost:8000/v1/pois -H "Authorization: Bearer "$TOKEN --data  '{"name" : "Wolfsburg", "longitude" : 10.780420, "latitude" : 52.427547}'
curl -v -X POST http://localhost:8000/v1/pois -H "Authorization: Bearer "$TOKEN --data  '{"name" : "Berlin", "longitude" : 13.404954, "latitude" : 52.520008}'
curl -v -X POST http://localhost:8000/v1/pois -H "Authorization: Bearer "$TOKEN --data  '{"name" : "Munich", "longitude" : 11.576124, "latitude" : 48.137154}'
```
If the creation was successful it is responded with a http 200 and a corresponding unique poi id (e.g. "3cba9846-aeea-4c2e-9f24-38289ef2b926").
This unique POI Id must be used for GET, UPDATE and DELETE requests.

#### Get Poi
Replace the id behind v1/pois/ to the one you got from the creation response.
Replace the bearer token by the one you got from the enrollment status response!
```shell
curl -v -X GET http://localhost:8000/v1/pois/3cba9846-aeea-4c2e-9f24-38289ef2b926 -H "Authorization: Bearer "$TOKEN
```

#### Update Poi
Replace the id behind v1/pois/ to the one you got from the creation response.
Replace the bearer token by the one you got from the enrollment status response!
```shell
curl -v -X PUT http://localhost:8000/v1/pois/3cba9846-aeea-4c2e-9f24-38289ef2b926 -H "Authorization: Bearer "$TOKEN --data  '{"name" : "Dresden Centre" "longitude" : 13.737262, "latitude" : 51.050407}'
```

#### Delete Poi
Replace the id behind v1/pois/ to the one you got from the creation response.
Replace the bearer token by the one you got from the enrollment status response!
```shell
curl -v -X DELETE http://localhost:8000/v1/pois/3cba9846-aeea-4c2e-9f24-38289ef2b926 -H "Authorization: Bearer "$TOKEN
```

#### Search Poi by a radius
Replace the id behind v1/pois/ to the one you got from the creation response.
Replace the bearer token by the one you got from the enrollment status response!
The radius is in meter!
```shell
curl -v -X POST http://localhost:8000/v1/pois/list -H "Authorization: Bearer "$TOKEN --data '{"longitude" : 13.737262, "latitude" : 51.050407, "radius" : 20000}'
```

#### Search all Poi
Use the same curl request but remove the data part!

## Open points
* OpenApi spec missing in ./api
* The api should be improved to use protobuf and not JSON
* Listing all Poi request should have paging concept since this data can get very huge
* Unit testing must be extended
* Integration tests must be implemented