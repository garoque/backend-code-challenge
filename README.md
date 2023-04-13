<img src="pkg/assets/snapfi.png" align="right" height="178" alt="Snapfi"/>

# Snapfi Backend Code Challenge

A Snapfi se propõe a resolver um problema identificado entre os MEI e autônomos: A dificuldade de gestão financeira que essas pessoa tem. Para isso, o desafio foi implementar uma API em Go capaz de simular uma transação.

## Dependências

- 🐳 [Docker](https://docs.docker.com/desktop/)
- [Golang](https://golang.org/doc/install)
- [Goose](https://github.com/pressly/goose)

## Como rodar

Esse projeto possui um makefile, após instaladas as dependências podemos rodar os seguintes comandos:
1° Rode o comando `make run` para iniciar o container docker e a API.
2° Em outro terminal, rode o comando `make mig-up` para criar as tabelas necessárias no banco de dados.

Pronto! Sua aplicação estará disponível rodando localhost na porta `:1323`.

Caso deseje parar o container docker, há disponível o comando `make stop`.

## Como usar a API

1° Criar dois usuários:
    - É necessário criar ao menos dois usuários para simularmos uma transação;
    - Para isso, temos o endpoint `http://localhost:1323/v1/user [POST]`, que aceita no body param um json com o campo `name`. Exemplo:
    ```javascript
        {
            "name": "Gabriel"
        }
    ```
    - Podemos obter a lista de usuários criados com o endpoint `http://localhost:1323/v1/user [GET]`;

2° Incrementar o saldo de ao menos um dos usuários criados:
    - Para simular uma transação, é necessário que o usuário tenha um saldo disponível;
    - Para isso, temos o endpoint `http://localhost:1323/v1/transaction/increase-balance [PUT]`, que aceita no body param um json com os campos `userId`, que é o ID do usuário que será incrementado o valor e `value`, que é o valor a ser incrementado no saldo. Exemplo:
    ```javascript
    {
        "userId": "3fe00197-4116-43a1-815d-4635acd4f3a2",
        "value": 100.00
    }
    ```

3° Realizar uma transação entre dois usuários:
    - Para realizarmos uma transação, temos o endpoint `http://localhost:1323/v1/transaction [POST]`, que aceita no body param um json com os campos `sourceUserId`, que é o ID do usuário que está realizando a transação, ou seja, de onde será debitado o valor, o outro campo é o `destinationUserId`, que é o ID do usuário que irá receber o valor e o campo `amount`, que é a quantia transacionada. Exemplo:
    ```javascript
    {
        "sourceUserId": "3fe00197-4116-43a1-815d-4635acd4f3a2",
        "destinationUserId": "b2eb6dc7-1fcd-444a-a973-74d022eb848b",
        "amount": 100.00
    }
    ```