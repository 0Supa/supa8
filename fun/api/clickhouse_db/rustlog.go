package logs_db

import (
	"context"
	logger "log"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
)

var log = logger.New(os.Stdout, "Clickhouse ", logger.LstdFlags)

var (
	ctx             = context.Background()
	Clickhouse, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "rustlog",
			Username: "default",
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "supa8", Version: "0.1"},
			},
		},
	})
)

func init() {
	if err != nil {
		log.Println(err)
	}

	if err := Clickhouse.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
			return
		}
		log.Println(err)
	}
}
