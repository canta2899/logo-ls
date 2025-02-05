# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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


