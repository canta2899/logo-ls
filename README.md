# logo-ls

This is a fork of the original [logo-ls repository](https://github.com/Yash-Handa/logo-ls). The modifications I made were

- Removing case sensitive paths in order to avoid issues on OSX/Windows machines
- Chaning UTF-8 space character in order to properly render the output on every terminal
- Fixing a bug on path separators that caused issues with the git status flag on Windows machines
- Implementing a symlink resolution that behaves like the one of the original [ls(coreutils)](https://www.gnu.org/software/coreutils/manual/html_node/ls-invocation.html#ls-invocation).

The original repository has been inactive for more than a year, my pull requests are still pending.
On this repository you can find separate branches for each feature/fix I worked on, if you want to use my fork you can just pull from the current branch (**main**) and run `go install .` inside the project's base directory.
