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
		return uuid.Nil, fmt.Errorf("instructor ID not found in context")
	}

	instructorID, ok := instructorIDStr.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid instructor ID type")
	}

	return uuid.Parse(instructorID)
}