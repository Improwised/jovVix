package database

import (
	stdSQL "database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

// custom schema for mysql database
// source: https://github.com/ThreeDotsLabs/watermill-sql/blob/master/pkg/sql/schema_adapter_mysql.go
type MySQLSchema struct {
	GenerateMessagesTableName func(topic string) string
	SubscribeBatchSize        int
}

func (s MySQLSchema) SchemaInitializingQueries(topic string) []string {
	createMessagesTable := strings.Join([]string{
		"CREATE TABLE IF NOT EXISTS " + s.MessagesTable(topic) + " (",
		"`offset` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,",
		"`uuid` VARCHAR(36) NOT NULL,",
		"`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,",
		"`payload` LONGBLOB DEFAULT NULL,",
		"`metadata` JSON DEFAULT NULL",
		");",
	}, "\n")

	return []string{createMessagesTable}
}

func (s MySQLSchema) InsertQuery(topic string, msgs message.Messages) (string, []interface{}, error) {
	insertQuery := fmt.Sprintf(
		`INSERT INTO %s (uuid, payload, metadata) VALUES %s`,
		s.MessagesTable(topic),
		strings.TrimRight(strings.Repeat(`(?,?,?),`, len(msgs)), ","),
	)

	args, err := defaultInsertArgs(msgs)
	if err != nil {
		return "", nil, err
	}

	return insertQuery, args, nil
}

func (s MySQLSchema) batchSize() int {
	if s.SubscribeBatchSize == 0 {
		return 100
	}

	return s.SubscribeBatchSize
}

func (s MySQLSchema) SelectQuery(topic string, consumerGroup string, offsetsAdapter sql.OffsetsAdapter) (string, []interface{}) {
	nextOffsetQuery, nextOffsetArgs := offsetsAdapter.NextOffsetQuery(topic, consumerGroup)
	selectQuery := `
		SELECT offset, uuid, payload, metadata FROM ` + s.MessagesTable(topic) + `
		WHERE
			offset > (` + nextOffsetQuery + `)
		ORDER BY
			offset ASC
		LIMIT ` + fmt.Sprintf("%d", s.batchSize())

	return selectQuery, nextOffsetArgs
}

func (s MySQLSchema) UnmarshalMessage(row sql.Scanner) (sql.Row, error) {
	r := sql.Row{}
	err := row.Scan(&r.Offset, &r.UUID, &r.Payload, &r.Metadata)
	if err != nil {
		return sql.Row{}, errors.Wrap(err, "could not scan message row")
	}

	msg := message.NewMessage(string(r.UUID), r.Payload)

	if r.Metadata != nil {
		err = json.Unmarshal(r.Metadata, &msg.Metadata)
		if err != nil {
			return sql.Row{}, errors.Wrap(err, "could not unmarshal metadata as JSON")
		}
	}

	r.Msg = msg

	return r, nil
}

func (s MySQLSchema) MessagesTable(topic string) string {
	if s.GenerateMessagesTableName != nil {
		return s.GenerateMessagesTableName(topic)
	}
	return fmt.Sprintf("`watermill_%s`", topic)
}

func (s MySQLSchema) SubscribeIsolationLevel() stdSQL.IsolationLevel {
	// MySQL requires serializable isolation level for not losing messages.
	return stdSQL.LevelSerializable
}

func MysqlDBConnection(cfg config.Sql) (*stdSQL.DB, error) {
	dbURL := cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + ")/" + cfg.Db

	db, err := stdSQL.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}

	return db, err
}
