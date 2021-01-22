package api

import (
	"time"

	log "github.com/sirupsen/logrus"
	"mogui.it/k8l/go/persistence"
)

func mapModelToDTO(logsRecord []persistence.LogRecord) []LogEntry {
	mapped := make([]LogEntry, 0)
	for _, val := range logsRecord {
		t, err := time.Parse("2006-01-02T15:04:05.000000000+00:00", val.Timestamp)
		if err != nil {
			log.Error(err)
			continue
		}
		e := LogEntry{
			Namespace: val.Namespace,
			Container: val.Container,
			Pod:       val.Pod,
			Message:   val.Message,
			Timestamp: t,
		}
		mapped = append(mapped, e)
	}
	return mapped
}
