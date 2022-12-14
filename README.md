# App Weight Management

- [Overview](#overview)
- [Run](#run)
    - [Database](#database)
    - [Application](#application)
- [API](#api)
    - [Weight Inputs](#weight-inputs)
        - [Create](#create)
        - [Update](#update)
        - [List](#list)
        - [Get](#get)
        - [Delete](#delete)

## Overview
Simple project to training and implementing simple concepts of Go Language. A simple API to manage weight evolution using as database a relational dbms (PostgreSQL). The content of the project will be development according the study of new concepts.

## Run

### Database
To run database container execute the follow command on terminal:
```sh
docker-compose up -d
```

### Application
To run the application execute the follow command in the root folder:
```sh
go run cmd/main.go
```

## API

### Weight Inputs
Weight inputs

|Field    |Type     |Description                  |
|---------|---------|-----------------------------|
|id       |int      |Unique id of input           |
|user_id  |int      |Id of user                   |
|weight   |float64  |Input weight                 |
|date     |string   |Input Date                   |


#### Create
Flow: create weight input
```mermaid
sequenceDiagram
Main ->> Handler: HandleCreateEntry(req, resp)
Handler ->> Handler: Decode Body
alt Error
    Handler -->> Main: HTTP Bad Request 400
else Success
    Handler ->> Service: Service.CreateEntry
        Service ->> Repository: Repository.Save
        Repository ->> Database: INSERT INTO t_entry
        Database -->> Repository: result, err
        Repository -->> Service: result, err
        Service -->> Handler: result, err
    alt err != nil
        Handler -->> Main: HTTP Internal Server Error 500
    else err == nil
        Handler ->> Handler: Encode JSON Response
        Handler -->> Main: HTTP Created 201
    end
end
```

Example of Request
```curlrc
curl -X POST \
  'localhost:8000/entries' \
  --header 'Content-Type: application/json' \
  --data-raw '{
	"user_id": 1,
	"weight": 115.5,
	"date": "2022-05-10 00:30:00"
}'
```

#### Update
Flow: update weight input
```mermaid
sequenceDiagram
Main ->> Handler: HandleUpdateEntry(req, resp)
Handler ->> Handler: Decode Body
alt Error
    Handler -->> Main: HTTP Bad Request 400
else Success
    Handler ->> Service: Service.UpdateEntry
        Service ->> Service: Service.GetEntry
        alt Error On Decode
            Service -->> Handler: Error
            Handler -->> Main: HTTP Bad Request 400
        else Error Entry Not Found
            Service -->> Handler: Entry Not Found
            Handler -->> Main: HTTP Not Found 404
        else Success
            Service ->> Repository: Repository.Update
            Repository ->> Database: UPDATE t_entry
            Database -->> Repository: result, err
            Repository -->> Service: result, err
            Service -->> Handler: result, err
        end
    alt err != nil
        Handler -->> Main: HTTP Internal Server Error 500
    else err == nil
        Handler ->> Handler: Encode JSON Response
        Handler -->> Main: HTTP Created 201
    end
end
```

Example of Request
```curlrc
curl -X PUT \
  'localhost:8000/entries/{id}' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "weight": 100.38,
  "date": "2022-01-01T10:00:00Z"
}'
```

Request Parameters
|Field    |Type     |Description               |
|---------|---------|--------------------------|
|id       |int      |ID of register to update  |

#### List
Flow: list weight inputs
```mermaid
sequenceDiagram
Main ->> Handler: HandleListEntries(req, resp)
Handler ->> Handler: Get Params
Handler ->> Service: Service.ListEntries
Service ->> Repository: Repository.ListAll
Repository ->> Database: SELECT FROM t_entry
Database -->> Repository: []results, err
Repository -->> Service: []results, err
Service -->> Handler: []results, err
alt err != nil
    Handler -->> Main: HTTP Internal Server Error 500
else err == nil
    Handler ->> Handler: Encode JSON Response
    Handler -->> Main: HTTP OK 200
end
```

Example of Request
```curlrc
curl -X GET \
  'localhost:8000/entries?start=0&count=20' \
```

Request Parameters
|Field    |Type     |Description                |
|---------|---------|---------------------------|
|start    |int      |Start index of search      |
|count    |int      |Quantity of registers      |

#### Get
Flow: get weight input by id
```mermaid
sequenceDiagram
Main ->> Handler: HandleGetEntry(req, resp)
Handler ->> Handler: Get Params
Handler ->> Service: Service.GetEntry
Service ->> Repository: Repository.FindById
Repository ->> Database: SELECT FROM t_entry
Database -->> Repository: result, err
Repository -->> Service: result, err
Service -->> Handler: result, err
alt err != nil
    Handler -->> Main: HTTP Internal Server Error 500
else entry == nil || entry.ID == 0
    Handler -->> Main: HTTP Not Found 404
else err == nil
    Handler ->> Handler: Encode JSON Response
    Handler -->> Main: HTTP OK 200
end
```

Example of Request
```curlrc
curl -X GET \
  'localhost:8000/entries/{id}' \
```

Request Parameters
|Field    |Type     |Description          |
|---------|---------|---------------------|
|id       |int      |ID of register       |

#### Delete
Flow: delete weight input
```mermaid
sequenceDiagram
Main ->> Handler: HandleDeleteEntry(req, resp)
Handler ->> Handler: Get Params
Handler ->> Service: Service.DeleteEntry
Service ->> Repository: Repository.DeleteById
Repository ->> Database: DELETE FROM t_entry
Database -->> Repository: err
Repository -->> Service: err
Service -->> Handler: err
alt err != nil
    Handler -->> Main: HTTP Internal Server Error 500
else err == nil
    Handler ->> Handler: Encode JSON Response
    Handler -->> Main: HTTP No Content 204
end
```

Example of Request
```curlrc
curl -X DELETE \
  'localhost:8000/entries/{id}' \
```

Request Parameters
|Field    |Type     |Description          |
|---------|---------|---------------------|
|id       |int      |ID of register       |