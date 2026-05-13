# 1. Início

← [Índice](README.md) | [Próximo: Arquitetura →](02-arquitetura.md)

---

## O que é o Crom Bitcoin Pix?

O **Crom Bitcoin Pix** é um sistema de pagamentos Bitcoin que replica a experiência do Pix brasileiro — instantâneo, sem custódia e sem intermediários. É um único binário Go que funciona como:

- **Carteira Bitcoin HD** — Suas chaves, seu dinheiro
- **Nó Lightning** — Pagamentos em segundos (futuro)
- **Identidade Nostr** — Descoberta descentralizada (futuro)

## Por que existe?

O Pix revolucionou pagamentos no Brasil, mas depende do Banco Central e bancos comerciais. O Crom Bitcoin Pix traz a mesma experiência usando Bitcoin:

| Pix (BACEN) | Crom Bitcoin Pix |
|---|---|
| Controlado pelo BC | Soberano (suas chaves) |
| Conta bancária | Apenas um binário |
| KYC obrigatório | Sem identificação |
| Pode ser censurado | Resistente a censura |
| Instantâneo | Instantâneo (Lightning) |

## Quickstart

```bash
# 1. Clonar e compilar
git clone https://github.com/MrJc01/crom-bitcoin-pix.git
cd crom-bitcoin-pix
make build

# 2. Criar carteira
./bin/crom-pay wallet create --network testnet

# 3. Ver endereço
./bin/crom-pay wallet address

# 4. Rodar testes
./tests/run_all.sh
```

## Requisitos

- **Go 1.22+** (testado com 1.25)
- **Linux/macOS** (Windows via WSL)
- ~50MB de espaço em disco

## Princípios

1. **Soberania** — O usuário controla suas chaves
2. **Privacidade** — Sem servidores, sem telemetria
3. **Simplicidade** — Um binário, zero configuração
4. **Segurança** — Criptografia auditada, testes rigorosos

---

← [Índice](README.md) | [Próximo: Arquitetura →](02-arquitetura.md)
