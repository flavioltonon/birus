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

- Classifier: é um classificador de textos para um tipo de documento específico (ex: cupons fiscais eletrônicos, RGs, ou livros do Tolkien)
```
{
    "id": string,
    "name": string
}
```

- Score: é o resultado da classificação de um texto. O grau de confiança no resultado (confidence) pode variar entre 0 e 100.
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
Content-Type: application/json
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
Content-Type: application/json
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
Content-Type: application/json
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
Content-Type: application/json
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

### Processar uma imagem:

**Request**

```
POST /api/ocr/read

1) Content-Type: application/json
{
    "base64": string // obrigatório
    "options": string // opcional; exemplo: "grayscale;resize:1080,720;adjust-contrast:50;adjust-brightness:50;blur:2.5;sharpen:1.2" 
}

2) Content-Type: multipart/form-data
- file: multipart file // obrigatório
- options: string // opcional; exemplo: "grayscale;resize:1080,720;adjust-contrast:50;adjust-brightness:50;blur:2.5;sharpen:1.2" 
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

> Cenário: texto extraído com sucesso
```
Status: 200
{
    "image": {
        "base64": string
    }
}

// Obs: Caso a API esteja rodando em modo development (ocr.development_mode: true), a imagem resultante será gravada no diretório /output/image.jpg
```

### Extrair texto de uma imagem (OCR):

**Request**

```
POST /api/ocr/read

1) Content-Type: application/json
{
    "base64": string // obrigatório
    "options": string // opcional; exemplo: "grayscale;resize:1080,720;adjust-contrast:50;adjust-brightness:50;blur:2.5;sharpen:1.2" 
}

2) Content-Type: multipart/form-data
- file: multipart file // obrigatório
- options: string // opcional; exemplo: "grayscale;resize:1080,720;adjust-contrast:50;adjust-brightness:50;blur:2.5;sharpen:1.2" 
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

> Cenário: texto extraído com sucesso
```
Status: 200
{
    "text": string
}
```

### Extrair texto de várias imagens (OCR):

**Request**

```
POST /api/ocr/read/batch

1) Content-Type: application/json
{
    "base64_list": []string // obrigatório
    "options": string // opcional; exemplo: "grayscale;resize:1080,720;adjust-contrast:50;adjust-brightness:50;blur:2.5;sharpen:1.2" 
}

2) Content-Type: multipart/form-data
- files: []multipart file // obrigatório
- options: string // opcional; exemplo: "grayscale;resize:1080,720;adjust-contrast:50;adjust-brightness:50;blur:2.5;sharpen:1.2" 
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

> Cenário: texto extraído com sucesso
```
Status: 200
{
    "texts": []string
}
```
