Cross WSL Git (xwslgit)
=======================

概要
----

`xwslgit` は自動的に Git on Windows と、 Windows Subsystem for Linux (WSL) 上の `git` コマンドを呼び分ける簡単なコマンドラインツールです。

`xwslgit` は呼び出し対象のパスが WSL 上のパス (`\\wsl$\Ubuntu-22.04\...` または `\\wsl.localhost\Ubuntu-22.04\...`) の場合、 `wsl -d Ubuntu-22.04 --shell-type none -- git ...` というコマンドを呼び出します。
そうでない場合、つまり対象ディレクトリが `C:\...` などの Windows のものの場合は、 `xwslgit` は Git for Windows の `git.exe` を呼び出します。

背景
----

[TortoiseGit](https://tortoisegit.org/) を WSL のディレクトリ上で使用すると、 [Git for Windows](https://gitforwindows.org/) が呼び出され、以下のような難点があります:

* めちゃ遅い。多分 Windows ファイル共有 (SMB) 経由でのアクセスになるせい。
* ファイルの権限がおかしくなる。特に実行ビットがちゃんと WSL 上で設定されない。
* ファイルの所有者の問題。 `git config --global --add safe.directory ...` を呼び出さないとちゃんと動かない。

[wslgit](https://github.com/andy-5/wslgit) も解決策になりますが、Windows ファイルシステム上のリポジトリに対しても WSL で `git` を起動してしまいます。

使用方法
--------

`xwslgit` のビルド済みバイナリを [Releases](https://github.com/ikedam/xwslgit/releases) からダウンロードできます。

### クイックスタート

#### TortoiseGit のユーザー向け

1. `xwslgit.exe` を `%APPDATA%\xwslgit\git.exe` として設置してください (ファイル名を `git.exe` に変更していることに注意)。
2. システム設定から環境変数 `WSLENV` に `GIT_SSH/p` を設定してください。
    * 設定を TortoiseGit に反映するのにコンピューターの再起動が必要です。
3. TortoiseGit の設定で、 "General > Git.exe Path" に手順 1 で `git.exe` を設置したディレクトリを指定してください。

### インストール

いくつか方法があります。

* A: `xwslgit.exe` を適当な場所、例えば `%APPDATA%\xwslgit\xwslgit.exe` に置きます。Git を使うツールで Git for Windows の `git.exe` (`C:\Program Files\Git\bin\git.exe` など) の代わりに、今回置いたパスを指定します。
* B: `xwslgit.exe` を `git.exe` にリネームして、適当な場所、例えば `%APPDATA%\xwslgit\git.exe` に置きます。Git を使うツールで Git for Windows のインストールディレクトリ (`C:\Program Files\Git` など) の代わりに、今回置いたディレクトリを指定します。
* C: `xwslgit.exe` を `git.exe` にリネームして、適当な場所、例えば `%APPDATA%\xwslgit\git.exe` に置きます。システムの `PATH` 環境変数に、今回置いたディレクトリが含まれるようにします。Git for Windows の `git.exe` があるディレクトリよりも前に指定する必要があることに注意してください。例えば `%APPDATA%\xwslgit:...;C:\Program Files\Git\bin;...` といった設定になります。

### 設定

`%APPDATA%\xwslgit\xwslgit.yaml` で動作を設定できます。
[config/xwslgit.yaml](config/xwslgit.yaml) に設定ファイルの例があるので、参照してください。

### 環境変数について

`xwslgit` は環境変数の面倒を見ません。
代わりに https://devblogs.microsoft.com/commandline/share-environment-vars-between-wsl-and-windows/ を参照して、WSL の機能で環境変数を扱ってください。

例えば、 `WSLENV` 環境変数を以下のように設定します:

```
WSLENV=GIT_SSH/p:GIT_DIR/p:GIT_WORK_TREE/p:GIT_AUTHOR_NAME:GIT_AUTHOR_EMAIL
```

動作の詳細
----------

* どのように対象の WSL ディストリビューションを判定しているか:
    * `detection.useArguments` が `true` に設定されている場合、 `xwslgit` は `git` コマンドに対する引数から WSL ディストリビューションを判定します。
        * 例えば `git clone https://github.com/ikedam/xwslgit \\wsl$\Ubuntu-22.04\home\ikedam\xwslgit` として呼ばれた場合、 `xwslgit` は `Ubuntu-22.04` を対象のディストリビューションと判断します。
    * `xwslgit` は現在の作業ディレクトリから対象の WSL ディストリビューションを判断します。
* WSL 上の `git` の呼び方:
    * コマンド引数のうちの WSL 上のパスを、WSL 内でのパスに変換します:
        * `\\wsl$\Ubuntu-22.04\home\ikedam\xwslgit` は `/home/ikedam/xwslgit` に変換されます。
    * コマンドを `wsl -d Ubuntu-22.04 --shell-type none -- git ...` で起動します。

そんだけ。

WSL で起動するコマンドのカスタマイズ
------------------------------------

WSL で git コマンドを起動する方法をカスタマイズできます。
例えば `direnv` を使う場合は以下のような設定を行います (テストしてない):

```yaml
distributions:
  Ubuntu-22.04:
    command:
      - wsl
      - -d
      - Ubuntu-22.04
      - --shell-type
      - none
      - --
      - direnv
      - exec
      - .
      - git
```

特別なコマンド
--------------

`xwslgit` のバージョン情報を以下のコマンドで確認できます:

```sh
xwslgit xwslgitversion
```

ローカルでのビルド方法
----------------------

```sh
GOOS=windows GOARCH=amd64 go build -o xwslgit.exe ./cmd/xwslgit
```

既知の問題
----------

* [Sourcetree for Windows](https://www.sourcetreeapp.com/) では使用できません。
    * Sourcetree は自動的に Git for Windows を検出して使用し、カスタムの Git クライアントを使用することができません。
    * いかなる手段を持ってしても `xwslgit` を Sourcetree で使用することはできませんでした。Sourcetree が参照している `git.exe` を `xwslgit` に置き換えても、 Sourcetree はなぜかそれを `git.exe` として扱うことはありませんでした。

ライセンス
----------

* `xwslgit` は [MIT License](LICENSE) で配布しています。
* `xwslgit logo` は [Git Logo](https://git-scm.com/downloads/logos) をベースに [Creative Commons Attribution 3.0 Unported License](https://creativecommons.org/licenses/by/3.0/) で配布しています。
