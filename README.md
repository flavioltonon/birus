# birus

Birus é uma API de classificação de textos escrita em Go.

## Requisitos:

- golang@1.17.1
- docker@latest
- docker-compose@latest

## Rodando o projeto:

Rodar o seguinte comando:

> make start

Para criar uma imagem Docker, o seguinte comando pode ser utilizado:

> make image

## Configurações padrão:

```
Servidor:
- server.address: :8080 // endereço padrão do servidor
- server.development_environment: true // modo de operação do servidor (true: desenvolvimento/false: release)

OCR:
- ocr.tessdata_prefix: /usr/share/tessdata/ // caminho para o diretório de dados de treinamento utilizados pela ferramenta de OCR Tesseract
- ocr.language: por // idioma utilizado pelo Tesseract

Banco de dados:
- database.kind: mongodb // tipo de banco de dados a ser utilizado
- database.uri: mongodb://localhost:27017 // endereço padrão para conexão do Birus com o banco de dados
```

A qualquer momento, é possível alterar (ambiente local) ou sobrescrever (container Docker) o arquivo de configurações da aplicação (config.yaml). No segundo caso, o arquivo deve ser colocado em `/config.yaml`.

## Funcionalidades:

### Definições

- Classifier:
```
{
    "id": string,
    "name": string
}
```

- Score:
```
{
    "name": string,
    "confidence": float64
}
```

### Criar classificador:

**Request**
```
POST /api/text-classification/classifiers
{
    "name": string // obrigatório
    "texts": []string  // obrigatório
}
```

**Response**

> Cenário: falha na validação do corpo da requisição
```
Status: 400
{
    "error": string
}
```

> Cenário: erros internos
```
Status: 500
{
    "error": string
}
```

> Cenário: classificador criado com sucesso
```
Status: 201
{
    "classifier": <classifier>
}
```

### Listar classificadores:

**Request**

```
GET /api/text-classification/classifiers
```

**Response**

> Cenário: erros internos
```
Status: 500
{
    "error": string
}
```

> Cenário: listagem realizada com sucesso
```
Status: 200
{
    "classifiers": []<classifier>
}
```

### Deletar classificador:

**Request**

```
DELETE /api/text-classification/classifiers/:classifier_id
```

**Response**

> Cenário: parâmetros de URL inválidos
```
Status: 400
{
    "error": string
}
```

> Cenário: classificador não encontrado
```
Status: 404
{
    "error": string
}
```

> Cenário: erros internos
```
Status: 500
{
    "error": string
}
```

> Cenário: classificador deletado com sucesso
```
Status: 204
```

### Classificar texto:

**Request**

```
POST /api/text-classification/classify
{
    "text": string // obrigatório
}
```

**Response**

> Cenário: parâmetros de URL inválidos
```
Status: 400
{
    "error": string
}
```

> Cenário: erros internos
```
Status: 500
{
    "error": string
}
```

> Cenário: texto classificado com sucesso
```
Status: 200
{
    "result": <score>
}
```
