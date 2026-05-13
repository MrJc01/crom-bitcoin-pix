# 📖 Glossário — Crom Bitcoin Pix

> Terminologia técnica de Bitcoin, Lightning Network e Nostr usada neste projeto.

---

## Bitcoin (Camada Base)

| Termo | Definição |
|---|---|
| **Satoshi (sat)** | Menor unidade de Bitcoin. 1 BTC = 100.000.000 sats. |
| **UTXO** | Unspent Transaction Output. "Moeda" não gasta na blockchain. Cada transação consome UTXOs e gera novos. |
| **SegWit** | Segregated Witness. Upgrade que separa assinaturas dos dados da transação, reduzindo taxas. |
| **Bech32** | Formato de endereço SegWit nativo. Começa com `bc1q...` (mainnet) ou `tb1q...` (testnet). |
| **BIP** | Bitcoin Improvement Proposal. Documento formal para propor mudanças no protocolo Bitcoin. |
| **HD Wallet** | Hierarchical Deterministic Wallet. Uma semente gera infinitas chaves de forma determinística. |
| **Mnemônico** | As 12 ou 24 palavras que codificam a entropia da semente (BIP-39). |
| **Derivation Path** | Caminho hierárquico para gerar chaves: `m/84'/0'/0'/0/0`. |
| **Regtest** | Regression Test Network. Rede Bitcoin local para desenvolvimento (sem valor real). |
| **Testnet** | Rede Bitcoin pública de testes (moedas sem valor real). |
| **Mainnet** | Rede Bitcoin principal (moedas com valor real). |
| **Mempool** | "Fila de espera" de transações Bitcoin aguardando confirmação por mineradores. |
| **Full Node** | Software que baixa e valida toda a blockchain (~500GB+). |
| **Light Client** | Software que verifica transações sem baixar a blockchain completa. |
| **SPV** | Simplified Payment Verification. Técnica de verificação leve (BIP-37, substituída por Neutrino). |
| **Neutrino** | Light client moderno (BIP-157/158) que usa compact block filters. Preserva privacidade. |
| **Compact Block Filter** | Filtro compacto (Golomb-Rice) que resume quais scripts existem em cada bloco. |

---

## Lightning Network (Camada 2)

| Termo | Definição |
|---|---|
| **Canal (Channel)** | Conexão de pagamento bidirecional entre dois nós Lightning. Aberto com TX on-chain. |
| **Invoice** | Fatura de pagamento Lightning no formato BOLT11. Contém hash, valor, destino e validade. |
| **BOLT** | Basis of Lightning Technology. Especificações técnicas da Lightning Network. |
| **BOLT11** | Padrão de encoding de invoices Lightning (bech32, campos tagged). |
| **HTLC** | Hash Time-Locked Contract. Contrato que roteia pagamento multi-hop com segurança atômica. |
| **Preimage** | Segredo criptográfico. Quando revelado, prova que um pagamento foi liquidado. |
| **Payment Hash** | Hash SHA-256 do preimage. Incluído na invoice para verificação. |
| **Routing** | Processo de encontrar caminho de canais entre pagador e recebedor na rede Lightning. |
| **Hop** | Cada nó intermediário no caminho de roteamento de um pagamento. |
| **Inbound Liquidity** | Capacidade de receber pagamentos em um canal. Depende do saldo remoto. |
| **Outbound Liquidity** | Capacidade de enviar pagamentos em um canal. Depende do saldo local. |
| **LSP** | Lightning Service Provider. Serviço que facilita liquidez e abertura de canais para usuários finais. |
| **SCB** | Static Channel Backup. Backup mínimo para recuperar fundos de canais. |
| **Watchtower** | Serviço terceiro que monitora a blockchain para detectar tentativas de fraude em canais. |
| **Breach Remedy** | Transação de penalidade executada quando um peer tenta publicar estado antigo de canal. |
| **LND** | Lightning Network Daemon. Implementação Go do protocolo Lightning (~70% da rede). |
| **CLN** | Core Lightning. Implementação C do protocolo Lightning (antes c-lightning). |
| **LDK** | Lightning Development Kit. Biblioteca Rust modular para embutir Lightning (sem bindings Go oficiais). |

---

## LNURL (Camada de Aplicação)

| Termo | Definição |
|---|---|
| **LNURL** | Lightning Network URL. Protocolo HTTP que automatiza interações Lightning. |
| **LNURL-Pay** | Sub-protocolo para solicitar pagamentos. Resolve endereço → invoice automaticamente. |
| **Lightning Address** | Formato `user@domain` que mapeia para endpoint LNURL-Pay (`/.well-known/lnurlp/user`). |
| **LUD** | LNURL Document. Especificações individuais do protocolo LNURL (similar aos BIPs). |
| **Callback URL** | URL retornada pelo servidor LNURL onde a carteira solicita a invoice final. |
| **millisatoshi (msat)** | 1/1000 de um satoshi. Unidade usada internamente na Lightning Network. |

---

## Nostr (Camada de Identidade)

| Termo | Definição |
|---|---|
| **Nostr** | Notes and Other Stuff Transmitted by Relays. Protocolo de comunicação descentralizado. |
| **Relay** | Servidor que armazena e retransmite eventos Nostr. Qualquer um pode rodar um relay. |
| **Event** | Unidade fundamental do Nostr. JSON assinado com tipo (kind), conteúdo e tags. |
| **Kind** | Tipo numérico de um evento Nostr (ex: 0=metadata, 1=text note, 9734=zap request). |
| **NIP** | Nostr Implementation Possibility. Especificações do protocolo (similar aos BIPs). |
| **npub** | Chave pública Nostr em formato bech32 legível (começa com `npub1...`). |
| **nsec** | Chave privada Nostr em formato bech32 (começa com `nsec1...`). NUNCA compartilhar. |
| **NIP-01** | Protocolo base: eventos, assinaturas, comunicação com relays. |
| **NIP-04** | Mensagens diretas cifradas (shared secret / DH). |
| **NIP-05** | Verificação de identidade via DNS (`user@domain` → pubkey). |
| **NIP-44** | Criptografia de mensagens versão 2 (substitui NIP-04, mais seguro). |
| **NIP-47** | Nostr Wallet Connect (NWC). Controle remoto de carteira via eventos cifrados. |
| **NIP-57** | Lightning Zaps. Gorjetas Lightning integradas ao Nostr. |
| **Zap** | Pagamento Lightning enviado via Nostr (gorjeta/doação). |
| **Zap Request** | Evento Kind 9734 enviado ao LNURL endpoint para iniciar um Zap. |
| **Zap Receipt** | Evento Kind 9735 publicado como prova de pagamento concluído. |

---

## Ecossistema CROM

| Termo | Definição |
|---|---|
| **Local-First** | Filosofia onde dados ficam primariamente no dispositivo do usuário, não em servidores. |
| **Zero-Backend** | Arquitetura sem servidor central. O binário é autossuficiente. |
| **Self-Custodial** | O usuário controla suas próprias chaves. Nenhum terceiro tem acesso aos fundos. |
| **E2EE** | End-to-End Encryption. Dados cifrados de ponta a ponta (só emissor e receptor leem). |
| **Soberania Digital** | Princípio de que o usuário deve ter controle total sobre seus dados e ativos. |
| **crom-p2pfile** | Projeto irmão: transferência P2P de arquivos via WebRTC com E2EE e QR Code. |
| **crom-pay** | Nome do binário executável deste projeto (Crom Bitcoin Pix). |
