# 3. Criptografia

← [Anterior: Arquitetura](02-arquitetura.md) | [Índice](README.md) | [Próximo: Comandos →](04-comandos.md)

---

## Visão Geral

A criptografia do Crom Bitcoin Pix segue as melhores práticas da indústria Bitcoin e as recomendações OWASP para armazenamento de segredos.

## Cadeia Criptográfica

```
crypto/rand (CSPRNG)
       │
       ▼
  256 bits de entropia
       │
       ▼
  BIP-39 Mnemônico (24 palavras)
       │
       ▼
  PBKDF2-HMAC-SHA512 (2048 iter)
       │
       ▼
  512 bits de Seed
       │
       ├──▶ BIP-32 Master Key (para derivação)
       │         │
       │         ▼
       │    BIP-84 Path: m/84'/coin'/0'/0/index
       │         │
       │         ▼
       │    secp256k1 Public Key (33 bytes comprimida)
       │         │
       │         ▼
       │    RIPEMD160(SHA256(pubkey)) → P2WPKH → bc1q...
       │
       └──▶ Argon2id (senha do usuário)
                 │
                 ▼
            AES-256-GCM (cifra a seed no disco)
```

## Primitivas Utilizadas

### 1. Geração de Entropia

| Parâmetro | Valor |
|---|---|
| Fonte | `crypto/rand` (CSPRNG do sistema operacional) |
| Tamanho | 256 bits (32 bytes) |
| Qualidade | Shannon Entropy 7.98/8.0 (verificado) |

**Nunca usamos `math/rand`** — todas as fontes de aleatoriedade são criptograficamente seguras.

### 2. Mnemônico BIP-39

| Parâmetro | Valor |
|---|---|
| Palavras | 24 (256 bits de entropia + 8 bits checksum) |
| Dicionário | 2048 palavras inglesas (BIP-39 wordlist) |
| Checksum | SHA-256 dos primeiros 8 bits |

### 3. Derivação de Seed

| Parâmetro | Valor |
|---|---|
| Algoritmo | PBKDF2-HMAC-SHA512 |
| Iterações | 2048 (padrão BIP-39) |
| Output | 512 bits (64 bytes) |

### 4. Derivação de Chaves (BIP-32/84)

| Parâmetro | Valor |
|---|---|
| Curva | secp256k1 |
| Path | `m/84'/coin'/0'/0/index` |
| Endereço | P2WPKH (SegWit nativo, Bech32) |

### 5. Cifra de Armazenamento

| Parâmetro | Valor |
|---|---|
| KDF | Argon2id |
| Memória | 64 MB |
| Iterações | 3 |
| Paralelismo | 4 threads |
| Salt | 128 bits (aleatório por cifra) |
| Cifra | AES-256-GCM |
| Nonce | 96 bits (aleatório por cifra) |
| Formato | `salt(16) \|\| nonce(12) \|\| ciphertext+tag` |

### 6. Proteção de Memória

| Mecanismo | Descrição |
|---|---|
| `zeroBytes()` | Zera slices de dados sensíveis após uso |
| `runtime.KeepAlive` | Impede que o compilador otimize o zeroing |
| `Keychain.Zero()` | Zera a master key ao fechar a wallet |
| Permissões `0700` | Diretório do DB acessível só pelo dono |

## Resistência a Ataques

| Ataque | Mitigação |
|---|---|
| Brute-force | Argon2id: ~1 tentativa/segundo |
| GPU cracking | Argon2id: 64MB RAM por tentativa |
| Side-channel | Argon2id variante "id" resiste a side-channel |
| Timing | Argon2id domina tempo (ratio 0.73 testado) |
| Tamper | AES-GCM: tag de autenticação detecta adulteração |
| Truncation | AES-GCM: rejeita dados incompletos |
| Cold boot | `zeroBytes()` + `Keychain.Zero()` minimiza exposição |

## Verificação BIP-84

O endereço gerado foi verificado contra o vetor oficial da especificação:

```
Mnemônico: abandon abandon ... about (12 palavras de teste)
Esperado:  bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu
Obtido:    bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu ✅
```

Isso garante **interoperabilidade** com qualquer wallet BIP-84 (Electrum, BlueWallet, Sparrow, etc).

---

← [Anterior: Arquitetura](02-arquitetura.md) | [Índice](README.md) | [Próximo: Comandos →](04-comandos.md)
