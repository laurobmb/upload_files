# API de Upload de Arquivos em Go

Uma API web simples, construída com a biblioteca padrão `net/http` do Go, que fornece um endpoint para o upload de arquivos com validações de tipo e tamanho.

## 📜 Descrição

Este projeto implementa um servidor HTTP que expõe um único endpoint (`/upload`) para receber arquivos através de requisições `POST` no formato `multipart/form-data`. O servidor foi projetado para ser leve e possui regras de validação para garantir que apenas arquivos com extensões permitidas e dentro de um limite de tamanho sejam aceitos.

Além disso, um middleware de logging foi implementado para registrar todas as requisições recebidas no console, facilitando o monitoramento e a depuração.

## ✨ Funcionalidades

* **Endpoint de Upload:** Rota `POST /upload` para receber arquivos.
* **Validação de Extensão:** Permite apenas o upload de arquivos com as seguintes extensões: `txt`, `adoc`, `md`, `yaml`, `yml`.
* **Validação de Tamanho:** Impõe um limite máximo de **10 MB** por arquivo.
* **Armazenamento Local:** Salva os arquivos recebidos em um diretório local `./uploads`, que é criado automaticamente.
* **Logging de Acesso:** Imprime no console os detalhes de cada requisição recebida (método, URI, endereço remoto) e o tempo de processamento.

## ⚙️ Pré-requisitos

* **Go:** É necessário ter o Go (versão 1.18 ou superior) instalado em sua máquina. Você pode baixá-lo em [go.dev](https://go.dev/dl/).

## 🚀 Como Executar

1.  **Salve o Código:**
    Salve o código-fonte em um arquivo chamado `main.go`.

2.  **Navegue até o Diretório:**
    Abra seu terminal e navegue até a pasta onde você salvou o arquivo `main.go`.

3.  **Execute o Servidor:**
    Rode o seguinte comando no terminal:
    ```sh
    go run main.go
    ```

4.  **Servidor no Ar:**
    Se tudo ocorrer bem, você verá a seguinte mensagem, indicando que o servidor está pronto para receber requisições na porta `8080`:
    ```
    Servidor iniciado em http://localhost:8080
    Endpoint de upload disponível em POST http://localhost:8080/upload
    ```

## 📡 Como Usar (API)

### Endpoint: `POST /upload`

Para enviar um arquivo, faça uma requisição `POST` para `http://localhost:8080/upload` com o corpo do tipo `multipart/form-data`. O campo do arquivo deve se chamar `file`.

### Exemplos de Uso com `cURL`

#### 1. Upload com Sucesso

Crie um arquivo de teste (ex: `teste.txt`) e execute o comando abaixo, substituindo `/caminho/para/seu/teste.txt` pelo caminho real do arquivo.

```sh
curl -X POST -F "file=@/caminho/para/seu/teste.txt" http://localhost:8080/upload
````

**Resposta esperada:**

```
Upload do arquivo 'teste.txt' realizado com sucesso!
```

#### 2\. Erro: Extensão de Arquivo Inválida

Tente enviar um arquivo com uma extensão não permitida, como `.zip`.

```sh
curl -X POST -F "file=@/caminho/para/seu/arquivo.zip" http://localhost:8080/upload
```

**Resposta esperada (Erro 400 Bad Request):**

```
Tipo de arquivo não permitido. Extensões aceitas: txt, adoc, md, yaml, yml
```

#### 3\. Erro: Arquivo Muito Grande

Tente enviar um arquivo com mais de 10 MB.

```sh
# (Assumindo que 'arquivo_grande.dat' tem mais de 10 MB)
curl -X POST -F "file=@/caminho/para/seu/arquivo_grande.dat" http://localhost:8080/upload
```

**Resposta esperada (Erro 400 Bad Request):**

```
O arquivo excede o limite de tamanho de 10 MB
```