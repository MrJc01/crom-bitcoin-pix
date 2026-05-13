# 7. Glossário

← [Anterior: Roadmap](06-roadmap.md) | [Índice](README.md) | [Próximo: Contribuição →](08-contribuicao.md)

---

## Bitcoin

| Termo | Definição |
|---|---|
| **Bitcoin** | Moeda digital descentralizada, sem banco central |
| **Satoshi (sat)** | Menor unidade de Bitcoin. 1 BTC = 100.000.000 sats |
| **UTXO** | Unspent Transaction Output — "moeda" não gasta |
| **SegWit** | Segregated Witness — formato de transação mais eficiente |
| **Bech32** | Formato de endereço SegWit (`bc1q...`) |
| **P2WPKH** | Pay-to-Witness-Public-Key-Hash — tipo de script SegWit |
| **Mainnet** | Rede principal Bitcoin (dinheiro real) |
| **Testnet** | Rede de testes (bitcoins sem valor) |
| **Regtest** | Rede local para desenvolvimento |

## Chaves e Derivação

| Termo | Definição |
|---|---|
| **BIP** | Bitcoin Improvement Proposal — documento de padrão |
| **BIP-39** | Padrão de mnemônico (lista de palavras para backup) |
| **BIP-32** | Derivação Hierárquica Determinística (HD) |
| **BIP-44** | Derivação multi-conta e multi-moeda |
| **BIP-84** | Derivação para endereços SegWit nativos |
| **HD Wallet** | Carteira Hierárquica Determinística — todas as chaves de uma seed |
| **Mnemônico** | 24 palavras que representam a chave mestra |
| **Seed** | 512 bits derivados do mnemônico (raiz de tudo) |
| **Master Key** | Chave privada mestra derivada da seed |
| **Derivation Path** | Caminho hierárquico (ex: `m/84'/0'/0'/0/0`) |
| **Hardened** | Derivação com `'` — protege contra leaked child key attack |
| **secp256k1** | Curva elíptica usada pelo Bitcoin |

## Criptografia

| Termo | Definição |
|---|---|
| **AES-256-GCM** | Cifra simétrica com autenticação (256 bits) |
| **AEAD** | Authenticated Encryption with Associated Data |
| **Argon2id** | KDF resistente a GPU e side-channel attacks |
| **KDF** | Key Derivation Function — transforma senha em chave |
| **CSPRNG** | Cryptographically Secure Pseudo-Random Number Generator |
| **Salt** | Valor aleatório que torna cada cifra única |
| **Nonce** | Number Used Once — evita replay de cifra |
| **PBKDF2** | Password-Based Key Derivation Function 2 |
| **HMAC-SHA512** | Hash-based Message Authentication Code com SHA-512 |
| **RIPEMD-160** | Hash de 160 bits usado no Bitcoin (Hash160) |

## Lightning Network

| Termo | Definição |
|---|---|
| **Lightning** | Protocolo de pagamento instantâneo sobre Bitcoin |
| **LND** | Lightning Network Daemon — implementação em Go |
| **Neutrino** | Light client BIP-157/158 (sem full node) |
| **Channel** | Canal de pagamento entre dois nós Lightning |
| **Invoice** | Pedido de pagamento Lightning (BOLT-11) |
| **BOLT** | Basis of Lightning Technology — specs do Lightning |

## Nostr

| Termo | Definição |
|---|---|
| **Nostr** | Notes and Other Stuff Transmitted by Relays |
| **NIP** | Nostr Implementation Possibilities — padrão Nostr |
| **Relay** | Servidor que armazena e retransmite eventos |
| **Zap** | Gorjeta Lightning via protocolo Nostr (NIP-57) |

## Infraestrutura

| Termo | Definição |
|---|---|
| **BadgerDB** | Banco de dados key-value embarcado em Go |
| **Cobra** | Framework CLI para Go |
| **Go Modules** | Sistema de dependências do Go |

---

← [Anterior: Roadmap](06-roadmap.md) | [Índice](README.md) | [Próximo: Contribuição →](08-contribuicao.md)
