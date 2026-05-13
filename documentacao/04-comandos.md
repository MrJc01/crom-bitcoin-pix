# 4. Comandos CLI

← [Anterior: Criptografia](03-criptografia.md) | [Índice](README.md) | [Próximo: Testes →](05-testes.md)

---

## Uso Geral

```bash
crom-pay [comando] [subcomando] [flags]
```

### Flags Globais

| Flag | Padrão | Descrição |
|---|---|---|
| `--data-dir` | `~/.crom-pay` | Diretório de dados da carteira |
| `--network` | `mainnet` | Rede Bitcoin (`mainnet`, `testnet`, `regtest`) |
| `--help` | — | Ajuda do comando |
| `--version` | — | Versão do binário |

---

## `wallet create`

Cria uma nova carteira Bitcoin HD.

```bash
crom-pay wallet create [--data-dir DIR] [--network REDE]
```

**Fluxo:**
1. Solicita senha (mínimo 8 caracteres)
2. Solicita confirmação de senha
3. Gera 24 palavras BIP-39
4. Exibe as palavras na tela (**anote em papel!**)
5. Salva a seed cifrada no disco
6. Exibe o endereço Bitcoin

**Exemplo:**
```
🔐 Criando carteira soberana...

╔════════════════════════════════════════════════════════════╗
║  ⚠️  ANOTE ESTAS PALAVRAS EM PAPEL — NUNCA COMPARTILHE!   ║
╠════════════════════════════════════════════════════════════╣
║   1. abandon      2. ability      3. able         4. about║
║   ...                                                      ║
╚════════════════════════════════════════════════════════════╝

⚡ Endereço Bitcoin: tb1qm35ns3zckfjke2588c7zmgq8pn3mkjwf6sg4ya
🌐 Rede: testnet
✅ Carteira criada com sucesso!
```

---

## `wallet balance`

Consulta o saldo da carteira.

```bash
crom-pay wallet balance [--data-dir DIR]
```

**Nota:** Saldo real é consultado via API mempool.space.

---

## `wallet address`

Exibe o endereço Bitcoin de recebimento.

```bash
crom-pay wallet address [--data-dir DIR]
```

---

## `wallet restore`

Restaura uma carteira a partir de um backup de 24 palavras.

```bash
crom-pay wallet restore [--data-dir DIR] [--network REDE]
```

---

## `receive [amount]`

Gera QR Code para recebimento Bitcoin.

```bash
crom-pay receive                 # QR com endereço
crom-pay receive 50000           # QR com valor (50000 sats)
```

---

## `pay <destino> <amount>`

Envia Bitcoin. Detecta automaticamente o tipo de destino:

```bash
crom-pay pay bc1q... 50000           # On-chain
crom-pay pay lnbc... 50000           # Lightning invoice
crom-pay pay alice@crom.run 5000     # NIP-05 → Lightning
```

---

## `lightning`

Gerencia o nó Lightning Network.

```bash
crom-pay lightning info              # Status do nó
crom-pay lightning invoice 5000      # Gerar invoice
crom-pay lightning channels          # Listar canais
crom-pay lightning setup             # Configurar conexão LND
```

---

## `nostr`

Identidade Nostr descentralizada.

```bash
crom-pay nostr identity              # Exibir npub (derivado do seed)
crom-pay nostr publish "mensagem"    # Publicar nota em relays
crom-pay nostr relays                # Listar relays conectados
crom-pay nostr verify user@crom.run  # Verificar NIP-05
```

---

## `contacts`

Gerenciar contatos de pagamento.

```bash
crom-pay contacts add alice --btc bc1q... --nip05 alice@crom.run
crom-pay contacts list
crom-pay contacts show alice
crom-pay contacts remove alice
```

---

## `tui`

Interface visual interativa no terminal (bubbletea + lipgloss).

```bash
crom-pay tui
```

Telas: Dashboard, Receber, Enviar, Nostr, Configurações.

---

## Exemplos com Pipe (para scripts)

```bash
# Criar carteira via pipe
printf "minha-senha-forte\nminha-senha-forte\n" | crom-pay wallet create --network testnet

# Ver saldo
printf "minha-senha-forte\n" | crom-pay wallet balance

# Receber com QR
printf "minha-senha-forte\n" | crom-pay receive 50000
```

---

← [Anterior: Criptografia](03-criptografia.md) | [Índice](README.md) | [Próximo: Testes →](05-testes.md)
