Cross WSL Git (xwslgit)
=======================

[Japanese text version](README_ja.md)

Concept
-------

`xwslgit` is a tiny command line tool to switch Git on Windows and git(s) on Windows Subsystem for Linux (WSL) depending on target directories.

`xwslgit` launches `wsl -d Ubuntu-22.04 --shell-type none -- git ...` if the target directory is the one on WSL (e.g., `\\wsl$\Ubuntu-22.04\...` or `\\wsl.localhost\Ubuntu-22.04\...`).
Otherwise (that is, the target directory is the one on Windows (e.g., `C:\...`) ), `xwslgit` launches `git.exe` of Git for Windows.

Motivation
----------

I want to use [TortoiseGit](https://tortoisegit.org/) even on WSL directories. But TortoiseGit uses [Git for Windows](https://gitforwindows.org/) and running Git for Windows against WSL directories results:

* Really slow. Assumed for it accesses via SMB.
* Permission issues. Executable bits aren't properly set on WSL filesystems.
* Ownership issues. You need to run `git config --global --add safe.directory ...` for WSL directories.

Usage
-----

### Quick usage

#### Tortoisegit users:

1. Put `xwslgit.exe` as `%APPDATA%\xwslgit\git.exe` (be aware renaming to `git.exe`).
2. Set `WSLENV` environment variable to `GIT_SSH/p`.
    * You need to restart the computer to have TortoiseGit read that.
3. Open TortoiseGit settings and set "General > Git.exe Path" to the path you put `git.exe` in step 1.

### Install

You have some options:

* A: Put `xwslgit.exe` to some appropriate path (e.g. `%APPDATA%\xwslgit\xwslgit.exe`). Configure your tools to use that instead of `git.exe` of Git for Windows (e.g. `C:\Program Files\Git\bin\git.exe`).
* B: Rename `xwslgit.exe` to `git.exe` and put it to some appropriate path (e.g. `%APPDATA%\xwslgit\git.exe`). Configure your tools to use that directory instead of the directory of `git.exe` of Git for Windows (e.g. `C:\Program Files\Git\bin`).
* C: Rename `xwslgit.exe` to `git.exe` and put it to some appropriate path (e.g. `%APPDATA%\xwslgit\git.exe`). Configure your system `PATH` environment variable to include that directory. NOTE: You must put that directory in front of the directory of `git.exe` of Git for Windows (e.g. `%APPDATA%\xwslgit;...;C:\Program Files\Git\bin;...`).

### Configuration

You can configure the behavior of xwslgit with `%APPDATA%\xwslgit\xwslgit.yaml`.
See [config/xwslgit.yaml](config/xwslgit.yaml) for the configuration example.

### Environment variables

`xwslgit` doesn't configure nor convert environment variables.
You can do that with https://devblogs.microsoft.com/commandline/share-environment-vars-between-wsl-and-windows/ .

For example, you can configure `WSLENV` like:

```
WSLENV=GIT_SSH/p:GIT_DIR/p:GIT_WORK_TREE/p:GIT_AUTHOR_NAME:GIT_AUTHOR_EMAIL
```

Detailed behavior
-----------------

* How to detect the target WSL distritubion:
    * `detection.useArguments` set to `true`, `xwslgit` detects target WSL distribution from arguments to `git` command.
        * If called with `git clone https://github.com/ikedam/xwslgit \\wsl$\Ubuntu-22.04\home\ikedam\xwslgit`, `xwslgit` considers `Ubuntu-22.04` as the target distribution.
    * `xwslgit` detects target WSL distribution from the current working directory.
* How to call `git` on WSL:
    * Convert paths pointing WSL filesystem in arguments to paths inside WSL.
        * `\\wsl$\Ubuntu-22.04\home\ikedam\xwslgit` is converted to `/home/ikedam/xwslgit`
    * Launch command with `wsl -d Ubuntu-22.04 --shell-type none -- git ...`.

That's all.

Special command
---------------

You can see the version information of `xwslgit` with:

```sh
xwslgit xwslgitversion
```

How to build locally
--------------------

```sh
GOOS=windows GOARCH=amd64 go build -o xwslgit.exe ./cmd/xwslgit
```

Known issues
------------

* Doesn't work with [Sourcetree for Windows](https://www.sourcetreeapp.com/).
    * Sourcetree automatically detects Git for Windows and doesn't accept custom git clients.
    * I couldn't have Sourcetree use `xwslgit`. Even replacing `git.exe` that Sourcetree refers to `xwslgit`, Sourcetree don't accept that as `git.exe`.
