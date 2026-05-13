# 🧑‍💻 Fluxos de Usuário — Crom Bitcoin Pix

> Jornadas completas do usuário com o binário, desde onboarding até pagamentos.

---

## 1. Onboarding (Primeira Execução)

```
$ crom-pay wallet create

🔐 Gerando carteira soberana...

╔══════════════════════════════════════════════════════════╗
║  SUA SEMENTE (ANOTE EM PAPEL — NUNCA COMPARTILHE!)      ║
║                                                          ║
║  1. abandon   7. diamond  13. music    19. vapor         ║
║  2. ability   8. effort   14. nation   20. village       ║
║  3. achieve   9. faith    15. ocean    21. voice         ║
║  4. balance  10. glory    16. perfect  22. weather       ║
║  5. captain  11. harbor   17. quality  23. wonder        ║
║  6. crystal  12. ivory    18. rocket   24. youth         ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝

⚡ Endereço Bitcoin (depósito): bc1q...xyz
🌐 Identidade Nostr: npub1...abc
📧 Lightning Address: user123@crom.run

✅ Carteira criada com sucesso!
💡 Envie Bitcoin para o endereço acima para começar a usar.
```

**O que acontece por baixo:**
1. `crypto/rand` gera 256 bits de entropia
2. BIP-39 converte em 24 palavras
3. BIP-32 deriva a master key
4. BIP-84 gera o primeiro endereço SegWit
5. Nostr keypair é derivada da mesma semente
6. BadgerDB inicializa o armazenamento local
7. Neutrino inicia sync com a rede Bitcoin

---

## 2. Restauração de Carteira

```
$ crom-pay wallet restore

🔑 Digite sua semente (24 palavras separadas por espaço):
> abandon ability achieve balance captain crystal diamond effort faith glory harbor ivory music nation ocean perfect quality rocket vapor village voice weather wonder youth

🔄 Sincronizando com a rede Bitcoin via Neutrino...
   ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 100%

✅ Carteira restaurada!
   Saldo: 0.00150000 BTC (150,000 sats)
   Canais Lightning: 1 ativo
```

---

## 3. Enviar Pagamento (Pagar como Pix)

### 3.1 Via Lightning Address (estilo e-mail)

```
$ crom-pay pay alice@crom.run 5000

🔍 Buscando alice@crom.run nos relays Nostr...
   → Encontrado! Pubkey: npub1alice...

💰 Valor: 5,000 sats (~R$ 25,00)
📋 Destino: alice@crom.run

Confirmar pagamento? [s/N]: s

⚡ Roteando pagamento via Lightning Network...
✅ Pago! Preimage: 0x3f...9a
   Tempo: 0.3 segundos
   Taxa de rota: 1 sat
```

### 3.2 Via QR Code (Invoice BOLT11)

```
$ crom-pay pay --qr

📷 Escaneie o QR Code ou cole a invoice:
> lnbc50n1pj...

💰 Valor: 5,000 sats
📋 Destino: 03abc...def (Nó Lightning)
📝 Descrição: "Café na padaria"
⏰ Expira em: 10 minutos

Confirmar pagamento? [s/N]: s

⚡ Roteando...
✅ Pago em 0.2 segundos!
```

---

## 4. Receber Pagamento (Cobrar como Pix)

```
$ crom-pay receive 10000

📥 Gerando cobrança de 10,000 sats...

╔═══════════════════════════════════╗
║          ▄▄▄▄▄▄▄▄▄▄▄▄▄           ║
║          █ QR CODE BOLT11 █       ║
║          █               █       ║
║          █  (escaneável)  █       ║
║          █               █       ║
║          ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀        ║
╚═══════════════════════════════════╝

Invoice: lnbc100n1pj...
Expira em: 1 hora

⏳ Aguardando pagamento...
✅ Recebido! +10,000 sats
   De: bob@crom.run
```

---

## 5. Consultar Saldo

```
$ crom-pay wallet balance

╔═══════════════════════════════════════╗
║        💰 SALDO CROM-PAY             ║
╠═══════════════════════════════════════╣
║                                       ║
║  On-Chain:     0.00500000 BTC         ║
║  Lightning:    0.00150000 BTC         ║
║  ─────────────────────────────        ║
║  Total:        0.00650000 BTC         ║
║                (~R$ 3.250)            ║
║                                       ║
║  Canais Ativos: 2                     ║
║  Capacidade Total: 0.01 BTC           ║
╚═══════════════════════════════════════╝
```

---

## 6. Gerenciar Canais Lightning

```
$ crom-pay channels list

┌────┬──────────────────┬────────────┬──────────┬────────┐
│ #  │ Peer             │ Capacidade │ Local    │ Status │
├────┼──────────────────┼────────────┼──────────┼────────┤
│ 1  │ ACINQ            │ 500k sats  │ 350k     │ ✅     │
│ 2  │ Bitrefill        │ 200k sats  │ 50k      │ ✅     │
└────┴──────────────────┴────────────┴──────────┴────────┘

$ crom-pay channels open --peer 03abc...def --amount 100000
⚡ Abrindo canal de 100,000 sats com 03abc...def
🔗 TX on-chain: txid:abc123...
⏳ Aguardando 3 confirmações...
```

---

## 7. Modo Daemon (Nó Persistente)

```
$ crom-pay daemon start

🚀 Iniciando Crom-Pay Daemon...
   ├── Neutrino sync: ✅ (block 850,000)
   ├── Lightning node: ✅ (2 canais ativos)
   ├── Nostr relays: ✅ (3 conectados)
   └── REST API: ✅ (localhost:8420)

📡 Nó online e pronto para receber pagamentos.
   Lightning Address: user@crom.run
   Pubkey: 03abc...

[Ctrl+C para parar]
```

---

## 8. Backup

```
$ crom-pay wallet backup --output ~/backup-crom.enc

🔐 Criando backup cifrado...
   ├── Semente (cifrada com senha)
   ├── Estado dos canais (SCB)
   └── Lista de contatos

🔑 Digite uma senha para o backup:
> ********

✅ Backup salvo em: ~/backup-crom.enc
   Tamanho: 12 KB
```
