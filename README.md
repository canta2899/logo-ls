# logo-ls

This is a fork of [logo-ls](https://github.com/Yash-Handa/logo-ls) which I ended up maintaining since the original repository went unmaintained some years ago.

My primary goals have been:

- Correct output (I should see all the files /bin/ls would show me)
- Informative output (I should see at least what /bin/ls would show me)
- MacOS / Linux / Windows compatibility (I mainly develop on OSX but I want to use logo-ls everywhere)

The modifications I made involved:

## What I did

- A lot of refactoring, at the point where the codebase is very different from the original one
- Removing files with case sensitive naming which where causing issues when cloning the repo on Windows/MacOs
- Changing UTF-8 space character in order to properly render the output on every terminal
- Fixing a bug on path separators that caused issues with the git status flag on Windows machines
- Implementing a symlink resolution that behaves like the one of the original [ls(coreutils)](https://www.gnu.org/software/coreutils/manual/html_node/ls-invocation.html#ls-invocation).
- Nerd Fonts Version 3 migration (this breaks compatibility with codepoints of previous versions)
- Added hard link count
- Added inode numbers
- Updated CLI args parsing so that it mostly behaves like the one of ls(coreutils)
- Re-implemented the git status feature from scratch reducing build size and improving performance
- Added extended attributes support
- Added sticky bit support

## What I plan to do

- Implement custom time stamp formatting
- Implement filter by regex pattern to avoid listing unwanted files or subdirectories
- Show size of directories as the sum of the sizes of their contents (this would be different than the behavior of _ls(coreutils)_)
- Add more tests
- Deploy to package managers
- Introduce a configuration file to customize the output

Feel free to contribute and I'll be more than happy to merge your changes. In case you want to add some new icons, please make a PR so that we can all benefit from that. The following sections explains how to do it.

---

# Adding Icons

In order to add new icons you need to:

1. Fork the repository
2. Add an entry in the map `IconSet` in `icons/icons_map.go`. The key should be the **name of the icon** and **not** its extension (i.e. `markdown` for markdown files). The value should be of type `IconInfo`, a struct indicating the unicode character and its color.
3. Add one or more entries in the map `IconExt` in `icons/icon_ext.go`, mapping the entry defined in `IconSet` to each one of the desired extensions

You can map an icon to a specific file name (i.e. `tsconfig.json`) by editing `icons/icons_files.go`. You can also override an icon for a specific file name (i.e. use a different icon for `gitlab-ci.yml` rather than the standard YML one) by editing `icons/icons_sub_ext.go`.

---

## Installation

Since no binary is available, you should compile yours by yourself (for now). You have to [install go](https://go.dev/doc/install), then run:

```bash
git clone https://github.com/canta2899/logo-ls
cd logo-ls
```

If you want to install directly to your `$GOPATH` you can use `go install`.

```bash
go install ./cmd/logo-ls
```

If you want to build the binary you can use `make`.

```bash
make logo-ls # outputs executable in bin/logo-ls
```

If you are not able to run make (i.e. you are on Windows) you can build the binary using the following command:

```bash
go build -o logo-ls -ldflags="-s -w" -tags=minimal -trimpath ./cmd/logo-ls
```
