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

// GetNamespaces returns all the namespaces sotred
func (r LogRepository) GetNamespaces() []string {
	namespaces := make([]string, 0)
	rows, err := r.Connection.Query("SELECT DISTINCT namespace_name FROM logs order by namespace_name")
	if err != nil {
		log.Error("Cannot retrieve namespaces", err)
	}

	defer rows.Close()
	for rows.Next() {
		var entry string
		err := rows.Scan(&entry)
		if err != nil {
			log.Error(err)
		}
		namespaces = append(namespaces, entry)
	}

	return namespaces
}

// GetContainers returns all the containers for a given namespace
func (r LogRepository) GetContainers(namespace string) []string {
	containers := make([]string, 0)
	rows, err := r.Connection.Query("SELECT DISTINCT container_name FROM logs WHERE namespace_name = ? ORDER BY container_name", namespace)
	if err != nil {
		log.Error("Cannot retrieve containers", err)
	}

	defer rows.Close()
	for rows.Next() {
		var entry string
		err := rows.Scan(&entry)
		if err != nil {
			log.Error(err)
		}
		containers = append(containers, entry)
	}

	return containers
}

func (r LogRepository) query(query string, args ...interface{}) (*sql.Rows, error) {
	log.Println(query, args)
	return r.Connection.Query(query, args...)
}

// GetLogs returns logs
func (r LogRepository) GetLogs(namespace string, container string, start string, length string, orderColumn string, orderDir string, match string) (int, []LogRecord) {
	logs := make([]LogRecord, 0)
	selectQuery := "SELECT namespace_name, container_name, pod_name, container_image, timestamp, message FROM logs "

	where := "WHERE namespace_name = ? AND container_name = ? "
	if match != "" {
		where += "AND id IN (SELECT id FROM logs_fts WHERE logs_fts MATCH '" + match + "' )"
	}

	selectQuery += where
	if orderDir == "asc" {
		selectQuery += "ORDER BY " + orderColumn + " ASC "
	} else {
		selectQuery += "ORDER BY " + orderColumn + " DESC "
	}
	selectQuery += "LIMIT ? OFFSET ?;"
	rows, err := r.query(selectQuery, namespace, container, length, start)
	if err != nil {
		log.Error("Cannot get log record ", err)
		return 0, logs
	}

	defer rows.Close()
	for rows.Next() {
		var entry LogRecord
		err := rows.Scan(&entry.Namespace, &entry.Container, &entry.Pod,
			&entry.Image, &entry.Timestamp, &entry.Message)
		if err != nil {
			log.Error(err)
		}
		logs = append(logs, entry)
	}
	row := r.Connection.QueryRow("SELECT count(*) FROM logs "+where, namespace, container)
	var count int
	row.Scan(&count)
	return count, logs
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
