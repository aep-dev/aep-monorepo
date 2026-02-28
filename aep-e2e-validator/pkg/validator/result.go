package validator

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type TestStatus string

const (
	StatusPass  TestStatus = "PASSED"
	StatusFail  TestStatus = "FAILED"
	StatusSkip  TestStatus = "SKIPPED"
	StatusError TestStatus = "ERROR"
)

type TestResult struct {
	Name     string
	Status   TestStatus
	Detail   string
	Duration time.Duration
}

const summaryWidth = 60

var renderer *lipgloss.Renderer

func getRenderer() *lipgloss.Renderer {
	if renderer == nil {
		renderer = lipgloss.NewRenderer(os.Stdout, termenv.WithProfile(termenv.ColorProfile()))
	}
	return renderer
}

func greenStyle() lipgloss.Style {
	return getRenderer().NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
}

func redStyle() lipgloss.Style {
	return getRenderer().NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
}

func yellowStyle() lipgloss.Style {
	return getRenderer().NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
}

func printSummary(results []TestResult, totalDuration time.Duration) {
	fmt.Println()
	fmt.Println(centerLine("test results", '='))

	maxNameLen := 0
	for _, r := range results {
		if len(r.Name) > maxNameLen {
			maxNameLen = len(r.Name)
		}
	}

	for _, r := range results {
		icon := statusIcon(r.Status)
		fmt.Printf(" %s  %-*s  %s\n", icon, maxNameLen, r.Name, r.Duration.Round(time.Millisecond))
		if r.Detail != "" {
			for _, line := range strings.Split(r.Detail, "\n") {
				fmt.Printf("         %s\n", line)
			}
		}
	}

	passed, failed, skipped, errored := countByStatus(results)
	var parts []string
	if passed > 0 {
		parts = append(parts, greenStyle().Render(fmt.Sprintf("%d passed", passed)))
	}
	if failed > 0 {
		parts = append(parts, redStyle().Render(fmt.Sprintf("%d failed", failed)))
	}
	if skipped > 0 {
		parts = append(parts, yellowStyle().Render(fmt.Sprintf("%d skipped", skipped)))
	}
	if errored > 0 {
		parts = append(parts, redStyle().Render(fmt.Sprintf("%d error", errored)))
	}
	summary := strings.Join(parts, ", ")
	summary = fmt.Sprintf("%s in %s", summary, totalDuration.Round(time.Millisecond))
	fmt.Println()
	fmt.Println(centerLine(summary, '='))
}

func statusIcon(s TestStatus) string {
	switch s {
	case StatusPass:
		return greenStyle().Render("PASSED ")
	case StatusFail:
		return redStyle().Render("FAILED ")
	case StatusSkip:
		return yellowStyle().Render("SKIPPED")
	case StatusError:
		return redStyle().Render("ERROR  ")
	default:
		return string(s)
	}
}

func countByStatus(results []TestResult) (passed, failed, skipped, errored int) {
	for _, r := range results {
		switch r.Status {
		case StatusPass:
			passed++
		case StatusFail:
			failed++
		case StatusSkip:
			skipped++
		case StatusError:
			errored++
		}
	}
	return
}

func centerLine(text string, pad byte) string {
	plainLen := lipgloss.Width(text)
	if plainLen+4 >= summaryWidth {
		return fmt.Sprintf("%c %s %c", pad, text, pad)
	}
	left := (summaryWidth - plainLen - 2) / 2
	right := summaryWidth - plainLen - 2 - left
	return fmt.Sprintf("%s %s %s", strings.Repeat(string(pad), left), text, strings.Repeat(string(pad), right))
}

func worstExitCode(results []TestResult) int {
	worst := ExitCodeSuccess
	for _, r := range results {
		var code int
		switch r.Status {
		case StatusPass, StatusSkip:
			continue
		case StatusFail:
			code = ExitCodeTestFailed
		case StatusError:
			code = ExitCodeSetupFailed
		}
		if code > worst {
			worst = code
		}
	}
	return worst
}
