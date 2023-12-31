package helpers

import (
	"github.com/google/uuid"
	"strings"
)

func UUID(prefix ...string) string {
	// Generate a UUID
	u, err := uuid.NewUUID()
	if err != nil {
		return ""
	}

	// Convert the UUID to a string and remove hyphens
	uuidString := strings.ReplaceAll(u.String(), "-", "")

	// Take the first 'length' characters of the UUID string
	shortUUID := uuidString[:16]

	return shortUUID
}
