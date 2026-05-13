# 📁 Estrutura do Projeto — Crom Bitcoin Pix

> Árvore de diretórios Go, organização de packages e convenções.

---

## Estrutura Proposta

```
crom-bitcoin-pix/
│
├── docs/                          # 📚 Documentação (você está aqui)
│   ├── 00-VISAO-GERAL.md
│   ├── 01-ARQUITETURA-TECNICA.md
│   ├── 02-PROTOCOLOS-PADROES.md
│   ├── 03-FLUXOS-USUARIO.md
│   ├── 04-ESTRUTURA-PROJETO.md
│   ├── 05-ROADMAP-MILESTONES.md
│   ├── 06-RISCOS-MITIGACOES.md
│   ├── 07-GLOSSARIO.md
│   └── 08-REFERENCIAS.md
│
├── cmd/                           # 🎮 Entrypoints CLI (cobra commands)
│   └── crom-pay/
│       ├── main.go                # func main() → executa rootCmd
│       ├── root.go                # Comando raiz, flags globais
│       ├── wallet.go              # wallet create|balance|restore|backup
│       ├── pay.go                 # pay <address> <amount>
│       ├── receive.go             # receive <amount> → gera invoice
│       ├── channels.go            # channels open|list|close
│       ├── contacts.go            # contacts add|list|remove
│       ├── config.go              # config set|get
│       └── daemon.go              # daemon start|stop|status
│
├── internal/                      # 🔒 Packages internos (não exportados)
│   │
│   ├── wallet/                    # 💰 Carteira Bitcoin
│   │   ├── seed.go                # BIP-39 mnemonic generation
│   │   ├── keychain.go            # BIP-32/44/84 key derivation
│   │   ├── address.go             # Address generation (Bech32)
│   │   ├── balance.go             # Balance via Neutrino
│   │   ├── backup.go              # Encrypted seed backup
│   │   └── wallet_test.go         # Unit tests
│   │
│   ├── lightning/                  # ⚡ Motor Lightning Network
│   │   ├── node.go                # LND embedded lifecycle
│   │   ├── channels.go            # Channel management
│   │   ├── invoices.go            # BOLT11 create/decode
│   │   ├── payments.go            # Payment routing
│   │   ├── lsp.go                 # LSP integration
│   │   ├── monitor.go             # Channel state monitor
│   │   └── lightning_test.go      # Unit tests
│   │
│   ├── nostr/                     # 🌐 Identidade Descentralizada
│   │   ├── identity.go            # Nostr keypair from BIP-39
│   │   ├── profile.go             # NIP-01 profile publish
│   │   ├── discovery.go           # User discovery via relays
│   │   ├── nwc.go                 # NIP-47 Wallet Connect
│   │   ├── zaps.go                # NIP-57 Lightning Zaps
│   │   ├── relays.go              # Relay management
│   │   └── nostr_test.go          # Unit tests
│   │
│   ├── storage/                   # 💾 Persistência Local
│   │   ├── db.go                  # BadgerDB initialization
│   │   ├── channels.go            # Channel state persistence
│   │   ├── transactions.go        # Transaction history
│   │   ├── contacts.go            # Contact book
│   │   ├── config.go              # User configuration
│   │   └── storage_test.go        # Unit tests
│   │
│   ├── lnurl/                     # 🔗 LNURL Protocol
│   │   ├── pay.go                 # LNURL-Pay (LUD-06) client
│   │   ├── address.go             # Lightning Address resolution
│   │   ├── server.go              # Mini HTTP server for receiving
│   │   └── lnurl_test.go          # Unit tests
│   │
│   └── api/                       # 🌍 REST API (localhost only)
│       ├── server.go              # HTTP server setup
│       ├── handlers.go            # Route handlers
│       ├── middleware.go          # Auth, CORS, logging
│       └── api_test.go           # Integration tests
│
├── pkg/                           # 📦 Packages exportáveis (opcionais)
│   ├── qrcode/                    # QR Code generation helpers
│   │   └── qr.go
│   └── crypto/                    # Crypto utility functions
│       └── aes.go
│
├── scripts/                       # 🛠️ Scripts de build e deploy
│   ├── build.sh                   # Cross-compilation script
│   ├── regtest.sh                 # Setup Bitcoin Regtest network
│   └── release.sh                 # GitHub Release automation
│
├── testdata/                      # 🧪 Fixtures de teste
│   ├── seeds.json                 # Known test seeds
│   └── invoices.json              # Sample BOLT11 invoices
│
├── .github/                       # 🤖 CI/CD
│   └── workflows/
│       ├── test.yml               # Run tests on push
│       └── release.yml            # Build + release on tag
│
├── go.mod                         # 📋 Go module definition
├── go.sum                         # 📋 Dependency checksums
├── Makefile                       # 🔧 Build commands
├── LICENSE                        # 📜 MIT License
└── README.md                      # 📖 Projeto overview
```

---

## Convenções

### Naming
- **Packages:** lowercase, singular (`wallet`, `lightning`, `nostr`)
- **Files:** snake_case (`channel_monitor.go`)
- **Funções exportadas:** PascalCase (`CreateWallet()`)
- **Funções internas:** camelCase (`deriveKey()`)

### Testes
- Cada package tem seu `*_test.go`
- Testes de integração no diretório `testdata/`
- CI roda `go test ./...` em cada push

### Build
```bash
# Desenvolvimento
go run ./cmd/crom-pay

# Build local
go build -o crom-pay ./cmd/crom-pay

# Release (static, cross-platform)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o crom-pay-linux-amd64 ./cmd/crom-pay
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o crom-pay-darwin-arm64 ./cmd/crom-pay
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o crom-pay-windows-amd64.exe ./cmd/crom-pay
```

### Makefile
```makefile
.PHONY: build test run clean

build:
	go build -o bin/crom-pay ./cmd/crom-pay

test:
	go test -v ./...

run:
	go run ./cmd/crom-pay

clean:
	rm -rf bin/

regtest:
	./scripts/regtest.sh

release:
	./scripts/release.sh
```
