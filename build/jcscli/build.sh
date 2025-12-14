#!/usr/bin/env bash
set -euo pipefail

# Configuration
APP_NAME="jcscli"
VERSION=${1:-"dev"}
DIST_DIR="dist"

# Clean dist directory
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

echo "Building $APP_NAME version $VERSION..."

# Target platforms
PLATFORMS=(
	"linux/amd64"
	"linux/arm64"
	"darwin/amd64"
	"darwin/arm64"
	"windows/amd64"
)

# Build binaries
for PLATFORM in "${PLATFORMS[@]}"; do
	GOOS=${PLATFORM%/*}
	GOARCH=${PLATFORM#*/}
	OUTPUT="$DIST_DIR/${APP_NAME}-${GOOS}-${GOARCH}"

	if [ "$GOOS" = "windows" ]; then
		OUTPUT="${OUTPUT}.exe"
	fi

	echo "→ Building $GOOS/$GOARCH"
	GOOS=$GOOS GOARCH=$GOARCH go build \
		-ldflags "-X main.version=$VERSION" \
		-o "$OUTPUT" ./cmd/jcscli

	# Package
	if [ "$GOOS" = "windows" ]; then
		zip -j "${OUTPUT%.exe}.zip" "$OUTPUT"
		rm "$OUTPUT"
	else
		tar -czf "${OUTPUT}.tar.gz" -C "$DIST_DIR" "$(basename "$OUTPUT")"
		rm "$OUTPUT"
	fi
done

# Update CHANGELOG.md
DATE=$(date +"%Y-%m-%d")
echo "Updating CHANGELOG.md..."
cat <<EOF >>CHANGELOG.md

## [$VERSION] - $DATE
### Added
- Release artifacts built for Linux, macOS, Windows.
- Canonical JSON CLI tool updates.

EOF

echo "All artifacts are in $DIST_DIR/"
