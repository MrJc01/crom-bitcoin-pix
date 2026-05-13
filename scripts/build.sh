#!/usr/bin/env bash
# Cross-compilation script para Crom-Pay
set -euo pipefail

BINARY="crom-pay"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS="-s -w -X main.Version=${VERSION}"
OUTPUT_DIR="dist"

echo "📦 Crom-Pay Cross-Compilation v${VERSION}"
echo "==========================================="

mkdir -p "${OUTPUT_DIR}"

platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    os="${platform%/*}"
    arch="${platform#*/}"
    output="${OUTPUT_DIR}/${BINARY}-${os}-${arch}"

    if [ "${os}" = "windows" ]; then
        output="${output}.exe"
    fi

    echo "  🔨 ${os}/${arch} → ${output}"
    CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" \
        go build -ldflags="${LDFLAGS}" -o "${output}" ./cmd/crom-pay
done

echo ""
echo "✅ Build completo!"
echo ""
ls -lh "${OUTPUT_DIR}/"
