# 🧪 Estrutura de Testes — Crom Bitcoin Pix

> Suíte completa de testes cobrindo unitários, integração, segurança/pentest, stress e E2E.

---

## Estrutura

```
tests/
├── README.md            # Este arquivo
├── run_all.sh           # Runner principal — executa tudo e gera relatórios
│
├── unit/                # Testes unitários Go (por package)
│   └── run.sh           # Runner isolado
│
├── integration/         # Testes de integração (wallet → storage → CLI)
│   ├── wallet_flow_test.go
│   └── run.sh
│
├── security/            # 🔴 Pentest e auditoria de segurança
│   ├── entropy_test.go       # Qualidade da entropia criptográfica
│   ├── crypto_audit_test.go  # Auditoria da cifra AES-GCM + Argon2id
│   ├── key_exposure_test.go  # Verificar vazamento de chaves em memória
│   ├── bruteforce_test.go    # Resistência a brute-force de senha
│   ├── path_traversal_test.go # Injection no data-dir
│   └── run.sh
│
├── stress/              # Testes de carga e limites
│   ├── concurrent_test.go    # Múltiplas carteiras simultâneas
│   ├── storage_test.go       # Limites do BadgerDB
│   └── run.sh
│
├── e2e/                 # End-to-End (binário compilado)
│   ├── cli_test.sh      # Teste do binário real via shell
│   └── run.sh
│
├── fixtures/            # Dados de teste conhecidos
│   └── test_vectors.json
│
└── reports/             # Relatórios gerados (gitignored)
    ├── .gitkeep
    └── (gerados pelo run_all.sh)
```

## Como Usar

```bash
# Rodar TODOS os testes e gerar relatório
./tests/run_all.sh

# Rodar apenas uma camada
./tests/unit/run.sh
./tests/security/run.sh
./tests/e2e/run.sh

# Ver relatórios
ls tests/reports/
cat tests/reports/report_YYYYMMDD_HHMMSS.md
```

## Categorias de Teste

| Camada | Foco | Ferramentas |
|---|---|---|
| **Unit** | Lógica isolada por package | `go test` |
| **Integration** | Fluxos entre packages | `go test -tags integration` |
| **Security** | Pentest, auditoria cripto | `go test -tags security` |
| **Stress** | Performance, limites | `go test -tags stress -timeout 120s` |
| **E2E** | Binário real, CLI completo | `bash` scripts |
