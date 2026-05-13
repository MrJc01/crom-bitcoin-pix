# ⚠️ Riscos e Mitigações — Crom Bitcoin Pix

> Análise técnica de viabilidade, armadilhas conhecidas e estratégias de mitigação.

---

## Painel de Especialistas Simulados

| # | Especialista | Contribuição Principal |
|---|---|---|
| 1 | **Eng. de Segurança Criptográfica** | Gestão de chaves privadas, side-channel attacks, entropy sources |
| 2 | **Arquiteto Lightning Network** | Channel management, routing fees, liquidity balancing |
| 3 | **Eng. de Sistemas Distribuídos** | Consistência de estado, crash recovery, data durability |
| 4 | **Especialista em UX Fintech** | Onboarding friction, error handling, regulatory compliance |
| 5 | **SRE / DevOps** | Build pipeline, cross-compilation, monitoring |

---

## Riscos Técnicos

### 🔴 CRÍTICO: Perda de Estado de Canais

**Risco:** Se o estado dos canais Lightning for corrompido (crash, disco cheio, kill -9), o usuário pode perder fundos.

**Mitigação:**
- Implementar SCB (Static Channel Backup) com flush automático a cada mudança
- Suporte a Watchtowers (terceiros monitoram breach attempts)
- Graceful shutdown handler com `signal.Notify(SIGTERM, SIGINT)`
- Testes de crash recovery com `kill -9` durante operações

**Severidade:** 🔴 CRÍTICA — Pode causar perda financeira irreversível.

---

### 🔴 CRÍTICO: Segurança da Semente

**Risco:** Semente armazenada em texto plano no disco pode ser roubada por malware.

**Mitigação:**
- Semente NUNCA é armazenada em texto plano
- Cifrar com AES-256-GCM usando senha do usuário
- Derivar chave de cifra via Argon2id (anti-brute-force)
- Zerar memória após uso (`crypto/subtle.ConstantTimeCompare`)
- Permissões de arquivo: `chmod 600`

**Severidade:** 🔴 CRÍTICA — Roubo de semente = roubo total de fundos.

---

### 🟡 ALTO: Liquidez Inicial (Chicken-and-Egg)

**Risco:** Para receber Lightning, o usuário precisa de inbound liquidity. Para ter inbound, alguém precisa abrir canal com ele. Usuário novo = zero liquidez.

**Mitigação:**
- Integrar com LSPs (Lightning Service Providers) como:
  - Voltage LSP
  - Breez SDK
  - ACINQ Phoenix-style
- LSP abre canal automaticamente no primeiro recebimento
- Custo: ~1-2% do valor do primeiro recebimento (transparente)

**Severidade:** 🟡 ALTA — Principal barreira de UX para novos usuários.

---

### 🟡 ALTO: Requisito de Estar Online

**Risco:** Lightning Network exige que os nós estejam online para receber. Se o binário estiver fechado, pagamentos falham.

**Mitigação:**
- Modo Daemon (`crom-pay daemon start`) para rodar em background
- Investigar Async Payments (proposta em desenvolvimento no protocolo LN)
- Fallback: exibir endereço on-chain quando offline
- Documentar claramente essa limitação

**Severidade:** 🟡 ALTA — Incompatível com UX "sempre disponível" do Pix.

---

### 🟡 ALTO: Complexidade do LND Embutido

**Risco:** Embutir o LND completo como biblioteca Go é complexo. O LND foi projetado como daemon standalone, não como library.

**Mitigação:**
- Estudar projetos que já fazem isso (Breez SDK, Blixt Wallet)
- Alternativa: rodar LND como subprocess gerenciado pelo binário
- Usar a API gRPC do LND via unix socket local
- Se inviável: considerar usar lnd como sidecar process

**Severidade:** 🟡 ALTA — Pode forçar mudança arquitetural significativa.

---

### 🟢 MÉDIO: Tamanho do Binário

**Risco:** Com LND + Neutrino + BadgerDB + go-nostr, o binário pode ficar grande (>50MB).

**Mitigação:**
- `ldflags="-s -w"` para strip debug symbols
- UPX compression se necessário
- Aceitar tamanho: 30-50MB é aceitável para uma carteira completa

**Severidade:** 🟢 MÉDIA — Inconveniência, não blocker.

---

### 🟢 MÉDIO: Taxas On-Chain para Abrir Canais

**Risco:** Cada abertura de canal requer uma transação on-chain (taxa de ~$1-10 dependendo do mempool).

**Mitigação:**
- Batch channel opens quando possível
- LSP pode absorver custo no primeiro canal
- Educar usuário sobre on-chain vs off-chain

**Severidade:** 🟢 MÉDIA — Custo de onboarding.

---

### 🟢 MÉDIO: Nostr Relay Availability

**Risco:** Se os relays Nostr que o usuário usa ficarem offline, a descoberta de endereços falha.

**Mitigação:**
- Conectar a múltiplos relays (3-5) simultaneamente
- Cache local de endereços conhecidos
- Fallback: permitir pagamento via invoice BOLT11 direta

**Severidade:** 🟢 MÉDIA — Degradação graciosa, não falha total.

---

## Análise de Viabilidade

| Aspecto | Viabilidade | Notas |
|---|---|---|
| **Carteira BIP-39/84** | ✅ Alta | btcsuite é maduro e bem testado |
| **Neutrino Light Client** | ✅ Alta | Usado em produção pelo LND mobile |
| **LND Embutido** | 🟡 Média | Complexo mas possível (Breez prova) |
| **Nostr Discovery** | ✅ Alta | go-nostr é simples e funcional |
| **LNURL-Pay** | ✅ Alta | Protocolo HTTP simples |
| **Cross-Compilation** | ✅ Alta | Go é excelente nisso |
| **UX "Pix-like"** | 🟡 Média | Limitado pelo requisito de estar online |

---

## Próximos Passos Recomendados

1. **POC Wallet Core** — Validar geração de seed + Neutrino sync em Regtest
2. **POC LND Embutido** — Testar se LND funciona como library ou precisa ser subprocess
3. **Definir LSP** — Escolher provedor de liquidez para onboarding
4. **Regtest Environment** — Script automatizado para levantar rede de teste
