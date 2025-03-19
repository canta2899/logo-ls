# logo-ls

This is a fork of [logo-ls](https://github.com/Yash-Handa/logo-ls) which I ended up maintaining since the original repository went unmaintained some years ago. Check the [Installation Guide](#installation) if you want to install it.

## What I did

I fixed several bugs and implemented some additional features, including Nerd Fonts v3 migration and symlinks count. You can check a complete list of features and fixes on the `CHANGELOG.md` file.

## What I plan to do

- Implement custom time stamp formatting
- Implement filter by regex pattern to avoid listing unwanted files or subdirectories
- Show size of directories as the sum of the sizes of their contents (this would be different than the behavior of _ls(coreutils)_)
- Add more tests
- Deploy to package managers
- Introduce a configuration file to customize the output

Feel free to contribute and I'll be more than happy to merge your changes. In case you want to add some new icons, please make a PR so that we can all benefit from that. The following sections explains how to do it.

---

## Adding Icons

In order to add new icons you need to:

1. Fork the repository
2. Add an entry in the map `IconSet` in `icons/icons_map.go`. The key should be the **name of the icon** and **not** its extension (i.e. `markdown` for markdown files). The value should be of type `IconInfo`, a struct indicating the unicode character and its color.
3. Add one or more entries in the map `IconExt` in `icons/icon_ext.go`, mapping the entry defined in `IconSet` to each one of the desired extensions

You can map an icon to a specific file name (i.e. `tsconfig.json`) by editing `icons/icons_files.go`. You can also override an icon for a specific file name (i.e. use a different icon for `gitlab-ci.yml` rather than the standard YML one) by editing `icons/icons_sub_ext.go`.

---

## Installation

### You don't want to build it yourself

If you are on Linux or OSX you can run this script on your shell to download the latest version and move it to ~/.local/bin:

```bash
curl -L https://raw.githubusercontent.com/canta2899/logo-ls/refs/heads/main/get.sh | sh
```

You can do the same on Windows by running the following command in powershell

```powershell
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
irm https://raw.githubusercontent.com/canta2899/logo-ls/main/get.ps1 | iex
```

Otherwise, you can manually download the binary for your platform from the [releases page](https://github.com/canta2899/logo-ls/releases/). Then, you have to extract the archive (md5 checksums are provided in case you want to verify the integrity of the file) and move the executable binary of logo-ls to a directory in your `$PATH` (or symlink it).

In case you want to replace the original `ls` command with `logo-ls`, I would suggest **adding an alias**:

On OSX/Linux

```bash
alias ls="logo-ls"
```

On Windows (powershell)
 
```powershell
Set-Alias ls logo-ls
```

### You have go installed and want to build it yourself

Clone the repository

```bash
git clone https://github.com/canta2899/logo-ls
```

If you want to install directly to your `$GOPATH` you can use `go install`.

```bash
go install ./cmd/logo-ls
```

If you want to build the binary you can use `make`:

```bash
# outputs executable 'logo-ls' in the root directory
make logo-ls
```

```bash
# cleans up the executable from the repo
make clean
```
