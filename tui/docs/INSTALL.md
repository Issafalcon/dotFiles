# Installing & Using DotFiles TUI

Download the AppImage from the [latest release](https://github.com/issafalcon/dotFiles/releases/latest), make it executable, and run it. That's it — no Go toolchain, no dependencies.

## Quick Start (Fresh Linux Install)

```bash
# Download the latest AppImage
curl -fsSL -o dotfiles-tui \
  "https://github.com/issafalcon/dotFiles/releases/latest/download/dotfiles-tui-$(curl -s https://api.github.com/repos/issafalcon/dotFiles/releases/latest | grep tag_name | cut -d'"' -f4)-x86_64.AppImage"

# Make it executable and run
chmod +x dotfiles-tui
./dotfiles-tui
```

Or pin a specific version:

```bash
VERSION=v1.0.0
curl -fsSL -o dotfiles-tui \
  "https://github.com/issafalcon/dotFiles/releases/download/${VERSION}/dotfiles-tui-${VERSION}-x86_64.AppImage"
chmod +x dotfiles-tui
./dotfiles-tui
```

## Alternative: Tarball

If you prefer a plain binary (e.g., inside WSL where FUSE may not be available):

```bash
VERSION=v1.0.0
curl -fsSL "https://github.com/issafalcon/dotFiles/releases/download/${VERSION}/dotfiles-tui-${VERSION}-linux-amd64.tar.gz" \
  | tar -xz
./dotfiles-tui
```

## Alternative: Build from Source

Requires Go 1.25+ and `make`:

```bash
git clone https://github.com/issafalcon/dotFiles.git
cd dotFiles/tui
make build
./build/dotfiles-tui
```

Or install to `/usr/local/bin`:

```bash
make install
dotfiles-tui
```

## System-wide Install

To install the AppImage so it's on your `PATH`:

```bash
sudo mv dotfiles-tui /usr/local/bin/dotfiles-tui
```

## Verify Download Integrity

Each release includes a `checksums-<version>.sha256` file:

```bash
VERSION=v1.0.0
curl -fsSL -O "https://github.com/issafalcon/dotFiles/releases/download/${VERSION}/checksums-${VERSION}.sha256"
sha256sum -c checksums-${VERSION}.sha256
```

## Prerequisites

The app checks prerequisites on startup, but for reference:

| Tool     | Purpose                          | Install                          |
|----------|----------------------------------|----------------------------------|
| `git`    | Clone/manage the dotfiles repo   | `sudo apt install git`           |
| `stow`   | Symlink config files             | `sudo apt install stow`          |
| `curl`   | Download packages                | `sudo apt install curl`          |
| `wget`   | Download packages                | `sudo apt install wget`          |

The TUI will warn you about any missing prerequisites before proceeding.

## FUSE Note (AppImage)

AppImages require FUSE to mount. Most modern distros include it. If you get a FUSE error:

```bash
# Ubuntu/Debian
sudo apt install libfuse2

# Or extract and run without FUSE
./dotfiles-tui --appimage-extract
./squashfs-root/AppRun
```

## Keyboard Controls

| Key          | Action                        |
|--------------|-------------------------------|
| `j` / `k`   | Navigate sidebar up/down      |
| `Enter`      | Select module                 |
| `i`          | Install selected module       |
| `d`          | Uninstall selected module     |
| `s` / `/`    | Search modules                |
| `Tab`        | Switch detail tab             |
| `Shift+Tab`  | Switch sidebar ↔ detail focus |
| `o`          | Open module website           |
| `?`          | Show help                     |
| `q`          | Quit                          |

## Creating a Release

Push a version tag to trigger the CI/CD pipeline:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The GitHub Actions workflow will automatically build, test, create the AppImage, and publish a GitHub Release.
