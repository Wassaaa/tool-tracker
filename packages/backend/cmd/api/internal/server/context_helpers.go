package server

import "github.com/gin-gonic/gin"

// SystemUserID is the fallback actor used when no authenticated user is present.
// This ID is seeded into the database by a migration.
const SystemUserID = "00000000-0000-0000-0000-000000000001"

// GetActorID extracts the actor identifier from the request context.
// Later, auth middleware should set a "currentUser" UUID string in the context.
func GetActorID(c *gin.Context) string {
	if v, ok := c.Get("currentUser"); ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return SystemUserID
}
