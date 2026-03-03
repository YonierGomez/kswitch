#!/usr/bin/env bash
# release.sh — automatiza el flujo completo de release para ksw
# Uso: ./scripts/release.sh <version> "<descripción>"
# Ejemplo: ./scripts/release.sh 1.3.2 "fix: corregir bug en alias"

set -e

VERSION="$1"
DESCRIPTION="$2"
HOMEBREW_TAP_DIR="/opt/homebrew/Library/Taps/yoniergomez/homebrew-ksw"

if [[ -z "$VERSION" || -z "$DESCRIPTION" ]]; then
  echo "Uso: $0 <version> \"<descripción>\""
  echo "Ejemplo: $0 1.3.2 \"fix: corregir bug en alias\""
  exit 1
fi

TAG="v$VERSION"

echo "🚀 Iniciando release $TAG..."

# 1. Asegurar que estamos en main actualizado
git checkout main
git pull origin main

# 2. Bump version constant en main.go y hacer PR antes del tag
echo "📝 Actualizando version constant en main.go..."
VERSION_BRANCH="chore/bump-version-$TAG"
git checkout -b "$VERSION_BRANCH"
sed -i '' "s|const version = \"[0-9]*\.[0-9]*\.[0-9]*\"|const version = \"$VERSION\"|g" main.go
git add main.go
git commit -m "chore: bump version to $TAG"
git push origin "$VERSION_BRANCH"

gh pr create \
  --repo YonierGomez/ksw \
  --base main \
  --head "$VERSION_BRANCH" \
  --title "chore: bump version to $TAG" \
  --body "$DESCRIPTION"

gh pr merge \
  --repo YonierGomez/ksw \
  --squash \
  "$VERSION_BRANCH"

git checkout main
git pull origin main

# 3. Crear y pushear el tag
echo "🏷️  Creando tag $TAG..."
git tag "$TAG"
git push origin "$TAG"

# 3. Esperar a que el tarball esté disponible en GitHub
echo "⏳ Esperando tarball en GitHub..."
TARBALL_URL="https://github.com/YonierGomez/ksw/archive/refs/tags/${TAG}.tar.gz"
for i in {1..20}; do
  STATUS=$(curl -sI "$TARBALL_URL" | head -1 | awk '{print $2}')
  if [[ "$STATUS" == "302" || "$STATUS" == "200" ]]; then
    echo "✔ Tarball disponible"
    break
  fi
  echo "  Intento $i/20..."
  sleep 5
done

# 4. Calcular sha256
echo "🔐 Calculando sha256..."
SHA256=$(curl -sL "$TARBALL_URL" | shasum -a 256 | awk '{print $1}')
echo "   sha256: $SHA256"

# 5. Actualizar Formula/ksw.rb en el repo principal
echo "📝 Actualizando Formula/ksw.rb..."
BRANCH="chore/bump-$TAG"
git checkout -b "$BRANCH"
sed -i '' "s|refs/tags/v[0-9]*\.[0-9]*\.[0-9]*.tar.gz|refs/tags/${TAG}.tar.gz|g" Formula/ksw.rb
sed -i '' "s|sha256 \"[a-f0-9]*\"|sha256 \"$SHA256\"|g" Formula/ksw.rb
git add Formula/ksw.rb
git commit -m "chore: bump formula to $TAG"
git push origin "$BRANCH"

# 6. Crear y mergear PR en repo principal
echo "🔀 Creando PR en YonierGomez/ksw..."
gh pr create \
  --repo YonierGomez/ksw \
  --base main \
  --head "$BRANCH" \
  --title "chore: bump formula to $TAG" \
  --body "$DESCRIPTION"

gh pr merge \
  --repo YonierGomez/ksw \
  --squash \
  "$BRANCH"

# 7. Actualizar formula en homebrew-ksw
echo "🍺 Actualizando homebrew-ksw..."
git -C "$HOMEBREW_TAP_DIR" checkout main
git -C "$HOMEBREW_TAP_DIR" pull origin main
git -C "$HOMEBREW_TAP_DIR" checkout -b "$BRANCH"
sed -i '' "s|refs/tags/v[0-9]*\.[0-9]*\.[0-9]*.tar.gz|refs/tags/${TAG}.tar.gz|g" "$HOMEBREW_TAP_DIR/Formula/ksw.rb"
sed -i '' "s|sha256 \"[a-f0-9]*\"|sha256 \"$SHA256\"|g" "$HOMEBREW_TAP_DIR/Formula/ksw.rb"
git -C "$HOMEBREW_TAP_DIR" add Formula/ksw.rb
git -C "$HOMEBREW_TAP_DIR" commit -m "chore: bump formula to $TAG"
git -C "$HOMEBREW_TAP_DIR" push origin "$BRANCH"

# 8. Crear y mergear PR en homebrew-ksw
echo "🔀 Creando PR en YonierGomez/homebrew-ksw..."
gh pr create \
  --repo YonierGomez/homebrew-ksw \
  --base main \
  --head "$BRANCH" \
  --title "chore: bump formula to $TAG" \
  --body "$DESCRIPTION"

gh pr merge \
  --repo YonierGomez/homebrew-ksw \
  --squash \
  "$BRANCH"

# 9. Actualizar brew local
echo "🔄 Actualizando brew..."
brew update && brew upgrade ksw

echo ""
echo "✅ Release $TAG completado"
echo "   - Tag pusheado"
echo "   - Formula actualizada en ksw y homebrew-ksw"
echo "   - brew upgrade ksw listo"
