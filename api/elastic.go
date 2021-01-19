package api

import (
	"encoding/json"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	p "mogui.it/k8l/persistence"
)

// BulkHandler Handler for bulk request
func BulkHandler(c *gin.Context) {
	d := json.NewDecoder(c.Request.Body)
	var count int = 0
	for {
		// Decode one JSON document.
		var v map[string]interface{}
		err := d.Decode(&v)

		if err != nil {
			// io.EOF is expected at end of stream.
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		// Do something with the value.
		_, isDoc := v["kubernetes"]
		if isDoc {
			repository := c.MustGet("repository").(*p.LogRepository)
			kube := v["kubernetes"].(map[string]interface{})
			record := p.LogRecord{
				Namespace: kube["namespace_name"].(string),
				Container: kube["container_name"].(string),
				Pod:       kube["pod_name"].(string),
				Image:     kube["container_image"].(string),
				Timestamp: v["@timestamp"].(string),
				Message:   v["log"].(string),
			}
			if repository.Save(record) {
				count++
			}
		}

	}
	code := 200
	if count > 0 {
		code = 201
	}
	c.JSON(code, gin.H{
		"indexed": count,
	})
}
