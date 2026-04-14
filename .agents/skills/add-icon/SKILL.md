---
name: add-icon
description: Add a new icon to logo-ls by updating the icon maps
---

# Add Icon to logo-ls

You are helping the user add a new icon to the logo-ls project. This project uses Nerd Fonts glyphs to display icons next to file and directory names in terminal listings.

## Architecture

The icon system consists of four files in the `icons/` directory:

- **`icons/icons_map.go`** — `IconSet` map: defines all available icons. Each entry has a name (key), a glyph (Unicode codepoint), and an RGB color. This is where every new icon must be registered first.
- **`icons/icons_ext.go`** — `IconExt` map: maps file extensions (e.g. `"js"`, `"py"`) to entries in `IconSet`.
- **`icons/icons_files.go`** — `IconFileName` map: maps specific file names (e.g. `"tsconfig.json"`, `".babelrc"`) to entries in `IconSet`.
- **`icons/icons_sub_ext.go`** — `IconSubExt` map: maps compound sub-extensions (e.g. `"spec.ts"`, `"d.ts"`, `"gitlab-ci.yml"`) to entries in `IconSet`. Used to override the default extension icon for specific patterns.

The `IconInfo` struct looks like this:

```go
type IconInfo struct {
    Glyph        string
    Color        [3]uint8 // RGB color
    IsExecutable bool
}
```

## Procedure

Follow these steps interactively:

### Step 1 — Gather icon details

Ask the user:

1. **Icon name**: The logical name for this icon in `IconSet` (e.g. `"markdown"`, `"docker"`, `"rust"`). Check `icons/icons_map.go` to see if it already exists. If it does, skip to Step 2.
2. **Nerd Font codepoint**: The Unicode codepoint for the glyph (e.g. `U+E7A8` or `e7a8`). The user should look this up at https://www.nerdfonts.com/cheat-sheet.
3. **RGB color**: Three values 0–255 for red, green, blue (e.g. `66, 165, 245`).

### Step 2 — Determine mapping type

Ask the user which type of mapping they want. They can choose one or more:

- **Extension** — Map one or more file extensions to this icon (added to `icons/icons_ext.go`)
- **File name** — Map one or more specific file names to this icon (added to `icons/icons_files.go`)
- **Sub-extension** — Map one or more compound sub-extensions to this icon (added to `icons/icons_sub_ext.go`)

Then ask for the specific extension(s), file name(s), or sub-extension(s) to map.

### Step 3 — Apply changes

1. If the icon is new, add an entry to `IconSet` in `icons/icons_map.go`. Insert it in **alphabetical order** within the map. Follow the existing formatting style:
   ```go
   "icon-name": {Glyph: "\U0000XXXX", Color: [3]uint8{R, G, B}},
   ```
   The codepoint must be formatted as an 8-digit uppercase Unicode escape (e.g. `\U0000E7A8`).

2. Add entries to the appropriate mapping file(s), also in **alphabetical order**, following the existing formatting style:
   - Extensions in `icons/icons_ext.go`: `"ext": IconSet["icon-name"],`
   - File names in `icons/icons_files.go`: `"filename": IconSet["icon-name"],`
   - Sub-extensions in `icons/icons_sub_ext.go`: `"sub.ext": IconSet["icon-name"],`

3. Verify the changes compile by running `go build ./...`.

### Step 4 — Summary

Show the user a summary of what was added:
- The icon name, glyph codepoint, and color
- All mappings that were created
- Confirm the build succeeded
