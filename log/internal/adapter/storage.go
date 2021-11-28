package adapter

import (
	"strings"

	"github.com/gocql/gocql"
	"github.com/okmaki/node/log/internal/core"
)

type StorageAdapter struct {
	session *gocql.Session
}

func NewStorageAdapter(hosts []string) (*StorageAdapter, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "log"

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &StorageAdapter{session: session}, nil
}

func (db *StorageAdapter) Close() {
	db.session.Close()
}

func (db *StorageAdapter) Record(log core.Log) error {
	source_query := `
	INSERT INTO log.source_logs (source_id, source_type, transaction_id, transaction_timestamp, level, timestamp, location, data)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	if err := db.session.Query(source_query, log.SourceId, log.SourceType, log.TransactionId, log.TransactionTimestamp, log.Level, log.Timestamp, log.Location, log.Data).Exec(); err != nil {
		return err
	}

	transaction_query := `
	INSERT INTO log.transaction_logs (source_id, source_type, transaction_id, transaction_timestamp, level, timestamp, location, data)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	if err := db.session.Query(transaction_query, log.SourceId, log.SourceType, log.TransactionId, log.TransactionTimestamp, log.Level, log.Timestamp, log.Location, log.Data).Exec(); err != nil {
		return err
	}

	return nil
}

func (db *StorageAdapter) SearchBySource(filter core.LogFilter, limit int) ([]core.Log, error) {
	var queryBuilder strings.Builder
	params := make([]interface{}, 0, 9)

	queryBuilder.WriteString(`
	SELECT source_id, source_type, transaction_id, transaction_timestamp, level, timestamp, location, data
	FROM log.source_logs
	WHERE source_type = ?`)
	params = append(params, filter.SourceType)

	if filter.After > 0 {
		queryBuilder.WriteString(`
		AND timestamp > ?`)
		params = append(params, filter.After)
	}

	if filter.Before > 0 {
		queryBuilder.WriteString(`
		AND timestamp < ?`)
		params = append(params, filter.Before)
	}

	if filter.SourceId != "" {
		queryBuilder.WriteString(`
		AND source_id = ?`)
		params = append(params, filter.SourceId)
	}

	addLevelsCheckToQuery(filter, &queryBuilder, &params)

	queryBuilder.WriteString(`
	LIMIT ?
	ALLOW FILTERING`)
	params = append(params, limit)

	return executeQuery(db, queryBuilder.String(), &params, limit)
}

func (db *StorageAdapter) SearchByTransaction(filter core.LogFilter, limit int) ([]core.Log, error) {
	var queryBuilder strings.Builder
	params := make([]interface{}, 0, 9)

	queryBuilder.WriteString(`
	SELECT source_id, source_type, transaction_id, transaction_timestamp, level, timestamp, location, data
	FROM log.transaction_logs
	WHERE transaction_id = ?`)
	params = append(params, filter.TransactionId)

	if filter.After > 0 {
		queryBuilder.WriteString(`
		AND timestamp > ?`)
		params = append(params, filter.After)
	}

	if filter.Before > 0 {
		queryBuilder.WriteString(`
		AND timestamp < ?`)
		params = append(params, filter.Before)
	}

	if filter.SourceId != "" {
		queryBuilder.WriteString(`
		AND source_id = ?`)
		params = append(params, filter.SourceId)
	} else if filter.SourceType != "" {
		queryBuilder.WriteString(`
		AND source_type = ?`)
		params = append(params, filter.SourceType)
	}

	addLevelsCheckToQuery(filter, &queryBuilder, &params)

	queryBuilder.WriteString(`
	LIMIT ?
	ALLOW FILTERING`)
	params = append(params, limit)

	return executeQuery(db, queryBuilder.String(), &params, limit)
}

// ------------------------------
// helpers
// ------------------------------

func addLevelsCheckToQuery(filter core.LogFilter, queryBuilder *strings.Builder, params *[]interface{}) {
	levels := make([]core.LogLevel, 0, 4)

	if filter.HasLevel(core.LevelDebug) {
		levels = append(levels, core.LevelDebug)
		*params = append(*params, core.LevelDebug)
	}

	if filter.HasLevel(core.LevelInfo) {
		levels = append(levels, core.LevelInfo)
		*params = append(*params, core.LevelInfo)
	}

	if filter.HasLevel(core.LevelWarning) {
		levels = append(levels, core.LevelWarning)
		*params = append(*params, core.LevelWarning)
	}

	if filter.HasLevel(core.LevelError) {
		levels = append(levels, core.LevelError)
		*params = append(*params, core.LevelError)
	}

	levelsCount := len(levels)

	if levelsCount == 1 {
		queryBuilder.WriteString(`
		AND level = ?`)
	} else if levelsCount > 1 {
		queryBuilder.WriteString(`
		AND level IN (`)

		for i := 0; i < levelsCount; i++ {
			if i == 0 {
				queryBuilder.WriteString("?")
			} else {
				queryBuilder.WriteString(", ?")
			}
		}

		queryBuilder.WriteString(")")
	}
}

func executeQuery(db *StorageAdapter, query string, params *[]interface{}, limit int) ([]core.Log, error) {
	scanner := db.session.Query(query, (*params)...).Iter().Scanner()

	logs := make([]core.Log, 0, limit)
	var err error = nil
	for scanner.Next() {
		log := core.Log{}
		err = scanner.Scan(&(log.SourceId), &(log.SourceType), &(log.TransactionId), &(log.TransactionTimestamp), &(log.Level), &(log.Timestamp), &(log.Location), &(log.Data))

		if err != nil {
			break
		}

		logs = append(logs, log)
	}

	return logs, err
}
