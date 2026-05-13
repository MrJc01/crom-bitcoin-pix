# 5. Testes

← [Anterior: Comandos](04-comandos.md) | [Índice](README.md) | [Próximo: Roadmap →](06-roadmap.md)

---

## Visão Geral

O projeto possui **72 testes automatizados** em **6 camadas**, cobrindo desde unitários até pentest automatizado.

## Como Rodar

```bash
# Suite completa (gera relatório em tests/reports/)
./tests/run_all.sh

# Apenas unitários (rápido, ~10s)
go test ./internal/...

# Com race detector
go test -race ./internal/...

# Apenas segurança
go test -v -tags security ./tests/security/...

# Apenas stress
go test -v -tags stress -timeout 300s ./tests/stress/...
```

## Estrutura

```
tests/
├── run_all.sh              # Runner principal
├── README.md               # Documentação da suite
├── fixtures/
│   └── test_vectors.json   # Vetores BIP-39 conhecidos
├── integration/
│   └── wallet_flow_test.go # Fluxos completos
├── security/
│   ├── entropy_test.go     # Qualidade da entropia
│   ├── crypto_audit_test.go# Auditoria criptográfica
│   └── key_exposure_test.go# Exposição de chaves
├── stress/
│   └── concurrent_test.go  # Carga e concorrência
├── e2e/
│   └── cli_test.sh         # Testes do binário real
└── reports/                # Relatórios gerados (.md + .log)
```

## Camadas de Teste

### Camada 1 — Unitários (Wallet): 24 testes

| Teste | O que valida |
|---|---|
| `TestGenerateMnemonic` | Gera 24 palavras válidas |
| `TestValidateMnemonic_*` | Rejeita mnemônicos inválidos |
| `TestMnemonicToSeed_*` | Seed determinística |
| `TestNewKeychain_*` | Criação e validação de keychain |
| `TestDeriveAddress_*` | Endereços SegWit corretos |
| `TestEncryptDecrypt_*` | Roundtrip AES-GCM |
| `TestWallet_*` | Create, Open, Close, Restore |
| `TestBIP84_OfficialVector` | Vetor oficial BIP-84 ✅ |

### Camada 1b — Unitários (Storage): 12 testes

| Teste | O que valida |
|---|---|
| `TestStore_PutGet` | CRUD básico |
| `TestStore_Delete` | Remoção funciona |
| `TestStore_Overwrite` | Sobrescrita de valores |
| `TestStore_PrefixIsolation` | Prefixos não colidem |
| `TestStore_Persistence` | Dados persistem entre aberturas |
| `TestStore_DirectoryPermissions` | Mode 0700 |
| `TestStore_ClosedStoreErrors` | Erros em store fechado |

### Camada 2 — Integração: 6 testes

| Teste | O que valida |
|---|---|
| `TestFlow_CreateOpenClose` | Ciclo completo |
| `TestFlow_CreateRestoreMatch` | Restore gera mesmo endereço |
| `TestFlow_MultipleAddresses` | 10 endereços únicos |
| `TestFlow_NetworkIsolation` | Mainnet ≠ Testnet |
| `TestFlow_DuplicateCreateFails` | Proteção contra duplicação |
| `TestFlow_DataPersistence` | 5 ciclos open/close |

### Camada 3 — Segurança/Pentest: 15 testes

| Teste | O que valida |
|---|---|
| `TestCrypto_CiphertextDiffers` | Salt/nonce aleatórios |
| `TestCrypto_TamperDetection` | AES-GCM rejeita adulteração |
| `TestCrypto_TruncationDetection` | Dados truncados rejeitados |
| `TestCrypto_WrongPasswordTiming` | Resistência a timing attack |
| `TestCrypto_Argon2idSlowness` | Rate ≤ 1/s |
| `TestEntropy_ShannonEntropy` | 7.98/8.0 bits |
| `TestEntropy_NoBiasInBytes` | Sem viés estatístico |
| `TestKeyExposure_NoPlaintextSeedOnDisk` | Sem seed em plaintext |
| `TestKeyExposure_FilePermissions` | Permissões 700 |
| `TestKeyExposure_PathTraversal` | Resistência a path traversal |

### Camada 4 — Stress: 4 testes

| Teste | O que valida |
|---|---|
| `TestStress_ConcurrentWalletCreation` | 20 wallets simultâneas |
| `TestStress_RapidOpenClose` | 50 ciclos open/close |
| `TestStress_ManyAddressDerivations` | 1000 endereços únicos |
| `TestStress_EncryptDecryptBatch` | 10 ciclos encrypt/decrypt |

### Camada 5 — E2E: 11 testes

| Teste | O que valida |
|---|---|
| `--help` | Flag funciona |
| `--version` | Mostra nome do binário |
| `wallet create` | Gera endereço `tb1q` |
| `wallet balance` | Mostra 0 sats |
| `wallet address` | Endereço consistente |
| Senha errada | Rejeitada com erro |
| Duplicação | Criação duplicada bloqueada |
| `wallet --help` | Lista 4 subcomandos |

## Relatórios

Cada execução do `run_all.sh` gera:

- **`report_YYYYMMDD_HHMMSS.md`** — Relatório formatado com resumo
- **`raw_YYYYMMDD_HHMMSS.log`** — Log bruto completo

---

← [Anterior: Comandos](04-comandos.md) | [Índice](README.md) | [Próximo: Roadmap →](06-roadmap.md)
