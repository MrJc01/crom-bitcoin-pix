# 8. Contribuição

← [Anterior: Glossário](07-glossario.md) | [Índice](README.md)

---

## Como Contribuir

### 1. Fork e Clone

```bash
git clone https://github.com/SEU_USUARIO/crom-bitcoin-pix.git
cd crom-bitcoin-pix
```

### 2. Setup

```bash
# Instalar dependências
go mod download

# Build
make build

# Verificar que tudo funciona
./tests/run_all.sh
```

### 3. Desenvolvimento

```bash
# Criar branch
git checkout -b feature/minha-feature

# Desenvolver...

# Rodar testes
go test ./...

# Verificar código
go vet ./...
```

### 4. Pull Request

- Descreva o que mudou e por quê
- Inclua testes para código novo
- Garanta que `go vet` e todos os testes passam

---

## Padrões de Código

### Estrutura de Diretórios

| Diretório | Conteúdo |
|---|---|
| `cmd/` | Entry points e CLI |
| `internal/` | Lógica de negócio (não exportada) |
| `pkg/` | Bibliotecas exportáveis (futuro) |
| `tests/` | Testes de integração, segurança e E2E |
| `documentacao/` | Documentação em português |
| `docs/` | Specs técnicas |

### Convenções Go

- **Erros:** Sempre retornar `error`, nunca `panic` em produção
- **Nomes:** `camelCase` para privado, `PascalCase` para público
- **Comentários:** Doc-comments em todas as funções públicas
- **Imports:** stdlib → terceiros → interno (separados por linha vazia)

### Segurança

- **NUNCA** use `math/rand` para dados criptográficos
- **SEMPRE** zere dados sensíveis da memória após uso
- **SEMPRE** valide inputs (mnemônico, senha, paths)
- **NUNCA** logue dados sensíveis (seeds, senhas, chaves)

### Testes

- Todo código novo deve incluir testes
- Use build tags para testes pesados: `//go:build security`
- Testes devem funcionar com `t.TempDir()` (não criar lixo)
- Nomes de teste: `Test[Componente]_[Cenário]`

---

## Estrutura de Commit

```
tipo(escopo): descrição curta

Corpo opcional com mais detalhes.
```

**Tipos:**
- `feat` — Nova funcionalidade
- `fix` — Correção de bug
- `security` — Correção de segurança
- `test` — Novos testes
- `docs` — Documentação
- `refactor` — Refatoração sem mudar comportamento

**Exemplos:**
```
feat(wallet): adicionar suporte a múltiplas contas
fix(crypto): corrigir zeroing de masterKey
security(storage): fix append mutation em Get/Delete
test(storage): adicionar 12 testes CRUD
docs: criar pasta documentacao/ com 8 guias
```

---

## Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a [licença MIT](../LICENSE).

---

← [Anterior: Glossário](07-glossario.md) | [Índice](README.md)
