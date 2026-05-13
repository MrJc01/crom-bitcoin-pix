# 🔗 Referências — Crom Bitcoin Pix

> Links para especificações, bibliotecas, projetos de referência e material de estudo.

---

## 1. Especificações de Protocolo

### Bitcoin (BIPs)
| BIP | Título | Link |
|---|---|---|
| BIP-32 | Hierarchical Deterministic Wallets | [github.com/bitcoin/bips](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki) |
| BIP-39 | Mnemonic Code for Generating Keys | [github.com/bitcoin/bips](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) |
| BIP-44 | Multi-Account Hierarchy | [github.com/bitcoin/bips](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki) |
| BIP-84 | Derivation for P2WPKH (SegWit) | [github.com/bitcoin/bips](https://github.com/bitcoin/bips/blob/master/bip-0084.mediawiki) |
| BIP-157 | Compact Block Filters (Client) | [github.com/bitcoin/bips](https://github.com/bitcoin/bips/blob/master/bip-0157.mediawiki) |
| BIP-158 | Compact Block Filters (Construction) | [github.com/bitcoin/bips](https://github.com/bitcoin/bips/blob/master/bip-0158.mediawiki) |

### Lightning Network (BOLTs)
| BOLT | Título | Link |
|---|---|---|
| BOLT #11 | Invoice Protocol (Payment Encoding) | [github.com/lightning/bolts](https://github.com/lightning/bolts/blob/master/11-payment-encoding.md) |
| BOLTs | Todas as especificações | [github.com/lightning/bolts](https://github.com/lightning/bolts) |

### LNURL (LUDs)
| LUD | Título | Link |
|---|---|---|
| LUD-06 | Pay Request (LNURL-Pay) | [github.com/lnurl/luds](https://github.com/lnurl/luds/blob/luds/06.md) |
| LUDs | Todas as especificações | [github.com/lnurl/luds](https://github.com/lnurl/luds) |
| Lightning Address | Spec | [lightningaddress.com](https://lightningaddress.com/) |

### Nostr (NIPs)
| NIP | Título | Link |
|---|---|---|
| NIP-01 | Basic Protocol | [github.com/nostr-protocol/nips](https://github.com/nostr-protocol/nips/blob/master/01.md) |
| NIP-05 | DNS-Based Verification | [github.com/nostr-protocol/nips](https://github.com/nostr-protocol/nips/blob/master/05.md) |
| NIP-47 | Nostr Wallet Connect | [github.com/nostr-protocol/nips](https://github.com/nostr-protocol/nips/blob/master/47.md) |
| NIP-57 | Lightning Zaps | [github.com/nostr-protocol/nips](https://github.com/nostr-protocol/nips/blob/master/57.md) |
| NIPs | Todas as especificações | [github.com/nostr-protocol/nips](https://github.com/nostr-protocol/nips) |

---

## 2. Bibliotecas Go

### Bitcoin / Lightning
| Biblioteca | Descrição | Link |
|---|---|---|
| btcsuite/btcd | Full node Bitcoin em Go + primitivas | [github.com/btcsuite/btcd](https://github.com/btcsuite/btcd) |
| btcsuite/btcutil | Utilidades (endereços, encoding, etc) | [github.com/btcsuite/btcutil](https://github.com/btcsuite/btcutil) |
| btcsuite/btcwallet | HD Wallet implementation | [github.com/btcsuite/btcwallet](https://github.com/btcsuite/btcwallet) |
| lightningnetwork/lnd | Lightning Network Daemon | [github.com/lightningnetwork/lnd](https://github.com/lightningnetwork/lnd) |
| lightninglabs/neutrino | BIP-157/158 Light Client | [github.com/lightninglabs/neutrino](https://github.com/lightninglabs/neutrino) |

### Nostr
| Biblioteca | Descrição | Link |
|---|---|---|
| nbd-wtf/go-nostr | Nostr SDK para Go | [github.com/nbd-wtf/go-nostr](https://github.com/nbd-wtf/go-nostr) |

### Infraestrutura
| Biblioteca | Descrição | Link |
|---|---|---|
| dgraph-io/badger | Embedded key-value DB | [github.com/dgraph-io/badger](https://github.com/dgraph-io/badger) |
| spf13/cobra | CLI framework | [github.com/spf13/cobra](https://github.com/spf13/cobra) |
| skip2/go-qrcode | QR Code generation | [github.com/skip2/go-qrcode](https://github.com/skip2/go-qrcode) |
| fatih/color | Terminal colors | [github.com/fatih/color](https://github.com/fatih/color) |

---

## 3. Projetos de Referência

| Projeto | O que faz | Relevância | Link |
|---|---|---|---|
| **Breez SDK** | SDK mobile para Lightning (self-custodial) | Referência de LND embutido + LSP | [github.com/AgoraIO-Extensions/breez-sdk-go](https://github.com/AgoraIO-Extensions/breez-sdk-go) |
| **Blixt Wallet** | Carteira Lightning mobile com LND embutido | Prova que LND embutido funciona | [github.com/AgoraIO-Extensions/blixt-wallet](https://github.com/AgoraIO-Extensions/blixt-wallet) |
| **Phoenix (ACINQ)** | Carteira Lightning auto-gerenciada | Referência de UX para onboarding | [phoenix.acinq.co](https://phoenix.acinq.co) |
| **Alby NWC** | Nostr Wallet Connect server | Implementação de NIP-47 | [github.com/getAlby/nostr-wallet-connect](https://github.com/getAlby/nostr-wallet-connect) |
| **crom-p2pfile** | P2P file transfer (Zero-Backend) | Projeto irmão, mesma filosofia | [github.com/MrJc01/crom-p2pfile](https://github.com/MrJc01/crom-p2pfile) |
| **LNbits** | Plataforma Lightning extensível | Referência de LNURL-Pay | [github.com/lnbits/lnbits](https://github.com/lnbits/lnbits) |

---

## 4. Material de Estudo

### Artigos e Guias
- [Mastering the Lightning Network (Book)](https://github.com/lnbook/lnbook) — O'Reilly, gratuito no GitHub
- [Understanding Lightning Network](https://bitcoinmagazine.com/technical/understanding-the-lightning-network-part-building-a-bidirectional-payment-channel-1464710791) — Bitcoin Magazine
- [LNURL Specifications Explained](https://voltage.cloud/blog/lightning-network-faq/what-are-lnurls/) — Voltage
- [Nostr Protocol Overview](https://nostr.com/) — Site oficial

### Vídeos
- [Lightning Network Deep Dive](https://www.youtube.com/watch?v=rrr_zPmEiME) — Whiteboard Crypto
- [Building with LND (Workshop)](https://www.youtube.com/results?search_query=building+with+lnd+go) — Lightning Labs

### Ferramentas de Desenvolvimento
- [Polar](https://lightningpolar.com/) — IDE visual para redes Lightning de teste
- [Bitcoin Core (Regtest)](https://bitcoincore.org/) — Para ambiente de desenvolvimento local
- [ThunderHub](https://thunderhub.io/) — Dashboard web para nós Lightning

---

## 5. Ecossistema CROM (Projetos Relacionados)

| Projeto | Descrição | Stack |
|---|---|---|
| crom-p2pfile | Transferência P2P de arquivos E2EE | React + WebRTC + PeerJS |
| crom-wiki | Wiki colaborativa com IA (Rosa) | Go + TUI |
| crom-dataset | Sistema de datasets descentralizado | React + Go API |
| crom-chat | Chat descentralizado P2P | Go daemon + WebSocket |
| cromia-api | API do ecossistema CROM | Go + REST |
| crom.me | Portal principal CROM | Web |
