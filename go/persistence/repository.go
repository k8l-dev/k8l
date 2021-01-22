package persistence

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// Storage is the object enclosing all the repository
type Storage struct {
	LogRepository *LogRepository
}

// STORAGE is the global object
var STORAGE Storage

// LogRepository is the main interface to save roecords on sqlite
type LogRepository struct {
	Connection *sql.DB
}

// Setup the database
func (r LogRepository) Setup() {
	_, err := r.Connection.Exec(`
	CREATE TABLE IF NOT EXISTS logs (	
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
		namespace_name TEXT NOT NULL,
		container_name TEXT NOT NULL,
		pod_name TEXT NOT NULL,
		container_image TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		message TEXT NOT NULL);
	
	CREATE INDEX IF NOT EXISTS timestamp_idx ON logs(timestamp);

	CREATE VIRTUAL TABLE IF NOT EXISTS logs_fts USING fts5(
		id UNINDEXED,
			namespace_name UNINDEXED,
			container_name,
			pod_name,
			container_image,
			timestamp UNINDEXED,
			message, 
			content='logs'
	);

	CREATE TRIGGER IF NOT EXISTS logs_ai AFTER INSERT ON logs
		BEGIN
			INSERT INTO logs_fts (container_name, pod_name, container_image, message)
			VALUES (new.container_name, new.pod_name, new.container_image, new.message);
		END;

	CREATE TRIGGER IF NOT EXISTS logs_ad AFTER DELETE ON logs
		BEGIN
			INSERT INTO logs_fts (logs_fts, container_name, pod_name, container_image, message)
			VALUES ('delete', old.container_name, old.pod_name, old.container_image, old.message);
		END;

	CREATE TRIGGER IF NOT EXISTS logs_au AFTER UPDATE ON logs
		BEGIN
			INSERT INTO logs_fts (logs_fts, container_name, pod_name, container_image, message)
			VALUES ('delete', old.container_name, old.pod_name, old.container_image, old.message);
			INSERT INTO logs_fts (container_name, pod_name, container_image, message)
			VALUES (new.container_name, new.pod_name, new.container_image, new.message);
		END;
				`)
	if err != nil {
		log.Fatal("Setup failed ", err)
	}
}

// GetLogs returns logs
func (r LogRepository) GetLogs(namespace string, container string) []LogRecord {
	logs := make([]LogRecord, 0)
	rows, err := r.Connection.Query("SELECT namespace_name, container_name, pod_name, container_image, timestamp, message FROM logs WHERE namespace_name = ? AND container_name = ?", namespace, container)
	if err != nil {
		log.Error("Cannot insert log record", err)
	}

	defer rows.Close()
	log.Info(">>>>")
	for rows.Next() {
		var entry LogRecord
		err := rows.Scan(&entry.Namespace, &entry.Container, &entry.Pod,
			&entry.Image, &entry.Timestamp, &entry.Message)
		log.Info("asd", entry)
		if err != nil {
			log.Error(err)
		}
		logs = append(logs, entry)
	}

	return logs
}

// Save a record on sqlite
func (r LogRepository) Save(record LogRecord) bool {
	stmt, err := r.Connection.Prepare(`
		INSERT INTO logs (namespace_name, container_name, pod_name, container_image, "timestamp", message)
		VALUES(?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		log.Error("Cannot insert log record", err)
		return false
	}

	_, err = stmt.Exec(record.Namespace,
		record.Container,
		record.Pod,
		record.Image,
		record.Timestamp,
		record.Message)

	if err != nil {
		log.Error("Cannot insert log record", err)
	}

	return true
}
