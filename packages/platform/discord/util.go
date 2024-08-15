package discord

import (
	"strconv"
	"time"
)

// Discord Epochは，Discord APIのSnowflakeで使用される基準日時
// 基準時刻は，2015/01/01 00:00:00 UTC
// 上位42ビットはタイムスタンプ，次の5ビットはデータセンターID，次の5ビットはワーカーID，最後の12ビットはシーケンス番号
// ref: https://discord.com/developers/docs/reference#snowflakes
const DiscordEpoch = int64(1420070400000)

func TimestampFromSnowflake(id string) time.Time {
	// IDは，確実にsnowflakeで有るため，簡略化の為にエラーチェックを省略
	snowflake, _ := strconv.ParseInt(id, 10, 64)
	return time.Unix(0, ((snowflake>>22)+DiscordEpoch)*int64(time.Millisecond))
}
