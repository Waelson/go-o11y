### Aplicação
Simples aplicação de exemplo sobre observabilidade utilizando Otel com Zipkin

### Pré-requisitos
- Go 1.21.1
- Docker
- Docker Compose

### Como executar?
Na raiz do projeto execute o comando abaixo
```
 docker-compose up --build 
```
### Como validar?
Espere o `docker-compose` construir a aplicação e carregar as dependências. Depois disso, execute o comando abaixo para realizar requisições para o endpoint.
```
curl -X POST -H "Content-Type: application/json" -d '{"cep": "29902555"}' http://localhost:8080/clima
```
### Como visualizar os traces?
Acesse o endereço abaixo em seu browser, clique na opção `Find a trace`, depois no botão `Run query`
```
http://127.0.0.1:9411/zipkin/
```