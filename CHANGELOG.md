# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## logo-ls [1.6.0]

### Features

- Icon overrides allow users to extend or change icons through a YAML configuration file (#35)

### Fixes

- Refactored the codebase and reduce the amount of syscalls made by `logo-ls` (measured a 25% faster average execution time) (#34)

## logo-ls [1.5.2]

### Features

- Updated `.log` icon to a one that renders after codepoint migration

## logo-ls [1.5.1]

- New CI/CD which builds binaries and ships them to Homebrew

## logo-ls [1.5.0]

### Features

- Sorting honors LC_COLLATE variable, so that logo-ls behaves more like coreutils ls
- Removed force pinning of . and ..
- AUR package ownership taken by @lcian
- Better testing with CI setup

### Fixes

- Changed dotfiles grouping so that it respects `-X` flag and keeps them grouped together after extensionless files/dirs

## logo-ls [1.4.3]

### Fixes

- Restore coreutils ls -X behavior: extensionless files grouped first (regression from f8bf870)

## logo-ls [1.4.2]

### Fixes

- logo-ls now keeps dotfiles and dotfolders on top of the list regardless of the sorting method used, except "none" and "mod time"

## logo-ls [1.4.1]

### Fixes

- Migrated to internal args parser with no external dependencies
- Improved error handling

## logo-ls [1.4.0]

### Features

- Symlink resolution
- Nerd Fonts Version 3 migration
- Added hard link count
- Added inode numbers
- Updated CLI args parsing so that arguments can be passed in any order and in between flags
- Re-implemented the git status feature from scratch reducing build size and improving performance
- Added extended attributes support
- Added sticky bit support

### Fixes

- UTF-8 space character was changed in order to properly render the output on every terminal
- Correct path separators are used on Windows machines to correctly compute git status


