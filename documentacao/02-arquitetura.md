# 2. Arquitetura

← [Anterior: Início](01-inicio.md) | [Índice](README.md) | [Próximo: Criptografia →](03-criptografia.md)

---

## Visão Geral

O sistema é composto por 4 camadas:

```
┌─────────────────────────────────────────────────┐
│                  CLI (Cobra)                    │
│          cmd/crom-pay/                          │
│    main.go → root.go → wallet.go                │
├─────────────────────────────────────────────────┤
│              Wallet Engine                      │
│           internal/wallet/                      │
│  ┌──────────┬────────────┬──────────────────┐   │
│  │  seed.go │keychain.go │    crypto.go     │   │
│  │  BIP-39  │   BIP-84   │  AES+Argon2id   │   │
│  └──────────┴────────────┴──────────────────┘   │
│                  wallet.go                      │
│          (orquestração de fluxos)                │
├─────────────────────────────────────────────────┤
│            Storage Layer                        │
│          internal/storage/                      │
│       db.go (BadgerDB wrapper)                  │
│       wallet.go (seed/config ops)               │
├─────────────────────────────────────────────────┤
│         Sistema de Arquivos                     │
│     ~/.crom-pay/wallet.db/ (BadgerDB)           │
└─────────────────────────────────────────────────┘
```

## Módulos

### `seed.go` — Geração de Entropia

Responsável por gerar e validar mnemônicos BIP-39:

- **Entropia:** 256 bits via `crypto/rand` (CSPRNG do OS)
- **Mnemônico:** 24 palavras do dicionário BIP-39
- **Seed:** PBKDF2-HMAC-SHA512 com 2048 iterações

### `keychain.go` — Derivação de Chaves

Implementa BIP-84 (Native SegWit):

```
m/84'/0'/0'/0/0   ← Mainnet, primeira address
m/84'/1'/0'/0/0   ← Testnet, primeira address
  │   │  │  │ └── índice do endereço
  │   │  │  └──── chain (0=externo, 1=troco)
  │   │  └─────── conta (sempre 0 por agora)
  │   └────────── coin type (0=BTC, 1=testnet)
  └──────────── purpose (84=SegWit nativo)
```

**Endereços gerados:** Bech32 (`bc1q...` / `tb1q...`) — P2WPKH (Pay-to-Witness-Public-Key-Hash)

### `crypto.go` — Cifra e Proteção

Responsável por cifrar/decifrar a seed no disco:

- **KDF:** Argon2id (64MB RAM, 3 iterações, 4 threads)
- **Cifra:** AES-256-GCM (AEAD)
- **Formato:** `salt(16) || nonce(12) || ciphertext`

→ Detalhes em [Criptografia](03-criptografia.md)

### `wallet.go` — Orquestração

Agrega seed + keychain + storage em operações de alto nível:

| Operação | Fluxo |
|---|---|
| `Create` | Gera mnemônico → Deriva seed → Cifra → Salva → Retorna endereço |
| `Open` | Carrega cifrada → Decifra → Cria keychain → Retorna wallet |
| `Restore` | Valida mnemônico → Deriva seed → Cifra → Salva |
| `Close` | Zero keychain → Fecha BadgerDB |

### `storage/` — Persistência

BadgerDB wrapper com prefixos para isolamento:

| Prefixo | Dados |
|---|---|
| `seed:` | Seed cifrada |
| `config:` | Rede, versão |
| `tx:` | Transações (futuro) |
| `contact:` | Contatos (futuro) |

## Fluxo de Dados: `wallet create`

```
Usuário                CLI                  Wallet              Storage
  │                    │                      │                    │
  │─── senha ─────────▶│                      │                    │
  │                    │─── Create() ────────▶│                    │
  │                    │                      │── GenerateMnemonic │
  │                    │                      │── MnemonicToSeed   │
  │                    │                      │── NewKeychain      │
  │                    │                      │── DeriveAddress(0) │
  │                    │                      │── EncryptData      │
  │                    │                      │── zeroBytes(seed)  │
  │                    │                      │── kc.Zero()        │
  │                    │                      │──────── SaveSeed ─▶│
  │                    │                      │──────── SaveConf ─▶│
  │◀── mnemônico ──────│◀── info ─────────────│                    │
```

---

← [Anterior: Início](01-inicio.md) | [Índice](README.md) | [Próximo: Criptografia →](03-criptografia.md)
