# tubu

tubu は Go で実装された Discord Bot です．

## Features

- メッセージリンクの展開機能
- `channel`コマンド
- `dice` コマンド
- `guild`コマンド
- `version`コマンド
- メッセージのピン留め機能

## Deployment

> [!IMPORTANT]
> 現在，tubu は Docker compose によるホスティングのみを公式にサポートしています．

Docker Compose が導入されていることを確認してください．

1. 以下のような`compose.yml`を用意します．

   ```yaml
   services:
     bot:
       container_name: tubu-bot
       image: ghcr.io/aqyuki/tubu:latest
       env_file:
         - .env
       restart: unless-stopped
   ```

2. `.env`を同一のフォルダに作成します．設定できるオプションについては[Options](#options)を確認してください．

3. `docker compose up`で起動できます．

## Options

**tubu**の設定は，環境変数もしくは起動オプションから変更できます．

| 設定名            | 概要                                          | 環境変数名           | フラグ名  | 既定値 | 必須  |
| :---------------- | :-------------------------------------------- | :------------------- | :-------- | :----: | :---: |
| **DISCORD_TOKEN** | Discord Bot の APi トークンを指定してください | `TUBU_DISCORD_TOKEN` | `token`   |  ---   | **◯** |
| **API_TIMEOUT**   | API のタイムアウトを指定してください          | `TUBU_TIMEOUT`       | `timeout` | `10s`  |       |
