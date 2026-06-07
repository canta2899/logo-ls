# Contributing

Contributions are welcome. Feel free to open an issue or a pull request if you have any questions or want to contribute.

## Getting started

Clone the repository and build the binary:

```bash
git clone https://github.com/canta2899/logo-ls
cd logo-ls
make logo-ls
```

## Making changes

- Open an issue before starting significant work so we can align on the approach.
- Keep pull requests focused on a single concern as much as possible
- Make sure all tests pass before submitting:

```bash
make test
```

- If you add new behaviour, add or update tests to cover it.
- Keep the existing code style and conventions, the project uses standard Go formatting (`gofmt`).

## Adding Icons

If you use any coding agent (OpenCode, Gemini CLI, Claude Code, etc.) there's a built in skill called `/add-icon` which you can use to let your agent do the job for you. If you want to do it manually, you can pretend to be a coding agent and read the skill file yourself at `add-icon-skill.md` in the repo root.

> **Note for Windows contributors:** the paths `.agents/skills/add-icon/SKILL.md` and `.claude/skills/add-icon/SKILL.md` are symlinks to `add-icon-skill.md`. Git for Windows does not create real symlinks by default, so these may be checked out as plain text files containing the link target. To get working symlinks, enable Developer Mode (or run as admin) and set `git config --global core.symlinks true` before cloning.
