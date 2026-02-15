# clijp

`clijp` は、標準入力で受け取った内容を GitHub Copilot SDK を使用して日本語に翻訳するコマンドラインツールです。コードや記号はそのまま保持し、自然な日本語訳を提供します。

## 前提条件

このツールを使用するには、以下の前提条件を満たす必要があります：

- **GitHub Copilot のライセンス**: GitHub Copilot SDK を利用するため、GitHub Copilot の有効なライセンスが必要です。GitHub Copilot は GitHub アカウントに紐づいており、個人または組織のサブスクリプションが必要です。
- **Go 環境**: Go 1.26.0 以上がインストールされていること。
- **GitHub Copilot の設定**: SDK が正しく動作するよう、GitHub Copilot がインストール・設定されている環境で実行してください。詳細は [GitHub Copilot のドキュメント](https://docs.github.com/en/copilot) を参照してください。

## インストール

```sh
go install github.com/otakakot/clijp@latest
```

## 使用方法

標準入力からテキストを渡して実行します。パイプやリダイレクトを使用して入力してください。

### 使用例

```sh
❯ go help | clijp
Go is a tool for managing Go source code.

Usage:

        go <command> [arguments]

The commands are:

        bug         start a bug report
        build       compile packages and dependencies
        clean       remove object files and cached files
        doc         show documentation for package or symbol
        env         print Go environment information
        fix         apply fixes suggested by static checkers
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         add dependencies to current module and install them
        install     compile and install packages and dependencies
        list        list packages or modules
        mod         module maintenance
        work        workspace maintenance
        run         compile and run Go program
        telemetry   manage telemetry data and settings
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         report likely mistakes in packages

Use "go help <command>" for more information about a command.

Additional help topics:

        buildconstraint build constraints
        buildjson       build -json encoding
        buildmode       build modes
        c               calling between Go and C
        cache           build and test caching
        environment     environment variables
        filetype        file types
        goauth          GOAUTH environment variable
        go.mod          the go.mod file
        gopath          GOPATH environment variable
        goproxy         module proxy protocol
        importpath      import path syntax
        modules         modules, module versions, and more
        module-auth     module authentication using go.sum
        packages        package lists and patterns
        private         configuration for downloading non-public code
        testflag        testing flags
        testfunc        testing functions
        vcs             controlling version control with GOVCS

Use "go help <topic>" for more information about that topic.

=== 日本語翻訳 ===

Goは、Goソースコードを管理するためのツールです。

使用方法:

        go <コマンド> [引数]

コマンド一覧:

        bug         バグレポートを開始する
        build       パッケージと依存関係をコンパイルする
        clean       オブジェクトファイルとキャッシュファイルを削除する
        doc         パッケージまたはシンボルのドキュメントを表示する
        env         Go環境情報を出力する
        fix         静的チェッカーが提案する修正を適用する
        fmt         gofmt（再フォーマット）でパッケージソースを整形する
        generate    ソースを処理してGoファイルを生成する
        get         現在のモジュールに依存関係を追加してインストールする
        install     パッケージと依存関係をコンパイルしてインストールする
        list        パッケージまたはモジュールを一覧表示する
        mod         モジュールのメンテナンス
        work        ワークスペースのメンテナンス
        run         Goプログラムをコンパイルして実行する
        telemetry   テレメトリデータと設定を管理する
        test        パッケージをテストする
        tool        指定されたgoツールを実行する
        version     Goのバージョンを出力する
        vet         パッケージ内の潜在的なミスを報告する

コマンドの詳細情報は "go help <コマンド>" を使用してください。

追加のヘルプトピック:

        buildconstraint ビルド制約
        buildjson       build -json エンコーディング
        buildmode       ビルドモード
        c               GoとC間の呼び出し
        cache           ビルドとテストのキャッシュ
        environment     環境変数
        filetype        ファイルタイプ
        goauth          GOAUTH環境変数
        go.mod          go.modファイル
        gopath          GOPATH環境変数
        goproxy         モジュールプロキシプロトコル
        importpath      インポートパス構文
        modules         モジュール、モジュールバージョンなど
        module-auth     go.sumを使用したモジュール認証
        packages        パッケージリストとパターン
        private         非公開コードのダウンロード設定
        testflag        テストフラグ
        testfunc        テスト関数
        vcs             GOVCSによるバージョン管理の制御

トピックの詳細情報は "go help <トピック>" を使用してください。
```

```sh
❯ go doc fmt.Sprintf | clijp
package fmt // import "fmt"

func Sprintf(format string, a ...any) string
    Sprintf formats according to a format specifier and returns the resulting
    string.

=== 日本語翻訳 ===

package fmt // import "fmt"

func Sprintf(format string, a ...any) string
    Sprintf はフォーマット指定子に従って書式化し、結果の文字列を返します。
```

## ライセンス

このプロジェクトは MIT License の下で公開されています。詳細は [LICENSE](LICENSE) ファイルを参照してください。
