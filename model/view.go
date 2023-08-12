package model

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/drewharris/shulker/styles"
)

func (m MainModel) View() string {
	if m.viewMode == startupView {
		errors := styles.Error.Render(strings.Join(m.errorMessages, "\n"))
		loadingStyle := lipgloss.NewStyle().Padding(2).Bold(true).Italic(true)
		loading := loadingStyle.Height(m.height / 2).Width(m.width / 3).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(m.spinner.View() + "   Starting Server...")

		var nonAlphanumericRegex = regexp.MustCompile(`[^\x20-\x7e]`)
		for i, str := range m.loadingOutput {
			m.loadingOutput[i] = nonAlphanumericRegex.ReplaceAllString(str, "")
		}
		output := styles.Dimmed.Render(strings.Join(m.loadingOutput, "\n"))
		screen := lipgloss.NewStyle().Margin(1).MaxWidth(m.width - 8).Height(m.height).AlignVertical(lipgloss.Center)
		return screen.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				errors,
				lipgloss.JoinHorizontal(lipgloss.Top, loading, output)),
		)
	} else if m.viewMode == shutdownView {
		return lipgloss.NewStyle().Width(m.width).Height(m.height).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(m.spinner.View() + " Shutting down... Please be patient")
	} else if m.viewMode == testView {
		return lipgloss.NewStyle().Width(5).MaxHeight(2).Render("He\nllo there indiviudal")
	}

	// Main interface
	doc := strings.Builder{}

	// Status bar
	statusStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color("#555"))
	var statusBar string
	if m.isBuilding {
		statusBar = statusStyle.Copy().Background(lipgloss.Color("#ab009d")).Render("  " + m.spinner.View() + "  BUILDING... ")
	} else {
		statusBar = statusStyle.Render("  IDLE")
	}

	menu := styles.LeftMenuContainer.Copy().Height(m.height - 3).Render("Shulker Menu\n* Rebuild\n* Restart")

	remainingWidth := m.width - lipgloss.Width(menu)

	var serverLogStrings []string
	if m.viewMode == buildView {
		for _, line := range lastLines(m.buildMessages, m.height-2) {
			serverLogStrings = append(serverLogStrings, styles.InlineLog.Copy().MaxWidth(remainingWidth-20).Render(line))
		}
	} else {
		for _, line := range lastLines(m.serverMessages, m.height-2) {
			serverLogStrings = append(serverLogStrings, styles.InlineLog.Copy().MaxWidth(remainingWidth-20).Render(line))
		}
	}

	serverLogs := styles.LogContainer.Copy().MaxWidth(remainingWidth - 1).Width(remainingWidth - 3).MaxHeight(m.height - 10).Render(lipgloss.JoinVertical(lipgloss.Left, serverLogStrings...))

	middle := lipgloss.JoinHorizontal(lipgloss.Bottom, menu, serverLogs)

	doc.WriteString(middle + "\n")
	doc.WriteString(statusBar)

	return doc.String()
}
