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

**Nota:** Saldo real requer integração Neutrino (Milestone 02). Atualmente mostra 0 sats.

---

## `wallet address`

Exibe o endereço Bitcoin de recebimento.

```bash
crom-pay wallet address [--data-dir DIR]
```

**Saída:**
```
⚡ Endereço Bitcoin: tb1qm35ns3zckfjke2588c7zmgq8pn3mkjwf6sg4ya
```

---

## `wallet restore`

Restaura uma carteira a partir de um backup de 24 palavras.

```bash
crom-pay wallet restore [--data-dir DIR] [--network REDE]
```

**Fluxo:**
1. Solicita as 24 palavras (em uma linha)
2. Valida o checksum BIP-39
3. Solicita nova senha
4. Salva a seed cifrada
5. Exibe o endereço (deve ser idêntico ao original)

**Importante:** O mesmo mnemônico sempre gera o mesmo endereço — isso é determinístico por design (BIP-32/84).

---

## Exemplos com Pipe (para scripts)

```bash
# Criar carteira via pipe
printf "minha-senha-forte\nminha-senha-forte\n" | crom-pay wallet create --network testnet

# Ver saldo
printf "minha-senha-forte\n" | crom-pay wallet balance

# Ver endereço
printf "minha-senha-forte\n" | crom-pay wallet address
```

---

← [Anterior: Criptografia](03-criptografia.md) | [Índice](README.md) | [Próximo: Testes →](05-testes.md)
