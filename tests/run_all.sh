#!/usr/bin/env bash
# ═══════════════════════════════════════════════════════════════════════════════
# Crom-Pay Test Runner — Executa todos os testes e gera relatório
# ═══════════════════════════════════════════════════════════════════════════════
set -uo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REPORT_DIR="$PROJECT_ROOT/tests/reports"
TIMESTAMP="$(date +%Y%m%d_%H%M%S)"
REPORT_FILE="$REPORT_DIR/report_${TIMESTAMP}.md"
LOG_FILE="$REPORT_DIR/raw_${TIMESTAMP}.log"
TMP_DIR="$(mktemp -d)"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; BOLD='\033[1m'; NC='\033[0m'

TOTAL_PASS=0
TOTAL_FAIL=0

mkdir -p "$REPORT_DIR"
cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

header() {
    echo -e "\n${CYAN}${BOLD}═══════════════════════════════════════════════${NC}"
    echo -e "${CYAN}${BOLD}  $1${NC}"
    echo -e "${CYAN}${BOLD}═══════════════════════════════════════════════${NC}\n"
}

run_go_tests() {
    local label="$1"
    local slug="$2"
    local tags="$3"
    local pkg="$4"
    local timeout="${5:-60s}"
    local outf="$TMP_DIR/${slug}.out"

    echo -e "${YELLOW}  ▶ ${label}${NC}"
    cd "$PROJECT_ROOT"
    go test -v -count=1 -tags "$tags" -timeout "$timeout" "$pkg" > "$outf" 2>&1 || true

    local p=0 f=0
    p="$(grep -c '^--- PASS' "$outf" || true)"
    f="$(grep -c '^--- FAIL' "$outf" || true)"

    TOTAL_PASS=$((TOTAL_PASS + p))
    TOTAL_FAIL=$((TOTAL_FAIL + f))

    if [ "$f" -eq 0 ] && [ "$p" -gt 0 ]; then
        echo -e "${GREEN}    ✅ ${p} testes passaram${NC}"
    elif [ "$f" -gt 0 ]; then
        echo -e "${RED}    ❌ ${f} falha(s) em $((p + f)) testes${NC}"
    else
        echo -e "${YELLOW}    ⚠️  Nenhum teste encontrado${NC}"
    fi

    echo "=== ${label} ===" >> "$LOG_FILE"
    cat "$outf" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
}

run_shell_test() {
    local label="$1"
    local script="$2"
    local outf="$TMP_DIR/e2e.out"

    echo -e "${YELLOW}  ▶ ${label}${NC}"
    bash "$script" > "$outf" 2>&1 || true

    local p=0 f=0
    p="$(grep -c '✅ PASS' "$outf" || true)"
    f="$(grep -c '❌ FAIL' "$outf" || true)"

    TOTAL_PASS=$((TOTAL_PASS + p))
    TOTAL_FAIL=$((TOTAL_FAIL + f))

    if [ "$f" -eq 0 ] && [ "$p" -gt 0 ]; then
        echo -e "${GREEN}    ✅ ${p} testes passaram${NC}"
    elif [ "$f" -gt 0 ]; then
        echo -e "${RED}    ❌ ${f} falha(s)${NC}"
    fi

    echo "=== ${label} ===" >> "$LOG_FILE"
    cat "$outf" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
}

# ═══════════════════════════════════════════════════════════════════════════════

echo ""
echo -e "${BOLD}🧪 Crom-Pay Test Suite — $(date)${NC}"
echo ""

echo "# Crom-Pay Test Log — $TIMESTAMP" > "$LOG_FILE"
echo "# Go: $(go version)" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# ─── Build ────────────────────────────────────────────────────────────────────
header "📦 Fase 0: Build"
echo -e "${YELLOW}  ▶ Compilando binário...${NC}"
cd "$PROJECT_ROOT"
if go build -o bin/crom-pay ./cmd/crom-pay 2>&1; then
    BIN_SIZE="$(ls -lh bin/crom-pay | awk '{print $5}')"
    echo -e "${GREEN}    ✅ Build OK (${BIN_SIZE})${NC}"
    BUILD_STATUS="OK"
else
    echo -e "${RED}    ❌ Build FALHOU${NC}"
    BUILD_STATUS="FALHOU"
fi

# ─── Testes ───────────────────────────────────────────────────────────────────
header "🔬 Fase 1: Testes Unitários (Wallet)"
run_go_tests "Wallet Core" "unit" "" "./internal/wallet/..." "60s"

header "🗄️  Fase 1b: Testes Unitários (Storage)"
run_go_tests "Storage CRUD" "storage" "" "./internal/storage/..." "60s"

header "🔗 Fase 1c: Testes Unitários (Chain)"
run_go_tests "Chain Client + TX" "chain" "" "./internal/chain/..." "60s"

header "⚡ Fase 1d: Testes Unitários (Lightning)"
run_go_tests "Lightning Client" "lightning" "" "./internal/lightning/..." "60s"

header "🌐 Fase 1e: Testes Unitários (Nostr)"
run_go_tests "Nostr Keys + Events" "nostr" "" "./internal/nostr/..." "60s"

header "📱 Fase 1f: Testes Unitários (UI)"
run_go_tests "QR Code + TUI" "ui" "" "./internal/ui/..." "60s"

header "📇 Fase 1g: Testes Unitários (Contacts)"
run_go_tests "Contacts CRUD" "contacts" "" "./internal/contacts/..." "60s"

header "🔄 Fase 2: Testes de Integração"
run_go_tests "Wallet Flow" "integration" "integration" "./tests/integration/..." "120s"

header "🔴 Fase 3: Pentest / Segurança"
run_go_tests "Entropy + Crypto + Keys" "security" "security" "./tests/security/..." "180s"

header "💪 Fase 4: Testes de Stress"
run_go_tests "Concorrência + Limites" "stress" "stress" "./tests/stress/..." "300s"

header "🌐 Fase 5: End-to-End (CLI)"
run_shell_test "CLI Binary Test" "$PROJECT_ROOT/tests/e2e/cli_test.sh"

# ═══════════════════════════════════════════════════════════════════════════════
# RELATÓRIO
# ═══════════════════════════════════════════════════════════════════════════════

TOTAL_TESTS=$((TOTAL_PASS + TOTAL_FAIL))
if [ "$TOTAL_TESTS" -gt 0 ]; then
    SUCCESS_RATE="$((TOTAL_PASS * 100 / TOTAL_TESTS))%"
else
    SUCCESS_RATE="N/A"
fi

section() {
    local f="$TMP_DIR/$1.out"
    if [ -f "$f" ]; then
        grep -E '(=== RUN|--- PASS|--- FAIL|^ok|✅|🔴|📊|⚠️)' "$f" 2>/dev/null || echo "Sem output"
    else
        echo "Sem output"
    fi
}

cat > "$REPORT_FILE" << ENDREPORT
# 📊 Relatório de Testes — Crom-Pay

| | |
|---|---|
| **Data** | $(date "+%Y-%m-%d %H:%M:%S") |
| **Go** | $(go version | awk '{print $3}') |
| **OS** | $(uname -s) $(uname -m) |
| **Build** | $BUILD_STATUS |
| **Binário** | ${BIN_SIZE:-N/A} |

---

## Resumo

| Métrica | Valor |
|---|---|
| **Total** | $TOTAL_TESTS |
| **Pass** | ✅ $TOTAL_PASS |
| **Fail** | ❌ $TOTAL_FAIL |
| **Taxa** | $SUCCESS_RATE |

---

## Fase 1: Unitários
\`\`\`
$(section unit)
\`\`\`

## Fase 2: Integração
\`\`\`
$(section integration)
\`\`\`

## Fase 3: Segurança
\`\`\`
$(section security)
\`\`\`

## Fase 4: Stress
\`\`\`
$(section stress)
\`\`\`

## Fase 5: E2E
\`\`\`
$(cat "$TMP_DIR/e2e.out" 2>/dev/null | grep -E '(✅|❌|🧪|═)' || echo "Sem output")
\`\`\`

---
Log bruto: \`$(basename "$LOG_FILE")\`
ENDREPORT

# ═══════════════════════════════════════════════════════════════════════════════

echo ""
header "📊 Resultado Final"
echo -e "  Total: ${BOLD}${TOTAL_TESTS}${NC} testes"
echo -e "  Pass:  ${GREEN}${TOTAL_PASS}${NC}"
echo -e "  Fail:  ${RED}${TOTAL_FAIL}${NC}"
echo ""
echo -e "  📄 Relatório: ${CYAN}${REPORT_FILE}${NC}"
echo -e "  📋 Log bruto: ${CYAN}${LOG_FILE}${NC}"
echo ""

if [ "$TOTAL_FAIL" -eq 0 ]; then
    echo -e "${GREEN}${BOLD}  ✅ TODOS OS TESTES PASSARAM!${NC}"
else
    echo -e "${RED}${BOLD}  ❌ ${TOTAL_FAIL} TESTE(S) FALHARAM${NC}"
fi
echo ""
exit "$TOTAL_FAIL"
