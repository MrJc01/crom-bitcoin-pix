# 🏛️ Arquitetura Técnica — Crom Bitcoin Pix

> Detalhamento dos componentes internos, stack tecnológico e decisões arquiteturais.

---

## 1. Stack Tecnológico

### Linguagem: Go (Golang)

**Por que Go?**
- Compilação estática: binário único sem dependências (CGO_ENABLED=0)
- Concorrência nativa: goroutines para canais Lightning + conexões Nostr
- Ecossistema Bitcoin nativo: LND, btcsuite, Neutrino são Go
- Cross-compilation trivial: um comando gera binários para Linux/Windows/macOS

### Dependências Core

| Biblioteca | Função | Repositório |
|---|---|---|
| `btcsuite/btcd` | Primitivas Bitcoin | github.com/btcsuite/btcd |
| `btcsuite/btcwallet` | Carteira HD | github.com/btcsuite/btcwallet |
| `lightningnetwork/lnd` | Nó Lightning embutido | github.com/lightningnetwork/lnd |
| `lightninglabs/neutrino` | Light client BIP-157/158 | github.com/lightninglabs/neutrino |
| `nbd-wtf/go-nostr` | Protocolo Nostr | github.com/nbd-wtf/go-nostr |
| `dgraph-io/badger` | DB embutido key-value | github.com/dgraph-io/badger |
| `spf13/cobra` | Framework CLI | github.com/spf13/cobra |
| `skip2/go-qrcode` | Geração de QR Codes | github.com/skip2/go-qrcode |

---

## 2. Componentes Internos

### 2.1 Wallet Core (`internal/wallet/`)

Lógica criptográfica e gerenciamento de chaves.

```
wallet/
├── seed.go          # Geração/restauração de sementes BIP-39
├── keychain.go      # Derivação de chaves BIP-32/BIP-44
├── address.go       # Endereços Bitcoin (P2WPKH/Bech32)
├── balance.go       # Saldo on-chain via Neutrino
└── backup.go        # Exportação segura da semente
```

**Fluxo de Criação:**
```
Entropia (256 bits) → Mnemônico (24 palavras) → Master Seed → Master Key → m/84'/0'/0'/0/0
```

Derivação BIP-84 (SegWit nativo) por ser padrão moderno com taxas menores.

### 2.2 Lightning Engine (`internal/lightning/`)

```
lightning/
├── node.go          # Lifecycle do nó LND embutido
├── channels.go      # Abertura/fechamento de canais
├── invoices.go      # Invoices BOLT11
├── payments.go      # Envio e roteamento
├── lsp.go           # Lightning Service Providers
└── monitor.go       # Channel Monitor (anti-fraude)
```

> **Por que LND e não LDK?** O LDK não possui bindings oficiais para Go. O LND é escrito nativamente em Go, possui API gRPC madura e ~70% market share na rede Lightning.

**Neutrino (Light Client):**
- NÃO baixa a blockchain (500GB+)
- Baixa apenas filtros compactos (~4MB/mês)
- Full node NÃO sabe quais endereços pertencem ao usuário (privacidade)

### 2.3 Nostr Discovery (`internal/nostr/`)

```
nostr/
├── identity.go      # Keypair Nostr (derivada do BIP-39)
├── profile.go       # Perfil público (NIP-01)
├── discovery.go     # Busca de endereços de pagamento
├── nwc.go           # Nostr Wallet Connect (NIP-47)
├── zaps.go          # Lightning Zaps (NIP-57)
└── relays.go        # Gerenciamento de relays
```

Nostr substitui servidor central: Alice publica endereço nos relays, Bob busca e paga. O relay apenas repassa eventos.

### 2.4 Local Storage (`internal/storage/`)

```
storage/
├── db.go            # BadgerDB init
├── channels.go      # Estado dos canais (CRITICAL)
├── transactions.go  # Histórico
├── contacts.go      # Contatos
└── config.go        # Configurações
```

> **Channel State é Crítico:** perda = perda de fundos. Implementar SCB (Static Channel Backup), watchtower support e graceful shutdown.

### 2.5 Interface Layer (`cmd/`)

```
cmd/
├── root.go          # cobra root
├── wallet.go        # wallet create | balance | restore
├── pay.go           # pay <endereço> <valor>
├── receive.go       # receive <valor> → QR Code
├── channels.go      # channels open | list | close
├── contacts.go      # contacts add | list
└── daemon.go        # daemon start → nó + API REST
```

**API REST (localhost:8420):**
```
GET  /api/v1/wallet/balance
POST /api/v1/pay
POST /api/v1/receive
GET  /api/v1/channels
GET  /api/v1/transactions
```

API roda somente em localhost para futuras UIs (Tauri, web local).

---

## 3. Requisitos de Sistema

| Recurso | Mínimo | Recomendado |
|---|---|---|
| RAM | 128 MB | 256 MB |
| Disco | 500 MB | 1 GB |
| Rede | Internet | Internet estável |
| CPU | 1 core | 2+ cores |
| OS | Linux x64 | Linux/macOS/Windows |

Binário final: ~15-30 MB (Go static build com LND embutido).
