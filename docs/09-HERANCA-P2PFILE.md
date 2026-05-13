# 🧬 Herança do crom-p2pfile — Análise de Padrões Reutilizáveis

> Estudo detalhado do projeto crom-p2pfile e como seus padrões migram para o crom-bitcoin-pix.

---

## 1. Visão Geral do crom-p2pfile

**Repositório:** [github.com/MrJc01/crom-p2pfile](https://github.com/MrJc01/crom-p2pfile)

**O que é:** MVP de Web App estático para transferência de arquivos Peer-to-Peer com criptografia End-to-End (E2EE) e arquitetura Zero-Backend.

**Stack:** React 18 + Vite + TypeScript (86.1%) + Tailwind CSS

**Princípios Chave:**
- Soberania Digital
- Local-First
- Zero-Backend / Zero-Database
- E2EE (AES-GCM 256-bit)
- Optical-OOB (QR Code para troca de chaves)

---

## 2. Arquitetura do crom-p2pfile (Como Funciona)

```
┌─────────────────┐                    ┌─────────────────┐
│   SENDER (PC)   │                    │ RECEIVER (Phone) │
│                 │                    │                  │
│ 1. Gera UUID    │                    │                  │
│ 2. Gera AES Key │    ┌──────────┐    │                  │
│ 3. Gera QR Code │───▶│ QR Code  │───▶│ 4. Lê QR (jsQR) │
│    #id=X&key=Y  │    │ (Ótico)  │    │ 5. Extrai key    │
│                 │    └──────────┘    │                  │
│                 │                    │                  │
│ 6. WebRTC ──────┼────────────────────┼──── DataChannel  │
│    PeerJS       │    (P2P direto)    │                  │
│                 │                    │                  │
│ 7. Chunking ────┼──▶ 64KB chunks ───┼──▶ 8. Reassembly │
│    AES-GCM      │    (cifrados)      │    Decrypt       │
│                 │                    │    ObjectURL      │
└─────────────────┘                    └──────────────────┘

Segurança:
- Fragmento #id=X&key=Y NUNCA sai do navegador (hash fragment)
- Chave AES trafega apenas pelo canal ótico (câmera → tela)
- PeerJS faz APENAS STUN discovery (NAT traversal)
- ZERO servidores de relay (TURN) propositalmente
```

---

## 3. Padrões que Migram para o crom-bitcoin-pix

### 3.1 Troca de Segredos via Canal Ótico (QR Code)

| crom-p2pfile | crom-bitcoin-pix |
|---|---|
| QR contém `#id=X&key=Y` (UUID + AES key) | QR contém Invoice BOLT11 (hash + destino + valor) |
| Segredo: chave AES-GCM 256-bit | Segredo: Payment Hash + Routing Hints |
| Leitor: jsQR (câmera do celular) | Leitor: QR scanner do CLI ou app |
| Propósito: estabelecer túnel cifrado | Propósito: iniciar pagamento Lightning |

**Lição:** O QR Code é o "handshake visual" entre dois peers. No P2P File é para cifrar dados. No Bitcoin Pix é para transmitir a fatura de pagamento.

### 3.2 Zero-Backend (Sem Servidor Central)

| crom-p2pfile | crom-bitcoin-pix |
|---|---|
| PeerJS broker apenas para STUN | Neutrino peers apenas para block filters |
| Sem banco de dados server-side | BadgerDB local (client-side) |
| Se o site cai, conexões ativas continuam | Se relays caem, pagamentos LN continuam |

**Lição:** O backend é apenas um "catalisador" para a conexão inicial. Depois, os peers são autônomos.

### 3.3 Chunking e Streams

| crom-p2pfile | crom-bitcoin-pix |
|---|---|
| Arquivo dividido em chunks de 64KB | Pagamento pode ser dividido via MPP (Multi-Path Payments) |
| Cada chunk cifrado individualmente | Cada HTLC é atomicamente seguro |
| Reassembly no receptor | Preimage revela o pagamento completo |

**Lição:** Fragmentar o payload para não sobrecarregar o canal é universal — seja WebRTC DataChannel ou Lightning channel capacity.

### 3.4 Criptografia Local-First

| crom-p2pfile | crom-bitcoin-pix |
|---|---|
| `Window.crypto.subtle` (Web Crypto API) | `crypto/rand` + `btcec` (Go stdlib + btcsuite) |
| AES-GCM 256-bit para dados | ECDSA secp256k1 para assinaturas Bitcoin |
| Chave nunca toca servidor HTTP | Chave privada nunca sai do disco local |
| Entropy: `crypto.getRandomValues()` | Entropy: `crypto/rand.Read()` |

**Lição:** A geração de entropia é a raiz de toda segurança. Ambos os projetos usam CSPRNG do sistema operacional.

---

## 4. O que NÃO Migra (Diferenças Fundamentais)

| Aspecto | crom-p2pfile | crom-bitcoin-pix |
|---|---|---|
| **Linguagem** | TypeScript (browser) | Go (binário nativo) |
| **Runtime** | Browser (V8/SpiderMonkey) | OS nativo (Go runtime) |
| **Persistência** | Nenhuma (sessão efêmera) | BadgerDB (persistente) |
| **Rede** | WebRTC (browser P2P) | TCP/IP Lightning (LN P2P) |
| **Valor** | Arquivos (bytes) | Bitcoin (satoshis) |
| **Consequência de bug** | Arquivo corrompido | Perda financeira |
| **Disponibilidade** | Ambos online simultaneamente | Receptor pode estar offline (async) |

> **Diferença crítica:** No P2P File, um bug causa "arquivo corrompido, tenta de novo". No Bitcoin Pix, um bug pode causar **perda irreversível de fundos**. O nível de rigor nos testes e na gestão de estado é ordens de magnitude maior.

---

## 5. Padrões de Código do crom-p2pfile para Reusar

### 5.1 Estrutura de Contexto (React → Go interfaces)

No crom-p2pfile, `PeerContext.tsx` centraliza toda lógica de conexão P2P.

**Equivalente Go:**
```go
// internal/lightning/node.go
type LightningNode interface {
    Start() error
    Stop() error
    CreateInvoice(amount int64, memo string) (*Invoice, error)
    PayInvoice(bolt11 string) (*PaymentResult, error)
    GetBalance() (*Balance, error)
    OpenChannel(peer string, amount int64) (*ChannelPoint, error)
}
```

### 5.2 Event-Driven Architecture

No crom-p2pfile, eventos WebRTC (`onopen`, `ondata`, `onclose`) disparam callbacks.

**Equivalente Go:**
```go
// Channels (Go idiom) para eventos assíncronos
type NodeEvents struct {
    PaymentReceived chan *Payment
    ChannelOpened   chan *Channel
    ChannelClosed   chan *Channel
    InvoicePaid     chan *Invoice
    PeerConnected   chan *Peer
}
```

### 5.3 Padrão de Fragmentação

No crom-p2pfile, chunking de 64KB com controle de fluxo.

**Equivalente Go:** Já embutido no LND (HTLCs com flow control nativo).

---

## 6. Conclusão: Evolução Natural

```
crom-p2pfile (2024)          crom-bitcoin-pix (2026)
─────────────────            ──────────────────────
Web App Estático      ──▶    Binário Go Nativo
Arquivos (bytes)      ──▶    Bitcoin (satoshis)
WebRTC Efêmero        ──▶    Lightning Persistente
AES-GCM Session       ──▶    BIP-39 Permanent Keys
QR = Chave Cripto     ──▶    QR = Invoice BOLT11
PeerJS Discovery      ──▶    Nostr Discovery
```

O crom-bitcoin-pix é a **versão com consequências financeiras** do mesmo DNA de soberania digital. A complexidade aumenta porque dinheiro é irreversível, mas os princípios (zero-backend, E2EE, P2P, QR Code) são idênticos.
