# Rate Limiter 🔐

Rate limiter simples e eficiente em Go, com suporte a controle por IP ou por token e persistência no Redis.

## 🧩 Funcionalidades

- ✅ Middleware para controle de requisições por `API_KEY` (token) ou IP
- 🔒 Bloqueio temporário de IPs ou tokens que excedem o limite
- 💾 Persistência em Redis
- ⚙️ Configuração via `.env`
- 🧪 Testes automatizados com cobertura para limites e bloqueios

---

## 🚀 Como Executar

### 1. Clone o projeto

```bash
git clone https://github.com/HenriqueCDessen/rate-limiter.git
cd rate-limiter
```

### 2. Crie um .env com base no arquivo .env.example
| Variável             | Descrição                                    | Padrão           |
| -------------------- | -------------------------------------------- | ---------------- |
| `REDIS_ADDR`         | Endereço do Redis                            | `localhost:6379` |
| `REDIS_PASSWORD`     | Senha do Redis (se houver)                   | `""`             |
| `REDIS_DB`           | Número do banco Redis                        | `0`              |
| `EXPIRATION_SECONDS` | Tempo de expiração do contador (em segundos) | `60`             |
| `IP_LIMIT`           | Limite de requisições por IP                 | `5`              |
| `TOKEN_LIMIT`        | Limite de requisições por token (`API_KEY`)  | `10`             |


### 3. Rode do Docker Compose
```bash
docker compose up --build
```

### 4. Tests
Ao ter o container em execução podemos executar
```bash
go test -v ./...
```
ou teste de integração 
```bash
go run load_validation.go
```
### 📦 Estrutura do Projeto
```bash
├── cmd/server/main.go       
│
├── config/                  # Carregamento de configurações via .env
│   └── config.go            # Função LoadConfig que carrega e valida as variáveis de ambiente

├── internal/
│   ├── limiter/             # Lógica principal de rate limiting
│   │   ├── redis_limiter.go # Implementação de Store usando Redis (com Allow, Block, IsBlocked)
│   │   └── store.go         # Interface Store para permitir diferentes implementações de limiter

│   ├── middleware/          # Middleware Gin
│   │   └── rate_limit.go    # Middleware que aplica as regras de rate limit por IP ou token

├── .env.example             # Exemplo de configuração do ambiente

├── load_validation.go       # teste de integração

├── go.mod / go.sum          # Gerenciamento de dependências Go
```



