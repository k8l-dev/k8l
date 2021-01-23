package api

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	p "mogui.it/k8l/go/persistence"
)

// BulkHandler Handler for bulk request
func BulkHandler(c *gin.Context) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	log.Debug(string(jsonData))
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

		_, isKubeLog := v["kubernetes"]
		_, isGenericLog := v["container_name"]
		repository := p.STORAGE.LogRepository

		if isKubeLog {
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
		} else if isGenericLog {
			record := p.LogRecord{
				Namespace: "_generic",
				Container: v["container_name"].(string),
				Pod:       "",
				Image:     "",
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
