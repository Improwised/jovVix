package database

import (
	stdSQL "database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Improwised/jovvix/api/config"
	"github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/pkg/errors"

	"github.com/ThreeDotsLabs/watermill/message"
)

// custom schema for postgres database
// source: https://github.com/ThreeDotsLabs/watermill-sql/blob/master/pkg/sql/schema_adapter_postgresql.go

type PostgreSQLSchema struct {
	GenerateMessagesTableName func(topic string) string
	SubscribeBatchSize        int
}

func (s PostgreSQLSchema) SchemaInitializingQueries(topic string) []string {
	createMessagesTable := `
		CREATE TABLE IF NOT EXISTS ` + s.MessagesTable(topic) + ` (
			"offset" SERIAL,
			"uuid" VARCHAR(36) NOT NULL,
			"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"payload" bytea DEFAULT NULL,
			"metadata" JSON DEFAULT NULL,
			"transaction_id" xid8 NOT NULL,
			PRIMARY KEY ("transaction_id", "offset")
		);
	`

	return []string{createMessagesTable}
}
func defaultInsertArgs(msgs message.Messages) ([]interface{}, error) {
	var args []interface{}
	for _, msg := range msgs {
		metadata, err := json.Marshal(msg.Metadata)
		if err != nil {
			return nil, errors.Wrapf(err, "could not marshal metadata into JSON for message %s", msg.UUID)
		}

		args = append(args, msg.UUID, msg.Payload, metadata)
	}

	return args, nil
}

func (s PostgreSQLSchema) InsertQuery(topic string, msgs message.Messages) (string, []interface{}, error) {
	insertQuery := fmt.Sprintf(
		`INSERT INTO %s (uuid, payload, metadata, transaction_id) VALUES %s`,
		s.MessagesTable(topic),
		defaultInsertMarkers(len(msgs)),
	)

	args, err := defaultInsertArgs(msgs)
	if err != nil {
		return "", nil, err
	}
	// log.Fatal("insertQuery", args)
	return insertQuery, args, nil
}

func defaultInsertMarkers(count int) string {
	result := strings.Builder{}

	index := 1
	for i := 0; i < count; i++ {
		result.WriteString(fmt.Sprintf("($%d,$%d,$%d,pg_current_xact_id()),", index, index+1, index+2))
		index += 3
	}

	return strings.TrimRight(result.String(), ",")
}

func (s PostgreSQLSchema) batchSize() int {
	if s.SubscribeBatchSize == 0 {
		return 100
	}

	return s.SubscribeBatchSize
}

func (s PostgreSQLSchema) SelectQuery(topic string, consumerGroup string, offsetsAdapter sql.OffsetsAdapter) (string, []interface{}) {
	// Query inspired by https://event-driven.io/en/ordering_in_postgres_outbox/

	nextOffsetQuery, nextOffsetArgs := offsetsAdapter.NextOffsetQuery(topic, consumerGroup)
	selectQuery := `
		WITH last_processed AS (
			` + nextOffsetQuery + `
		)

		SELECT "offset", transaction_id, uuid, payload, metadata FROM ` + s.MessagesTable(topic) + `

		WHERE
		(
			(
				transaction_id = (SELECT last_processed_transaction_id FROM last_processed)
				AND
				"offset" > (SELECT offset_acked FROM last_processed)
			)
			OR
			(transaction_id > (SELECT last_processed_transaction_id FROM last_processed))
		)
		AND
			transaction_id < pg_snapshot_xmin(pg_current_snapshot())
		ORDER BY
			transaction_id ASC,
			"offset" ASC
		LIMIT ` + fmt.Sprintf("%d", s.batchSize())

	return selectQuery, nextOffsetArgs
}

func (s PostgreSQLSchema) UnmarshalMessage(row sql.Scanner) (sql.Row, error) {
	r := sql.Row{}
	var transactionID int64

	err := row.Scan(&r.Offset, &transactionID, &r.UUID, &r.Payload, &r.Metadata)
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
	r.ExtraData = map[string]any{
		"transaction_id": transactionID,
	}

	return r, nil
}

func (s PostgreSQLSchema) MessagesTable(topic string) string {
	if s.GenerateMessagesTableName != nil {
		return s.GenerateMessagesTableName(topic)
	}
	return fmt.Sprintf(`"watermill_%s"`, topic)
}

func (s PostgreSQLSchema) SubscribeIsolationLevel() stdSQL.IsolationLevel {
	// For Postgres Repeatable Read is enough.
	return stdSQL.LevelSerializable
}

func PostgresDBConnection(cfg config.Sql) (*stdSQL.DB, error) {
	dbURL := "postgres://" + cfg.Username + ":" + cfg.Password + "@" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + "/" + cfg.Db + "?" + cfg.QueryString

	db, err := stdSQL.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, err

}
