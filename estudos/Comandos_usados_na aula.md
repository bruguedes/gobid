#### 1. Instalar o Air
```bash
go install github.com/air-verse/air@latest


Objetivo:
Instalar o Air diretamente do repositório remoto para facilitar o desenvolvimento de aplicações Go, permitindo que alterações no código sejam refletidas automaticamente sem a necessidade de reiniciar manualmente o servidor.

Observação:
Certifique-se de que o diretório $GOPATH/bin esteja no seu PATH para que o comando air fique acessível no terminal após a instalação.

Observação:
Certifique-se de que o diretório $GOPATH/bin esteja no seu PATH para que o comando air fique acessível no terminal após a instalação.
```
#### 2. Criar o arquivo de configuração do Air
```bash
air --build.cmd "go build -o ./bin/api ./cmd/api" --build.bin "./bin/api"

Passo a Passo:

--build.cmd: Define o comando para compilar o código (go build), gerando o binário em ./bin/api.
--build.bin: Especifica o caminho do binário gerado para execução.
Estrutura esperada do projeto:

Código principal: ./cmd/api/main.go.
Binário gerado: ./bin/api.
Objetivo:
Monitorar alterações no código e reiniciar automaticamente a aplicação durante o desenvolvimento.

Pronto! Agora você pode focar no desenvolvimento sem precisar reiniciar manualmente o servidor.
```

#### 3. Criando migrações com tern
```bash
 tern new create_users_table
```

#### 4. Criando queries com sqlc
```bash
sqlc generate -f ./internal/store/pgstore/sqlc.yml
```



Aulas: bidsService
```bash
Dentro da pasta migrations, crie um arquivo de migração com o comando:
tern new create_bids_table
```
Rodar a migração:
```bash
tern up
```
criar arquivo bids.sqlc.yml dentro da pasta internal/store/pgstore/queries com o seguinte conteúdo, e escrever as queries necessárias.


rodar o sqlc generate, para gerar o código Go a partir das queries SQL:

```bash
sqlc generate -f ./internal/store/pgstore/sqlc.yml

Esse comando irá ler o arquivo sqlc.yml e gerar o código Go correspondente às queries definidas, facilitando a interação com o banco de dados PostgreSQL de forma tipada e segura.
O arquivo gerado esta localizado em internal/store/pgstore/queries/bids.sql.go
```
