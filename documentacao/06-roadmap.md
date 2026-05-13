# 6. Roadmap

← [Anterior: Testes](05-testes.md) | [Índice](README.md) | [Próximo: Glossário →](07-glossario.md)

---

## Milestones

### ✅ Milestone 01 — Wallet Core (Completo)

**Objetivo:** Carteira HD funcional com criptografia auditada.

| Entrega | Status |
|---|---|
| BIP-39 mnemônico (24 palavras) | ✅ |
| BIP-84 endereços SegWit | ✅ |
| AES-256-GCM + Argon2id | ✅ |
| CLI (create/balance/address/restore) | ✅ |
| BadgerDB storage | ✅ |
| 72 testes em 6 camadas | ✅ |
| Auditoria de segurança | ✅ |
| Documentação completa | ✅ |

---

### 🔜 Milestone 02 — Lightning Engine

**Objetivo:** Pagamentos instantâneos via Lightning Network.

| Entrega | Descrição |
|---|---|
| LND embutido | Nó Lightning dentro do binário |
| Neutrino | Sincronização leve (sem full node) |
| Invoices BOLT-11 | Gerar e pagar invoices |
| `pay` command | Pagar endereço Lightning |
| `receive` command | Gerar invoice para receber |
| Saldo real | Consultar via Neutrino |

---

### ⬜ Milestone 03 — Nostr Identity

**Objetivo:** Identidade descentralizada para descoberta.

| Entrega | Descrição |
|---|---|
| Keypair Nostr | Gerar chaves Nostr do mesmo seed |
| NIP-05 | Verificação de identidade |
| NIP-57 | Zaps (gorjetas Lightning) |
| `alice@crom.run` | Endereços legíveis |
| Relay discovery | Encontrar peers via relays |

---

### ⬜ Milestone 04 — Pix UX

**Objetivo:** Experiência de pagamento igual ao Pix.

| Entrega | Descrição |
|---|---|
| QR Code | Gerar e ler QR de pagamento |
| TUI | Interface visual no terminal |
| 1-tap pay | Pagamento com confirmação única |
| Contatos | Lista de endereços salvos |
| Histórico | Transações locais |

---

### ⬜ Milestone 05 — Multi-plataforma

**Objetivo:** Funcionar em qualquer lugar.

| Entrega | Descrição |
|---|---|
| Mobile (Gomobile) | Android/iOS via bindings |
| Web (WASM) | Interface web |
| Desktop | GUI nativa |
| Docker | Container para servidor |

---

## Visão de Longo Prazo

```
2026 Q2: Wallet Core ✅
2026 Q3: Lightning + Neutrino
2026 Q4: Nostr + Identidade
2027 Q1: Pix UX + QR + TUI
2027 Q2: Mobile + Multi-plataforma
```

---

← [Anterior: Testes](05-testes.md) | [Índice](README.md) | [Próximo: Glossário →](07-glossario.md)
