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
| Auditoria de segurança (6 vulns) | ✅ |

---

### ✅ Milestone 02 — Lightning Engine (Completo)

**Objetivo:** Pagamentos instantâneos via Lightning Network.

| Entrega | Status |
|---|---|
| Saldo real via mempool.space API | ✅ |
| TX Builder P2WPKH (on-chain) | ✅ |
| LND REST Client (macaroon auth) | ✅ |
| Invoices BOLT-11 (criar/decodificar) | ✅ |
| `pay` command (BTC/LN/NIP-05) | ✅ |
| `receive` command (QR Code) | ✅ |
| `lightning info/invoice/channels/setup` | ✅ |

---

### ✅ Milestone 03 — Nostr Identity (Completo)

**Objetivo:** Identidade descentralizada para descoberta.

| Entrega | Status |
|---|---|
| Derivação NIP-06 (m/44'/1237'/0'/0/0) | ✅ |
| Relay pool com fallback | ✅ |
| Publicação de eventos (kind 0, 1) | ✅ |
| Verificação NIP-05 | ✅ |
| Zap request (kind 9734) | ✅ |
| `nostr identity/publish/relays/verify` | ✅ |

---

### ✅ Milestone 04 — Pix UX (Completo)

**Objetivo:** Experiência de pagamento igual ao Pix.

| Entrega | Status |
|---|---|
| QR Code ASCII (normal + invertido) | ✅ |
| URI BIP-21 | ✅ |
| TUI bubbletea + lipgloss (5 telas) | ✅ |
| Contatos (CRUD BadgerDB) | ✅ |
| Resolução automática de destino | ✅ |
| `contacts add/list/show/remove` | ✅ |
| `tui` dashboard interativo | ✅ |

---

### ✅ Milestone 05 — CI/CD (Completo)

**Objetivo:** Automação e distribuição.

| Entrega | Status |
|---|---|
| GitHub Actions CI (test + vet) | ✅ |
| Build matrix 5 plataformas | ✅ |
| Release automático com SHA256 | ✅ |
| 79 testes em 6 camadas | ✅ |

---

### ⬜ Milestone 06 — Multi-plataforma (Futuro)

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
2026 Q2: Wallet Core ✅ + Lightning ✅ + Nostr ✅ + Pix UX ✅ + CI/CD ✅
2026 Q3: Mobile + WASM + Desktop
2026 Q4: Integração P2P completa via crom.run
```

---

← [Anterior: Testes](05-testes.md) | [Índice](README.md) | [Próximo: Glossário →](07-glossario.md)
