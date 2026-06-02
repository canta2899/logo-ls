![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Version](https://img.shields.io/github/v/release/canta2899/logo-ls?style=for-the-badge)
![Tests](https://img.shields.io/github/actions/workflow/status/canta2899/logo-ls/test.yml?label=Tests&style=for-the-badge)

<div align="center">
<img src="./.github/assets/screen.png" alt="logo-ls screenshot" width="600"/>
</div>

<h1 align="center">
    logo-ls
</h1>

A fork of [logo-ls](https://github.com/Yash-Handa/logo-ls) which I ended up maintaining since the original repository went unmaintained some years ago. Feel free to open an issue or a pull request if you have any questions or want to contribute. If you want to add icons, check the [Adding Icons](#adding-icons) section below for instructions on how to do so.

---

## Installation

### Prerequisites

- Ensure your terminal is using a Nerd Font to see the icons properly. You can download your preferred Nerd Font from [here](https://www.nerdfonts.com/font-downloads). Some terminal emulators such as [Ghostty](https://ghostty.org) come with built in support for Nerd Fonts, so you don't have to worry about it.
- The command will be installed as `logo-ls`, so you can optionally set an alias for `ls` to `logo-ls` if you want to use it as a drop in replacement for `ls`.

### Arch Linux

Install the logo-ls [AUR package](https://aur.archlinux.org/packages/logo-ls).

```bash
yay -S logo-ls
```

### Homebrew (tap)

```bash
brew install canta2899/homebrew-tap/logo-ls
```

### Binary Release (Linux/OSX/Windows)

Optionally, you can set the variables `LOGO_LS_INSTALL_DIR` and/or `LOGO_LS_VERSION` to specify a custom installation directory and/or version to install. By default, the scripts will install the latest version of logo-ls to `~/.local/bin`.

#### Linux/OSX

```bash
curl -L https://raw.githubusercontent.com/canta2899/logo-ls/refs/heads/main/get.sh | sh
```

#### Windows

```powershell
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
irm https://raw.githubusercontent.com/canta2899/logo-ls/main/get.ps1 | iex
```

#### Manual Install

You can download the binary for your platform from the [releases page](https://github.com/canta2899/logo-ls/releases/).

### Build from source

Clone the repository

```bash
git clone https://github.com/canta2899/logo-ls
```

Build the binary, which is outputted to the root directory of the repository:

```bash
make logo-ls
```
---

## Icons

`logo-ls` comes with a huge built in set of icons for common file types and directories. You can either contribute to the built in set by opening a pull request with an update to the skill file, or you can add custom icons or overrides for your own use.

### Contribute with new Icons

If you use any coding agent (OpenCode, Gemini CLI, Claude Code, etc.) there's a built in skill called `/add-icon` which you can use to let your agent do the job for you. If you want to do it manually, you can pretend to be a coding agent and read the skill file yourself.

> **Note for Windows contributors:** the skill file lives at `add-icon-skill.md` in the repo root, and the paths under `.agents/skills/add-icon/SKILL.md` and `.claude/skills/add-icon/SKILL.md` are symlinks to it. Git for Windows does not create real symlinks by default, so these may be checked out as plain text files containing the link target. To get working symlinks, enable Developer Mode (or run as admin) and set `git config --global core.symlinks true` before cloning.

### Overrides

If you want to add or override icons, you can drop a YAML file in either of the following locations:

- `$XDG_CONFIG_HOME/logo-ls/logo-ls-overrides.yaml` (defaults to `~/.config/logo-ls/logo-ls-overrides.yaml`)
- `~/.logo-ls-overrides.yaml`

The first existing file wins. The file is read **once at startup** and merged with the built-in icon set, so it does not impact per-directory listing speed. If the file is missing or empty the loader just skips it, if it cannot be parsed `logo-ls` prints a single warning and continues with built-in icons only.

#### Schema

Each entry can set a `glyph`, a `color`, or both. Whichever field you leave out keeps the built-in value for that match, so you can just recolor an existing icon, swap only its glyph, or define a brand-new icon by setting both.

```yaml
extensions:        # matched by file extension (without the leading dot)
  rs:
    glyph: "U+E7A8"    # full override: new glyph and new color
    color: "#dea584"
  go:
    color: "#ff5555"   # color-only: keep the built-in Go glyph, just recolor it
  py:
    glyph: "0xE606"    # glyph-only: keep the built-in Python color

files:             # matched by full file name (case-insensitive)
  .envrc:
    glyph: "U+E60B"
    color: "#ecd53f"

directories:       # matched by full directory name (case-insensitive)
  myproject:
    color: "#42a5f5" # recolor a directory without changing its glyph

sub_extensions:    # matched by "<last-segment><ext>"
  d.ts:
    glyph: "U+E628"
    color: "#3178c6"
```

Glyphs can be written as `U+XXXX` or `0xXXXX` (parsed as hex codepoints) or as literal strings. Colors are `#RRGGBB` (or the shorthand `#RGB`). User entries take priority over the built-in icons.

#### Flags

Two flags control the override loader:

- `--no-override` skips override loading entirely and forces built-in icons only.
- `--override-file <path>` loads overrides from an explicit YAML path instead of the default discovery locations.

If neither flag is passed the default discovery paths above are used.

