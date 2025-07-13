package components

import (
	"time"
)

type Sun struct {
	Sunrise    time.Time
	Sunset     time.Time
	DayLength  time.Duration
	CurrentPos string // "day" or "night"
	Icon       string
}

func NewSun() Sun {
	// Mock data for now
	now := time.Now()
	sunrise := time.Date(now.Year(), now.Month(), now.Day(), 5, 47, 0, 0, now.Location())
	sunset := time.Date(now.Year(), now.Month(), now.Day(), 20, 21, 0, 0, now.Location())

	currentPos := "day"
	icon := "☀️"
	if now.Before(sunrise) || now.After(sunset) {
		currentPos = "night"
		icon = "🌙"
	}

	return Sun{
		Sunrise:    sunrise,
		Sunset:     sunset,
		DayLength:  sunset.Sub(sunrise),
		CurrentPos: currentPos,
		Icon:       icon,
	}
}
