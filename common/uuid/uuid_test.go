package uuid

import (
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	uuid, err := New()
	if err != nil {
		t.Fatalf("New() returned an error: %v", err)
	}

	if len(uuid) != 36 {
		t.Errorf("Expected UUID length to be 36, but got %d", len(uuid))
	}

	// Regex to check UUID v4 format
	// xxxxxxxx-xxxx-4xxx-[89ab]xxx-xxxxxxxxxxxx
	re := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !re.MatchString(uuid) {
		t.Errorf("Generated UUID %s is not a valid V4 UUID", uuid)
	}
}
