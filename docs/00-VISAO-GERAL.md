# 🔥 Crom Bitcoin Pix — Visão Geral do Projeto

> **Missão:** Construir um binário Go auto-contido que transforma Bitcoin na experiência de pagamento instantâneo do Pix — sem servidores, sem custódia, sem intermediários.

---

## 📋 Sumário Executivo

O **Crom Bitcoin Pix** é um sistema de pagamentos descentralizado, **local-first** e **self-custodial**, escrito em Go. O usuário baixa um único binário executável que funciona simultaneamente como:

1. **Carteira Bitcoin** (chaves derivadas localmente via BIP-39/44)
2. **Nó Lightning** (pagamentos instantâneos via LND embutido + Neutrino)
3. **Interface de Descoberta** (identidades publicadas via protocolo Nostr)

O projeto é uma evolução natural do ecossistema CROM, aplicando os mesmos princípios de **soberania digital** e **zero-backend** já demonstrados no [crom-p2pfile](https://github.com/MrJc01/crom-p2pfile).

---

## 🎯 Princípios Fundamentais

| Princípio | Significado |
|---|---|
| **Self-Custodial** | O usuário é dono das suas chaves. A semente (12/24 palavras) nunca sai do dispositivo local. |
| **Zero-Backend** | Nenhum servidor central. O binário fala direto com a rede Bitcoin e relays Nostr. |
| **Local-First** | Todos os dados (saldos, histórico, canais) ficam no disco local em DB embutido. |
| **Binário Único** | Um executável compilado estaticamente para Linux/Windows/macOS. Sem Docker, sem dependências. |
| **Experiência Pix** | QR Code → Pagamento instantâneo em <1 segundo. Simples como o Pix do Banco Central. |

---

## 🧬 DNA do Projeto (Herança do crom-p2pfile)

O `crom-p2pfile` provou que é possível construir sistemas **E2EE, Zero-Backend** e **peer-to-peer** puramente no lado do cliente. As lições que migramos:

| crom-p2pfile (Web/P2P) | crom-bitcoin-pix (Go/Lightning) |
|---|---|
| WebRTC DataChannel | Lightning Network Payment Channels |
| QR Code com chave AES-GCM | QR Code com Invoice BOLT11 |
| PeerJS (STUN discovery) | Neutrino (Bitcoin SPV discovery) |
| Window.crypto.subtle | BIP-32/39/44 key derivation |
| Chunking 64KB (arquivos) | HTLCs (pagamentos multi-hop) |
| Zero-Backend Web | Zero-Backend Binary |
| Optical-OOB (câmera) | Nostr Relay Discovery |

---

## 🏗️ Arquitetura de Alto Nível

```
┌─────────────────────────────────────────────────┐
│                CROM-PAY BINARY                  │
│                                                 │
│  ┌─────────┐  ┌──────────┐  ┌───────────────┐  │
│  │ Wallet  │  │ Lightning│  │    Nostr       │  │
│  │ Core    │  │ Engine   │  │   Discovery    │  │
│  │         │  │          │  │               │  │
│  │ BIP-39  │  │ LND/     │  │ go-nostr     │  │
│  │ BIP-44  │  │ Neutrino │  │ NIP-47/57    │  │
│  │ btcutil │  │ Channels │  │ Relay Mgmt   │  │
│  └────┬────┘  └────┬─────┘  └──────┬────────┘  │
│       │            │               │            │
│  ┌────┴────────────┴───────────────┴─────────┐  │
│  │           Local Storage (BadgerDB)         │  │
│  │    Canais · Histórico · Contatos · Config  │  │
│  └────────────────────────────────────────────┘  │
│                                                 │
│  ┌────────────────────────────────────────────┐  │
│  │      Interface Layer (CLI + REST API)      │  │
│  │   cobra CLI  ·  localhost:8420 REST API    │  │
│  └────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
         │              │              │
    ┌────┴────┐   ┌─────┴─────┐  ┌────┴──────┐
    │ Bitcoin │   │ Lightning │  │   Nostr   │
    │ Network │   │  Network  │  │  Relays   │
    │ (P2P)   │   │   (P2P)   │  │  (P2P)    │
    └─────────┘   └───────────┘  └───────────┘
```

---

## 📂 Índice da Documentação

| # | Documento | Conteúdo |
|---|---|---|
| 00 | [VISÃO GERAL](./00-VISAO-GERAL.md) | Este documento. Missão, princípios, arquitetura macro. |
| 01 | [ARQUITETURA TÉCNICA](./01-ARQUITETURA-TECNICA.md) | Stack tecnológico, componentes internos, fluxos de dados. |
| 02 | [PROTOCOLOS E PADRÕES](./02-PROTOCOLOS-PADROES.md) | BIP-39, BIP-44, BOLT11, LNURL, Neutrino, Nostr NIPs. |
| 03 | [FLUXOS DE USUÁRIO](./03-FLUXOS-USUARIO.md) | Onboarding, enviar, receber, backup, restauração. |
| 04 | [ESTRUTURA DO PROJETO](./04-ESTRUTURA-PROJETO.md) | Árvore de diretórios Go, packages, dependências. |
| 05 | [ROADMAP E MILESTONES](./05-ROADMAP-MILESTONES.md) | Fases de desenvolvimento com critérios de aceitação. |
| 06 | [RISCOS E MITIGAÇÕES](./06-RISCOS-MITIGACOES.md) | Análise técnica de viabilidade e armadilhas. |
| 07 | [GLOSSÁRIO](./07-GLOSSARIO.md) | Terminologia Bitcoin, Lightning e Nostr. |
| 08 | [REFERÊNCIAS](./08-REFERENCIAS.md) | Links para specs, libs, e projetos de referência. |
| 09 | [HERANÇA P2PFILE](./09-HERANCA-P2PFILE.md) | DNA do crom-p2pfile: padrões reutilizáveis e diferenças. |

---

## 🚀 Como Começar

```bash
# Clone o repositório
git clone https://github.com/MrJc01/crom-bitcoin-pix.git
cd crom-bitcoin-pix

# Leia a documentação na ordem
ls docs/
```

> **Próximo passo:** Leia o [01-ARQUITETURA-TECNICA.md](./01-ARQUITETURA-TECNICA.md) para entender como cada peça se encaixa.
