# tubu

tubu は Go で実装された Discord Bot です．

## Features

- メッセージリンクの展開機能
- `dice` コマンド

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

2. `.env`を同一のフォルダに作成します．

   | 環境変数名           | 概要                                            | 既定値 |
   | :------------------- | :---------------------------------------------- | :----- |
   | `TUBU_DISCORD_TOKEN` | Discord Bot の API トークンを指定してください． | ---    |
   | `TUBU_API_TIMEOUT`   | API のタイムアウト時間を指定してください．      | `5s`   |

3. `docker compose up`で起動できます．
