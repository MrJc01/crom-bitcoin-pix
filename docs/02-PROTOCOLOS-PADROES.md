# 📡 Protocolos e Padrões — Crom Bitcoin Pix

> Especificações técnicas de todos os protocolos utilizados pelo binário.

---

## 1. Camada Bitcoin (On-Chain)

### BIP-39 — Mnemonic Code (Semente)
- **O que é:** Converte entropia criptográfica em 12/24 palavras legíveis.
- **Uso no projeto:** Onboarding do usuário. As palavras são a "senha mestra" da carteira.
- **Spec:** [BIP-39](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)

### BIP-32 — Hierarchical Deterministic Wallets
- **O que é:** Derivação de chaves filhas a partir de uma chave mestre.
- **Uso no projeto:** Uma semente → infinitas chaves (endereços) diferentes.
- **Spec:** [BIP-32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)

### BIP-44 — Multi-Account Hierarchy
- **O que é:** Padrão de caminhos de derivação: `m/purpose'/coin'/account'/change/index`
- **Uso no projeto:** `m/84'/0'/0'/0/0` para o primeiro endereço SegWit nativo.
- **Spec:** [BIP-44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)

### BIP-84 — Native SegWit (Bech32)
- **O que é:** Endereços que começam com `bc1q...` — mais baratos e eficientes.
- **Uso no projeto:** Padrão de endereço para depósitos on-chain.
- **Spec:** [BIP-84](https://github.com/bitcoin/bips/blob/master/bip-0084.mediawiki)

### BIP-157/158 — Compact Block Filters (Neutrino)
- **O que é:** Protocolo de light client que baixa apenas filtros compactos dos blocos.
- **BIP-157:** Define o protocolo P2P para solicitar filtros.
- **BIP-158:** Define a construção dos filtros usando codificação Golomb-Rice.
- **Uso no projeto:** Permite verificar transações sem baixar 500GB da blockchain.
- **Privacidade:** O full node NÃO sabe quais endereços o binário está rastreando.
- **Specs:** [BIP-157](https://github.com/bitcoin/bips/blob/master/bip-0157.mediawiki) | [BIP-158](https://github.com/bitcoin/bips/blob/master/bip-0158.mediawiki)

---

## 2. Camada Lightning (Off-Chain)

### BOLT11 — Invoice Protocol
- **O que é:** Formato padronizado para faturas de pagamento Lightning.
- **Encoding:** Bech32 com prefixo `lnbc` (mainnet) ou `lntb` (testnet).
- **Campos principais:**
  - Payment Hash (hash do preimage)
  - Recipient Public Key (nó destino)
  - Amount (valor em millisatoshis)
  - Expiration Time (validade)
  - Routing Hints (dicas de rota)
  - Signature (assinatura do payee)
- **Uso no projeto:** QR Code gerado pelo receptor contém uma invoice BOLT11.
- **Spec:** [BOLT #11](https://github.com/lightning/bolts/blob/master/11-payment-encoding.md)

### LNURL-Pay (LUD-06)
- **O que é:** Protocolo HTTP que automatiza a troca de invoices.
- **Fluxo:**
  1. Decodificar Lightning Address → URL HTTP
  2. `GET https://domain/.well-known/lnurlp/user` → metadata (min/max sats)
  3. Usuário escolhe valor
  4. `GET callback?amount=X` → recebe invoice BOLT11
  5. Carteira paga a invoice
- **Uso no projeto:** Permite pagar `alice@crom.run` como se fosse Pix por e-mail.
- **Spec:** [LUD-06](https://github.com/lnurl/luds/blob/luds/06.md)

### Lightning Address
- **O que é:** Abstração que mapeia `user@domain` → endpoint LNURL-Pay.
- **Conversão:** `alice@crom.run` → `https://crom.run/.well-known/lnurlp/alice`
- **Uso no projeto:** "Chave Pix" do sistema descentralizado.

### HTLCs (Hash Time-Locked Contracts)
- **O que é:** Contratos que permitem pagamentos multi-hop (A→B→C→D).
- **Segurança:** Fundos são liberados apenas quando o preimage correto é revelado.
- **Timeout:** Se o preimage não é revelado a tempo, os fundos retornam ao pagador.
- **Uso no projeto:** Base de toda transação Lightning.

### SCB (Static Channel Backup)
- **O que é:** Backup mínimo que permite recuperar fundos de canais em caso de perda de dados.
- **Uso no projeto:** Backup automático após cada mudança de estado de canal.

---

## 3. Camada Nostr (Identidade/Descoberta)

### NIP-01 — Basic Protocol
- **O que é:** Define eventos, assinaturas e comunicação com relays.
- **Uso no projeto:** Publicar perfil do usuário e descobrir outros.

### NIP-05 — DNS-Based Verification
- **O que é:** Mapeia `user@domain` para uma chave pública Nostr.
- **Uso no projeto:** Verificação de identidade (`alice@crom.run` → pubkey).

### NIP-47 — Nostr Wallet Connect (NWC)
- **O que é:** Permite controlar uma carteira Lightning remotamente via mensagens Nostr cifradas.
- **Eventos:**
  - Kind 23194: Request (ex: `pay_invoice`)
  - Kind 23195: Response (resultado do pagamento)
- **Criptografia:** NIP-04 ou NIP-44 (E2E encryption)
- **Uso no projeto:** Integração futura com apps externos que queiram disparar pagamentos.

### NIP-57 — Lightning Zaps
- **O que é:** Protocolo para "gorjetas" via Lightning usando Nostr.
- **Eventos:**
  - Kind 9734: Zap Request (intenção de pagamento)
  - Kind 9735: Zap Receipt (confirmação)
- **Uso no projeto:** Integração com ecossistema Nostr existente.

---

## 4. Tabela Resumo dos Protocolos

| Camada | Protocolo | Papel no Sistema |
|---|---|---|
| Bitcoin | BIP-39 | Semente/backup |
| Bitcoin | BIP-32/44/84 | Derivação de chaves |
| Bitcoin | BIP-157/158 | Light client (Neutrino) |
| Lightning | BOLT11 | Faturas de pagamento |
| Lightning | LNURL-Pay | Automação de invoices |
| Lightning | HTLCs | Pagamentos multi-hop |
| Lightning | SCB | Backup de canais |
| Nostr | NIP-01 | Eventos base |
| Nostr | NIP-05 | Identidade DNS |
| Nostr | NIP-47 | Wallet Connect |
| Nostr | NIP-57 | Zaps (gorjetas) |
