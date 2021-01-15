package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

// BulkHandler Handler for bulk request
func BulkHandler(c *gin.Context) {
	d := json.NewDecoder(c.Request.Body)

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
			db := c.MustGet("db").(*sql.DB)
			kube := v["kubernetes"].(map[string]interface{})
			stmt, err := db.Prepare(`INSERT INTO logs (namespace_name, container_name, pod_name, container_image, "timestamp", message)
			VALUES(?, ?, ?, ?, ?, ?);
			`)
			if err != nil {
				log.Fatal(err)
			}

			_, err = stmt.Exec(kube["namespace_name"],
				kube["container_name"],
				kube["pod_name"],
				kube["container_image"],
				v["@timestamp"],
				v["log"])

			if err != nil {
				log.Println("Not written ", err)
			}
		}

	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
