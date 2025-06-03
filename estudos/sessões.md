# Relatório: Implementação de Sessões, Login e Logout com SCS no Go

## 1. Biblioteca e Instalação

Para gerenciar sessões em Go, utilize a biblioteca `scs` do Alex Edwards, popular e fácil de usar.

```bash
go get github.com/alexedwards/scs/v2
go get github.com/alexedwards/scs/postgresstore
```

## 2. Migration da Tabela de Sessões

Dentro da pasta `migrations` do seu projeto, crie um arquivo usando o comando:

```bash
tern new create_sessions_table
```

E adicione o seguinte conteúdo:

```sql
CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  data BYTEA NOT NULL,
  expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```

> **Obs:** Sempre use migrations para versionar alterações no banco.

## 3. Configuração do SessionManager

No seu código Go, inicialize o SessionManager e configure o storage para usar o PostgreSQL:

```go
import (
    "github.com/alexedwards/scs/v2"
    "github.com/alexedwards/scs/postgresstore"
    "time"
)

s := scs.New()
s.Store = postgresstore.New(pool) // pool = conexão com o banco
s.Lifetime = 24 * time.Hour
s.Cookie.SameSite = http.SameSiteLaxMode
s.Cookie.HttpOnly = true
```

> **Dica:** O atributo SameSite Lax ajuda a proteger contra ataques CSRF, mas não é suportado por todos os browsers.

## 4. Registro de Tipos Customizados no GOB

Se você for salvar tipos customizados (ex: `uuid.UUID`) na sessão, registre-os:

```go
import "encoding/gob"
import "github.com/google/uuid"

func init() {
    gob.Register(uuid.UUID{})
}
```

## 5. Middleware Global

Adicione o middleware do SCS para carregar e salvar a sessão automaticamente em todas as rotas:

```go
router.Use(api.Sessions.LoadAndSave)
```

## 6. Fluxo de Login

1. Receba e valide o JSON de login (`email` e `password`).
2. Autentique o usuário (busque por email, compare senha com bcrypt).
3. Se autenticado:
   - Renove o token da sessão:
     ```go
     api.Sessions.RenewToken(r.Context())
     ```
   - Salve o ID do usuário na sessão:
     ```go
     api.Sessions.Put(r.Context(), "AuthenticatedUserID", userID)
     ```
   - Retorne mensagem de sucesso.
4. Sempre trate erros de autenticação de forma genérica para não expor se o email existe ou não.

## 7. Fluxo de Logout

1. Renove o token da sessão antes de remover dados:
   ```go
   api.Sessions.RenewToken(r.Context())
   ```
2. Remova o ID do usuário da sessão:
   ```go
   api.Sessions.Remove(r.Context(), "AuthenticatedUserID")
   ```
3. Retorne mensagem de sucesso.

## 8. Middleware de Autenticação

Para proteger rotas que exigem usuário logado:

```go
func (api *API) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !api.Sessions.Exists(r.Context(), "AuthenticatedUserID") {
            jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]any{
                "message": "must be logged in",
            })
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## 9. Resumo dos Fluxos

- **Login:** valida → autentica → renova token → salva ID na sessão → responde sucesso.
- **Logout:** renova token → remove ID da sessão → responde sucesso.
- **Proteção de rotas:** middleware verifica se ID está na sessão.

## 10. Observações Importantes

- **CSRF:** Ataque de falsificação de solicitação entre sites, é um tipo de ataque em que um invasor engana um usuário autenticado para enviar uma solicitação não intencional a um site em que o usuário está autenticado. Isso pode levar a ações indesejadas, como transferências de dinheiro, alterações de senha ou outras operações sensíveis. O uso de SameSite Lax ajuda, mas para máxima proteção implemente CSRF tokens em endpoints sensíveis.
- **Erros GOB:** Sempre registre tipos customizados que serão salvos na sessão.
- **Tratamento de Erros:** Não exponha detalhes internos em erros inesperados.
- **Mensagens Padronizadas:** Use mensagens claras e padronizadas para facilitar o consumo pelo frontend.
- **Migrations:** Sempre use migrations para versionar alterações no banco.

## 11. Dica Extra: SQLC

Sempre que precisar criar uma nova query, primeiro adicione o SQL puro no arquivo `users.sql` e depois gere o código com o comando:

```bash
sqlc generate -f ./internal/store/pgstore/sqlc.yml
```

## 12. Extras e Boas Práticas

- **Comentários no Código:** Comente pontos críticos para facilitar manutenção.
- **Testes:** Implemente testes para login/logout e middleware de autenticação.
- **Logs:** Considere adicionar logs para tentativas de login/logout para auditoria e troubleshooting.

---

Com esse guia, você cobre não só o passo a passo da implementação, mas também pontos de segurança, manutenção e boas práticas para um sistema de autenticação robusto em Go usando SCS!
