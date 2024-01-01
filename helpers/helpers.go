package helpers

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

func UUID(prefix ...string) string {
	u, err := uuid.NewUUID()
	if err != nil {
		return time.Now().Format(time.StampNano)
	}
	// Convert the UUID to a string and remove hyphens
	uuidString := strings.ReplaceAll(u.String(), "-", "")
	shortUUID := uuidString[:16]
	return shortUUID
}
