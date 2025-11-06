package ui

import (
	"strings"

	"github.com/fatih/color"
)

type ColorFuncs struct {
	Cyan   func(a ...interface{}) string
	Yellow func(a ...interface{}) string
	Green  func(a ...interface{}) string
	Red    func(a ...interface{}) string
	Blue   func(a ...interface{}) string
	Bold   func(a ...interface{}) string
	Gray   func(a ...interface{}) string
}

func NewColorFuncs() *ColorFuncs {
	return &ColorFuncs{
		Cyan:   color.New(color.FgCyan).SprintFunc(),
		Yellow: color.New(color.FgYellow).SprintFunc(),
		Green:  color.New(color.FgGreen).SprintFunc(),
		Red:    color.New(color.FgRed).SprintFunc(),
		Blue:   color.New(color.FgBlue).SprintFunc(),
		Bold:   color.New(color.Bold).SprintFunc(),
		Gray:   color.New(color.FgHiBlack).SprintFunc(),
	}
}

func GetStatusColor(status string) func(a ...interface{}) string {
	status = strings.ToLower(status)

	if strings.Contains(status, "done") || strings.Contains(status, "closed") ||
		strings.Contains(status, "resolved") {
		return color.New(color.FgGreen).SprintFunc()
	}

	if strings.Contains(status, "progress") || strings.Contains(status, "review") {
		return color.New(color.FgYellow).SprintFunc()
	}

	if strings.Contains(status, "to do") || strings.Contains(status, "open") ||
		strings.Contains(status, "backlog") {
		return color.New(color.FgBlue).SprintFunc()
	}

	return color.New(color.FgWhite).SprintFunc()
}

func GetPriorityColor(priority string) func(a ...interface{}) string {
	priority = strings.ToLower(priority)

	if strings.Contains(priority, "highest") || strings.Contains(priority, "critical") {
		return color.New(color.FgRed, color.Bold).SprintFunc()
	}

	if strings.Contains(priority, "high") {
		return color.New(color.FgRed).SprintFunc()
	}

	if strings.Contains(priority, "medium") {
		return color.New(color.FgYellow).SprintFunc()
	}

	if strings.Contains(priority, "low") || strings.Contains(priority, "lowest") {
		return color.New(color.FgGreen).SprintFunc()
	}

	return color.New(color.FgWhite).SprintFunc()
}
