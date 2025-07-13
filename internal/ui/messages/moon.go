package messages

import (
	"wms/internal/ui/components"
)

// MoonDataMsg is sent when moon phase data is fetched.
type MoonDataMsg struct {
	Data  *components.MoonResponse
	Error error
}
