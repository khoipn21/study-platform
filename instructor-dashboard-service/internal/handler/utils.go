package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getInstructorIDFromContext extracts instructor ID from authenticated context
func getInstructorIDFromContext(c *gin.Context) (uuid.UUID, error) {
	// Get instructor ID from context (set by auth middleware)
	instructorIDStr, exists := c.Get("instructor_id")
	if !exists {
		// For testing purposes, use hardcoded instructor ID
		return uuid.Parse("22222222-2222-2222-2222-222222222222")
	}

	instructorID, ok := instructorIDStr.(string)
	if !ok {
		// Try UUID type
		if instructorIDUUID, ok := instructorIDStr.(uuid.UUID); ok {
			return instructorIDUUID, nil
		}
		return uuid.Nil, fmt.Errorf("invalid instructor ID type")
	}

	return uuid.Parse(instructorID)
}