```txt
List information about the FILEs with ICONS and GIT STATUS (the current dir 
by default). Sort entries alphabetically if none of -tvSUX is specified.
Usage: logo-ls [-1?aAcdDgGhiloRrSstUvVX] [-T value] [files ...]
 -1                    list one file per line.
 -?, --help            display this help and exit
 -a, --all             do not ignore entries starting with .
 -A, --almost-all      do not list implied . and ..
 -d, --directory       list directories themselves, not their contents
 -D, --git-status      print git status of files
 -g                    like -l, but do not list owner
 -G, --no-group        in a long listing, don't print group names
 -h, --human-readable
                       with -l and -s, print sizes like 1K 234M 2G etc.
 -e, --disable-icon    don't print icons of the files
 -i, --inode           print the inode number of each file
 -l                    use a long listing format
 -o                    like -l, but do not list group information
 -R, --recursive       list subdirectories recursively
 -r, --reverse         reverse order while sorting
 -S                    sort by file size, largest first
 -s, --size            print the allocated size of each file, in blocks
 -t                    sort by modification time, newest first
 -T, --time-style      display complete time information
 -U                    do not sort; list entries in directory order
 -v                    natural sort of (version) numbers within text
 -V, --version         output version information and exit
 -X                    sort alphabetically by entry extension

Exit status:
 0  if OK,
 1  if minor problems (e.g., cannot access subdirectory),
 2  if serious trouble (e.g., cannot access command-line argument).
```
