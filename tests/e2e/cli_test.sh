#!/usr/bin/env bash
# ─── E2E Test: Testa o binário compilado real ──────────────────────────────────
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
BINARY="$PROJECT_ROOT/bin/crom-pay"
TEST_DIR=$(mktemp -d)
PASS="e2e-test-senha-forte"
ERRORS=0

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_pass() { echo -e "${GREEN}  ✅ PASS: $1${NC}"; }
log_fail() { echo -e "${RED}  ❌ FAIL: $1${NC}"; ERRORS=$((ERRORS + 1)); }
log_test() { echo -e "${YELLOW}🧪 $1${NC}"; }

cleanup() { rm -rf "$TEST_DIR"; }
trap cleanup EXIT

echo ""
echo "═══════════════════════════════════════════════"
echo "  🧪 E2E Tests — Crom-Pay Binary"
echo "═══════════════════════════════════════════════"
echo ""

# ─── Test 1: Help flag ────────────────────────────────────────────────────
log_test "Test 1: --help retorna sucesso"
if "$BINARY" --help > /dev/null 2>&1; then
    log_pass "--help funciona"
else
    log_fail "--help falhou"
fi

# ─── Test 2: Version flag ────────────────────────────────────────────────
log_test "Test 2: --version retorna versão"
VERSION_OUT=$("$BINARY" --version 2>&1 || true)
if echo "$VERSION_OUT" | grep -qi "crom-pay"; then
    log_pass "--version mostra nome do binário"
else
    log_fail "--version não contém 'crom-pay': $VERSION_OUT"
fi

# ─── Test 3: Wallet create ───────────────────────────────────────────────
log_test "Test 3: wallet create gera carteira"
CREATE_OUT=$(printf "%s\n%s\n" "$PASS" "$PASS" | "$BINARY" wallet create --data-dir "$TEST_DIR/wallet1" --network testnet 2>&1)
if echo "$CREATE_OUT" | grep -q "tb1q"; then
    log_pass "wallet create gerou endereço tb1q"
else
    log_fail "wallet create não gerou endereço SegWit"
    echo "  Output: $CREATE_OUT"
fi

# Extrair endereço para comparação
ADDR_CREATE=$(echo "$CREATE_OUT" | grep -oP 'tb1q[a-z0-9]+' | head -1)

# ─── Test 4: Wallet balance ─────────────────────────────────────────────
log_test "Test 4: wallet balance retorna saldo"
BAL_OUT=$(printf "%s\n" "$PASS" | "$BINARY" wallet balance --data-dir "$TEST_DIR/wallet1" 2>&1)
if echo "$BAL_OUT" | grep -q "0 sats"; then
    log_pass "wallet balance mostra 0 sats"
else
    log_fail "wallet balance não mostra saldo"
fi

# ─── Test 5: Wallet address ─────────────────────────────────────────────
log_test "Test 5: wallet address retorna mesmo endereço"
ADDR_OUT=$(printf "%s\n" "$PASS" | "$BINARY" wallet address --data-dir "$TEST_DIR/wallet1" 2>&1)
ADDR_SHOW=$(echo "$ADDR_OUT" | grep -oP 'tb1q[a-z0-9]+' | head -1)

if [ "$ADDR_CREATE" = "$ADDR_SHOW" ]; then
    log_pass "Endereço consistente: $ADDR_CREATE"
else
    log_fail "Endereço inconsistente: create=$ADDR_CREATE address=$ADDR_SHOW"
fi

# ─── Test 6: Senha errada é rejeitada ────────────────────────────────────
log_test "Test 6: Senha errada é rejeitada"
WRONG_OUT=$(printf "senha-errada\n" | "$BINARY" wallet balance --data-dir "$TEST_DIR/wallet1" 2>&1 || true)
if echo "$WRONG_OUT" | grep -qi "incorreta\|corrompid\|erro\|error\|wrong"; then
    log_pass "Senha errada rejeitada"
else
    log_fail "Senha errada pode ter sido aceita!"
fi

# ─── Test 7: Criar duplicado falha ───────────────────────────────────────
log_test "Test 7: Criar carteira duplicada falha"
DUP_OUT=$(printf "%s\n%s\n" "$PASS" "$PASS" | "$BINARY" wallet create --data-dir "$TEST_DIR/wallet1" --network testnet 2>&1 || true)
if echo "$DUP_OUT" | grep -qi "existe\|already\|duplica"; then
    log_pass "Criação duplicada rejeitada"
else
    log_fail "Criação duplicada pode ter sido aceita"
fi

# ─── Test 8: Wallet subcommand help ──────────────────────────────────────
log_test "Test 8: wallet --help mostra subcomandos"
WHELP=$("$BINARY" wallet --help 2>&1)
for cmd in create balance restore address; do
    if echo "$WHELP" | grep -q "$cmd"; then
        log_pass "Subcomando '$cmd' listado"
    else
        log_fail "Subcomando '$cmd' não encontrado"
    fi
done

# ─── Resultado ───────────────────────────────────────────────────────────
echo ""
echo "═══════════════════════════════════════════════"
if [ $ERRORS -eq 0 ]; then
    echo -e "  ${GREEN}✅ TODOS OS TESTES E2E PASSARAM${NC}"
else
    echo -e "  ${RED}❌ $ERRORS TESTE(S) FALHARAM${NC}"
fi
echo "═══════════════════════════════════════════════"
echo ""

exit $ERRORS
