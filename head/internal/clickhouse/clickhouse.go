package clickhouse

import (
	"database/sql"
	"github.com/ClickHouse/clickhouse-go/v2"
	"time"
)

type Client interface {
	GetRow(userID, videoID uint32) (UserVideoTimesRow, error)
}

type UserVideoTimesRow struct {
	UserID         uint32
	EventTime      time.Time
	EventType      string //Enum8('watch' = 0, 'scroll' = 1),
	VideoID        uint32
	VideoTimestamp uint32
}

type clickhouseClient struct {
	conn *sql.DB
}

func (c clickhouseClient) GetRow(userID, videoID uint32) (res UserVideoTimesRow, err error) {
	row := c.conn.QueryRow(
		"SELECT * FROM user_video_times WHERE user_id=$1 AND video_id=$2 ORDER BY event_time DESC LIMIT 1",
		userID, videoID)

	err = row.Scan(&res.UserID, &res.EventTime, &res.EventType, &res.VideoID, &res.VideoTimestamp)
	return
}

func Init(host string) Client {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{host},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug: true,
	})
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)
	return clickhouseClient{conn: conn}
}
