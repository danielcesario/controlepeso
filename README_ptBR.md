# App Controle de Peso

- [Overview](#overview)
- [Executando](#executando)
    - [Base de Dados](#base-de-dados)
    - [Rodando Aplicação](#aplicação)
- [API](#api)
    - [Entradas](#entradas)
        - [Criação](#criação)
        - [Atualização](#atualização)
        - [Listagem](#listagem)
        - [Obtenção](#obtenção)
        - [Exclusão](#exclusão)

## Overview
Projeto simples para treinar e implementar simples conceitos da linguagem GO. Uma simples API para controle de peso utilizando como base de dados um banco relacional (PostgreSQL).
O conteúdo do projeto será incrementado conforme os estudos de novos conceitos.

## Executando

### Base de Dados
Para iniciar o container do banco de dados executar o seguinte comando no terminal
```sh
docker-compose up -d
```

### Aplicação
Para iniciar a API executar o comando na raiz do projeto
```sh
go run cmd/main.go
```

## API

### Entradas
Lançamentos de peso

|Campo    |Tipo     |Descrição                    |
|---------|---------|-----------------------------|
|id       |int      |Id único de cada entrada     |
|user_id  |int      |Id do usuário relacionado    |
|weight   |float64  |Peso registrado              |
|date     |string   |Data do Registro             |


#### Criação
Fluxo da criação de um novo lançamento de peso
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

Exemplo de Request
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

#### Atualização
Fluxo da atualização de um lançamento de peso
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

Exemplo de Request
```curlrc
curl -X PUT \
  'localhost:8000/entries/{id}' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "weight": 100.38,
  "date": "2022-01-01T10:00:00Z"
}'
```

Parâmetros da Request
|Campo    |Tipo     |Descrição                        |
|---------|---------|---------------------------------|
|id       |int      |ID do Registro a ser atualizado  |

#### Listagem
Fluxo da listagem de lançamentos de peso
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

Exemplo de Request
```curlrc
curl -X GET \
  'localhost:8000/entries?start=0&count=20' \
```

Parâmetros da Request
|Campo    |Tipo     |Descrição                    |
|---------|---------|-----------------------------|
|start    |int      |Registro inicial da pesquisa |
|count    |int      |Quantidade de registros      |

#### Obtenção
Fluxo de obtenção de um lançamentos de peso
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

Exemplo de Request
```curlrc
curl -X GET \
  'localhost:8000/entries/{id}' \
```

Parâmetros da Request
|Campo    |Tipo     |Descrição                    |
|---------|---------|-----------------------------|
|id       |int      |ID do Registro buscado       |

#### Exclusão
Fluxo de exclusão de um lançamentos de peso
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

Exemplo de Request
```curlrc
curl -X DELETE \
  'localhost:8000/entries/{id}' \
```

Parâmetros da Request
|Campo    |Tipo     |Descrição                        |
|---------|---------|---------------------------------|
|id       |int      |ID do Registro a ser excluído    |