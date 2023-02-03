Cross WSL Git (xwslgit)
=======================

Concept
-------

xwslgit is a tiny tool to switch git on Windows and git(s) on Windows Subsystem for Linux (WSL) depending on target directories.

The target directory is the one on Windows (e.g., `C:\...`), xwslgit launches `git.exe` on Windows.
The target directory is the one on WSL (e.g., `\\wsl$\Ubuntu-22.04\...` or `\\wsl.localhost\Ubuntu-22.04\...`, xwslgit launches `wsl -d Ubuntu-22.04 -- git ...`.
