# onamae-ddns

## Overview

実行すると、お名前.com の 指定のドメインの A レコードを実行マシンの現在のグローバル IP に更新する。

これを cron などで定期的に実行することで IP アドレスが変わった時に自動更新して DDNS を実現できる。

自動更新のクライアントは公式で存在するが Windows 専用であり、これは Linux や Mac でも使える。

## Usage

1. `config.json` の設定を自身のお名前.com のものに書き換える。

1. `go build` でビルドする。

1. できた `onamae-ddns` を cron などで定期的に実行するようにする。

## Notes

`go build` でビルドすると `config.json` の中身は実行バイナリに含まれる。`config.json` はセキュリティの観点からも削除してしまってよい。
