package model

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/drewharris/shulker/styles"
)

func (m MainModel) View() string {
	if m.isLoading {
		errors := styles.Error.Render(strings.Join(m.errorMessages, "\n"))
		loadingStyle := lipgloss.NewStyle().Padding(2).Bold(true).Italic(true)
		loading := loadingStyle.Height(m.height / 2).Width(m.width / 3).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(m.loadingModel.spinner.View() + "   Starting Server...")

		var nonAlphanumericRegex = regexp.MustCompile(`[^\x20-\x7e]`)
		for i, str := range m.loadingModel.loadingOutput {
			m.loadingModel.loadingOutput[i] = nonAlphanumericRegex.ReplaceAllString(str, "")
		}
		output := styles.Dimmed.Render(strings.Join(m.loadingModel.loadingOutput, "\n"))
		screen := lipgloss.NewStyle().Margin(1).MaxWidth(m.width - 8).Height(m.height).AlignVertical(lipgloss.Center)
		return screen.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				errors,
				lipgloss.JoinHorizontal(lipgloss.Top, loading, output)),
		)
	} else if m.isShuttingDown {
		return lipgloss.NewStyle().Width(m.width).Height(m.height).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(m.loadingModel.spinner.View() + " Shutting down... Please be patient")
	} else if m.isViewingBuildLogs {
		return strings.Join(m.buildMessages, "\n")
	}

	// Main interface
	doc := strings.Builder{}
	var remainingHeight int = m.height

	// Status bar
	statusStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color("#555"))
	var statusBar string
	if m.isBuilding {
		statusBar = statusStyle.Copy().Background(lipgloss.Color("#ab009d")).Render(m.loadingModel.spinner.View() + "  BUILDING... ")
	} else {
		statusBar = statusStyle.Render("IDLE")
	}
	remainingHeight = remainingHeight - lipgloss.Height(statusBar)

	leftMenuContainer := lipgloss.NewStyle().Padding(1).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#fff")).Height(9)
	menu := leftMenuContainer.Render("Shulker Menu\n* Rebuild\n* Restart")

	logContainer := lipgloss.NewStyle().Padding(1).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#fff")).Width(m.width - lipgloss.Width(menu) - 2).AlignVertical(lipgloss.Top)
	serverLogs := logContainer.Render(lastLines(m.serverMessages, m.height-1))
	// serverLogs := logContainer.Render("Helllo there")
	middle := lipgloss.JoinHorizontal(lipgloss.Bottom, menu, serverLogs)

	doc.WriteString(middle + "\n")
	doc.WriteString(statusBar)

	return doc.String()
}
