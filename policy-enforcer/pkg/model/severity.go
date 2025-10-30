package model

import "strings"

type Severity uint8

const (
	NoneSeverity Severity = iota
	LowSeverity
	MediumSeverity
	HighSeverity
	CriticalSeverity

	UnsetSeverity Severity = 255
)

func (s Severity) String() string {
	switch s {
	case LowSeverity:
		return "low"
	case MediumSeverity:
		return "medium"
	case HighSeverity:
		return "high"
	case CriticalSeverity:
		return "critical"
	default:
		return "none"
	}
}

func (s *Severity) Set(str string) {
	switch strings.ToLower(str) {
	case LowSeverity.String():
		*s = 1
	case MediumSeverity.String():
		*s = 2
	case HighSeverity.String():
		*s = 3
	case CriticalSeverity.String():
		*s = 4
	default:
		*s = 0
	}
}
