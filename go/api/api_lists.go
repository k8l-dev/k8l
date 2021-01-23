/*
 * k8l
 *
 * Kubernetes light logs API
 *
 * API version: 1.0
 * Contact: mogui83@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mogui.it/k8l/go/persistence"
)

// GetContainers - Get Containers for a given namespace
func GetContainers(c *gin.Context) {
	repository := persistence.STORAGE.LogRepository
	containers := repository.GetContainers(c.Param("namespace"))
	c.JSON(http.StatusOK, containers)
}

// GetNamespaces - Get All Namespaces
func GetNamespaces(c *gin.Context) {
	repository := persistence.STORAGE.LogRepository
	namespaces := repository.GetNamespaces()
	c.JSON(http.StatusOK, namespaces)
}