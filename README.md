# Tasks API (Golang)

API simples e moderna para gerenciamento de tarefas, com autenticação JWT e usuários protegidos.

## 📚 Sobre o Projeto

- CRUD de tasks, com validação robusta e feedback amigável.
- Login e cadastro de usuários, senha sempre criptografada.
- Rotas de tasks protegidas com JWT (exceto GET geral).
- Testes práticos via arquivos `.http` (facilidade para frontend ou testes manuais de API)
- Código bem comentado e preparado para manutenção!

## 🚀 Stack Utilizada
- [Go](https://golang.org/) 1.20+
- [Gin](https://github.com/gin-gonic/gin) (framework web)
- [SQLite (Gorm)](https://gorm.io/)
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) para senha
- [JWT](https://github.com/golang-jwt/jwt) para autenticação
- [go-playground/validator](https://github.com/go-playground/validator) para validação

## ⚡ Como rodar local

1. Clone o projeto
    ```bash
    git clone <repo_url>
    cd tasks
    ```
2. Instale as dependências:
    ```bash
    go mod tidy
    ```
3. Rode a aplicação:
    ```bash
    go run main.go
    ```
4. Acesse via [http://localhost:8080](http://localhost:8080)

> A primeira execução irá criar o banco `app.db` automaticamente.

## 🔑 Variáveis de ambiente

- `JWT_SECRET` (opcional): chave secreta para assinar os tokens. Use em produção:
    ```bash
    export JWT_SECRET="sua-senha-super-secreta"
    ```

## 💡 Principais Rotas

### Usuários
- `POST /auth/register` — cadastro de usuário `{ "name": "user", "password": "senha" }`
- `POST /auth/login` — retorna JWT ao passar login correto

### Tasks
- `GET /tasks/` — pública, lista todas as tarefas
- Rotas protegidas (usar `Authorization: Bearer <token>`):
    - `POST /tasks/` — cria tarefa
    - `GET /tasks/:id` — consulta uma tarefa específica
    - `PUT /tasks/:id` — atualiza uma tarefa
    - `DELETE /tasks/:id` — remove uma tarefa

### Exemplos de testes rápidos
Use os arquivos em `tests/tasks.http` e `tests/users.http` com o plugin [REST Client VSCode](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) ou qualquer cliente HTTP (Insomnia/Postman/cURL).

---

## 🙌 Boas práticas
- Sempre hash a senha do usuário (nunca armazene puro)
- Use JWT para garantir segurança nas rotas críticas
- Siga os exemplos de respostas padronizadas (sempre tem `ok: true`, e campo `message`)
- O código está organizado em camadas: `models/`, `routers/`, `middlewares/`, `utils/`, `tests/`

---

## 🛠️ Para contribuir
1. Fork esse repositório, crie sua branch e envie PRs;
2. Siga o padrão de mensagens de commit e mantenha os comentários claros;
3. Se encontrar bug, crie um teste `.http` para reproduzir.

---

Projeto criado para estudos e estruturação de boas práticas em backend Go — fique à vontade para usar, sugerir melhorias ou adaptar ao seu contexto!
