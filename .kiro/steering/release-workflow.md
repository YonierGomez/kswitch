# ksw Maintainer & Release Workflow

## Proyecto
CLI para Kubernetes context switching. Repo: `git@github.com:YonierGomez/ksw.git`

## Reglas siempre activas
- Siempre trabajar desde `main`, crear rama para cada cambio, nunca pushear directo a `main`
- Después de cada merge: volver a `main`, borrar rama local y quedar limpio
- Git remote usa SSH (`git@github.com:...`), nunca HTTPS
- `gh` CLI instalado y autenticado como `YonierGomez` via SSH
- No incluir nombres reales de clusters (nequi, sufi, clientes, apic, etc.) en landing ni demos
- No incluir `rm -f ~/.ksw.json` en tapes de VHS

## Tipos de rama
`feat/`, `fix/`, `chore/`, `docs/`

## Flujo cambios normales
```bash
git checkout main && git pull origin main
git checkout -b <tipo>/<descripcion>
# hacer cambios
git add . && git commit -m "<tipo>: <descripcion>"
git push origin <tipo>/<descripcion>
gh pr create --repo YonierGomez/ksw --base main --head <rama> --title "<titulo>" --body "<descripcion>"
gh pr merge --repo YonierGomez/ksw --squash <rama>
git checkout main && git pull origin main && git branch -D <rama>
```

## Flujo release (con versión)
Usar el script automatizado:
```bash
./scripts/release.sh <version> "<descripcion>"
# Ejemplo: ./scripts/release.sh 1.3.4 "fix: corregir bug en alias"
```

El script valida antes de hacer cualquier cosa:
1. `gh` autenticado
2. `main` limpio y actualizado
3. Tag no existe local ni remoto
4. `const version` en `main.go` coincide con la versión
5. Archivos clave presentes (`main.go`, `go.mod`, `Formula/ksw.rb`)
6. Homebrew tap dir existe en `/opt/homebrew/Library/Taps/yoniergomez/homebrew-ksw`
7. Tarball disponible en GitHub
8. Tarball contiene el `const version` correcto
9. Espera merge de PRs antes de continuar
10. Verifica `ksw -v` al final

El script también actualiza `index.html` automáticamente.

## Consistencia de versión — siempre verificar que coincidan
- `const version` en `main.go`
- `Formula/ksw.rb` — url tag y sha256
- `/opt/homebrew/Library/Taps/yoniergomez/homebrew-ksw/Formula/ksw.rb`
- `softwareVersion` en `index.html`
- Badge `⎈ vX.Y.Z · AI-Powered` en `index.html`
- Footer en `index.html`

## Repos involucrados
- Código fuente: `git@github.com:YonierGomez/ksw.git`
- Homebrew tap: `git@github.com:YonierGomez/homebrew-ksw.git` (local: `/opt/homebrew/Library/Taps/yoniergomez/homebrew-ksw`)
- Landing page: `https://yoniergomez.github.io/ksw/` (se despliega automáticamente desde `main`)

## Flujo PR en homebrew-ksw (manual si es necesario)
```bash
git -C /opt/homebrew/Library/Taps/yoniergomez/homebrew-ksw checkout main
git -C /opt/homebrew/Library/Taps/yoniergomez/homebrew-ksw pull origin main
git -C /opt/homebrew/Library/Taps/yoniergomez/homebrew-ksw checkout -b <rama>
# editar Formula/ksw.rb
git -C /opt/homebrew/Library/Taps/yoniergomez/homebrew-ksw add Formula/ksw.rb
git -C /opt/homebrew/Library/Taps/yoniergomez/homebrew-ksw commit -m "chore: bump formula to vX.Y.Z"
git -C /opt/homebrew/Library/Taps/yoniergomez/homebrew-ksw push origin <rama>
gh pr create --repo YonierGomez/homebrew-ksw --base main --head <rama> --title "..." --body "..."
gh pr merge --repo YonierGomez/homebrew-ksw --squash <rama>
```

## Si el tag apunta al commit incorrecto
```bash
git tag -d vX.Y.Z && git push origin :refs/tags/vX.Y.Z
git tag vX.Y.Z && git push origin vX.Y.Z
git show vX.Y.Z:main.go | grep "const version"  # verificar
```

## brew reinstall vs upgrade
Si `brew upgrade` dice "already installed" pero `ksw -v` muestra versión incorrecta:
```bash
brew reinstall ksw
```
