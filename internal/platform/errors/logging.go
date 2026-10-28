package errors

import (
	"context"
	"log"
)

func (r *applicationError) Log(ctx context.Context) {
	logFields := map[string]interface{}{
		"internal_code":    r.InternalCode(),
		"status_code":      r.StatusCode(),
		"message":          r.Message(),
		"original_message": r.OriginalMessage(),
	}

	for k, v := range r.extraFields {
		logFields[k] = v
	}

	log.Printf("[ERROR] %+v", logFields)
}
