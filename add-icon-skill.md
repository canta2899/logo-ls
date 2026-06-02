---
name: add-icon
description: Add a new icon to logo-ls by updating the icon maps
---

# Add Icon to logo-ls

You are helping the user add a new icon to the logo-ls project. This project uses Nerd Fonts glyphs to display icons next to file and directory names in terminal listings.

There are two ways to add an icon, and the user should pick one **before** you start editing files:

1. **Personal YAML override** (no contribution) — drop entries into the user's `~/.config/logo-ls/logo-ls-overrides.yaml` (or `~/.logo-ls-overrides.yaml`). Loaded once at startup; no rebuild required. Use this when the user just wants the icon on their own machine, or wants to try a glyph/color combination before committing to a PR.
2. **Built-in contribution** — edit the Go maps in `internal/icons/` and rebuild. Use this when the user wants the icon shipped upstream.

Workflow tip: a user can iterate on a glyph/color in YAML, then ask you to **promote** the working entries to the built-in maps for a contribution. If they hand you a YAML file (or paste its contents) and say "make this a contribution", read it, then run the Step 3 procedure below for each entry.

## Architecture

The icon system lives in the `internal/icons/` package.

**Built-in (contribution) maps:**

- **`internal/icons/icons_map.go`** — `IconSet` map: defines all available icons. Each entry has a name (key), a glyph (Unicode codepoint), and an RGB color. This is where every new icon must be registered first.
- **`internal/icons/icons_ext.go`** — `IconExt` map: maps file extensions (e.g. `"js"`, `"py"`) to entries in `IconSet`.
- **`internal/icons/icons_files.go`** — `IconFileName` map: maps specific file names (e.g. `"tsconfig.json"`, `".babelrc"`) to entries in `IconSet`.
- **`internal/icons/icons_sub_ext.go`** — `IconSubExt` map: maps compound sub-extensions (e.g. `"spec.ts"`, `"d.ts"`, `"gitlab-ci.yml"`) to entries in `IconSet`. Used to override the default extension icon for specific patterns.

**User override loader** (do not edit unless the user asks):

- **`internal/icons/override.go`** — parses the user's YAML file; the four top-level keys (`extensions`, `files`, `directories`, `sub_extensions`) mirror the four built-in maps above. User entries take priority over built-ins at lookup time. Each entry is sparse: a user can set `glyph` only, `color` only, or both — unset fields keep the built-in value for that match.

The `IconInfo` struct looks like this:

```go
type IconInfo struct {
    Glyph        string
    Color        [3]uint8 // RGB color
    IsExecutable bool
}
```

## Procedure

First, decide which path to take:

- If the user says "let me try it first", "I want it just for me", "don't open a PR", or similar → follow **Procedure A (YAML override)**.
- If the user says "contribute", "upstream", "PR", or is editing the repo source tree → follow **Procedure B (built-in)**.
- If they hand you a YAML override and ask to promote it, read the file, infer the same fields you would otherwise ask about, and skip to Procedure B Step 3.

When in doubt, ask which one they want before editing anything.

---

## Procedure A — User YAML override

### Step 1 — Locate or create the YAML file

Check, in order:

1. `$XDG_CONFIG_HOME/logo-ls/logo-ls-overrides.yaml` (defaults to `~/.config/logo-ls/logo-ls-overrides.yaml`)
2. `~/.logo-ls-overrides.yaml`

The first existing file wins at runtime. If neither exists, create the XDG path. Do not move an existing file — append to whichever one is already in use.

### Step 2 — Gather icon details

Ask the user:

1. **Match type**: `extensions`, `files`, `directories`, or `sub_extensions` (same semantics as the built-in maps).
2. **Pattern**: the extension (without leading dot), file name, directory name, or `last-segment.ext` sub-extension. Case-insensitive at lookup time.
3. **Glyph** *(optional)*: a Nerd Font codepoint (`U+E7A8`, `0xe7a8`) or a literal string. Codepoint forms with `U+` / `0x` prefix are parsed as hex; everything else is treated as a literal.
4. **Color** *(optional)*: `#RRGGBB` or `#RGB`.

At least one of `glyph` or `color` must be set. If the user only wants to recolor an existing icon, ask only for the color and skip the glyph; if they only want to change the glyph, do the inverse. The unset field will fall through to the built-in for the same match.

### Step 3 — Append the entry

Add (or extend) the relevant top-level section. Example:

```yaml
extensions:
  rs:
    glyph: "U+E7A8"
    color: "#dea584"
  go:
    color: "#ff5555"     # color-only override: keeps built-in Go glyph
```

If the section already exists, append the new key under it — do **not** create a duplicate top-level key.

### Step 4 — Verify

Tell the user to run `logo-ls` in a directory containing a matching file. If the override misbehaves, suggest `logo-ls --override-file <path>` to test an alternate file, or `--no-override` to confirm the issue is in the override (not the built-ins). Parse errors are printed once at startup as `logo-ls: ignoring icon overrides: ...`.

### Step 5 — Offer promotion

If the user iterated on the YAML and seems happy with the result, offer to **promote** the entries to the built-in maps via Procedure B so the icon can be contributed upstream.

---

## Procedure B — Built-in contribution

Follow these steps interactively:

### Step 1 — Gather icon details

Ask the user:

1. **Icon name**: The logical name for this icon in `IconSet` (e.g. `"markdown"`, `"docker"`, `"rust"`). Check `internal/icons/icons_map.go` to see if it already exists. If it does, skip to Step 2.
2. **Nerd Font codepoint**: The Unicode codepoint for the glyph (e.g. `U+E7A8` or `e7a8`). The user should look this up at https://www.nerdfonts.com/cheat-sheet.
3. **RGB color**: Three values 0–255 for red, green, blue (e.g. `66, 165, 245`).

### Step 2 — Determine mapping type

Ask the user which type of mapping they want. They can choose one or more:

- **Extension** — Map one or more file extensions to this icon (added to `internal/icons/icons_ext.go`)
- **File name** — Map one or more specific file names to this icon (added to `internal/icons/icons_files.go`)
- **Sub-extension** — Map one or more compound sub-extensions to this icon (added to `internal/icons/icons_sub_ext.go`)

Then ask for the specific extension(s), file name(s), or sub-extension(s) to map.

### Step 3 — Apply changes

1. If the icon is new, add an entry to `IconSet` in `internal/icons/icons_map.go`. Insert it in **alphabetical order** within the map. Follow the existing formatting style:
   ```go
   "icon-name": {Glyph: "\U0000XXXX", Color: [3]uint8{R, G, B}},
   ```
   The codepoint must be formatted as an 8-digit uppercase Unicode escape (e.g. `\U0000E7A8`).

2. Add entries to the appropriate mapping file(s), also in **alphabetical order**, following the existing formatting style:
   - Extensions in `internal/icons/icons_ext.go`: `"ext": IconSet["icon-name"],`
   - File names in `internal/icons/icons_files.go`: `"filename": IconSet["icon-name"],`
   - Sub-extensions in `internal/icons/icons_sub_ext.go`: `"sub.ext": IconSet["icon-name"],`

3. Verify the changes compile by running `go build ./...`.

### Step 4 — Summary

Show the user a summary of what was added:
- The icon name, glyph codepoint, and color
- All mappings that were created
- Confirm the build succeeded

If you arrived here by promoting a YAML override, also remind the user that they can now remove the matching entries from their personal YAML (the built-in will take over once installed) — but only if they want to; leaving it in place is harmless because user entries always override built-ins.
