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
		return lipgloss.NewStyle().Width(40).MaxHeight(30).Margin(3).Padding(1).Border(lipgloss.NormalBorder()).Render(lastLines(strings.Join(m.serverMessages, "\n"), 10))
	}

	// Main interface
	doc := strings.Builder{}

	// Status bar
	statusStyle := lipgloss.NewStyle().Padding(1, 1).Width(m.width)
	var statusBar string
	var statusText string
	if m.isBuilding {
		statusText = styles.Purple.Render("  " + m.spinner.View() + "  BUILDING... ")
	} else {
		if m.reloadSpigotOnBuild {
			statusText = "  Reload Spigot on build: " + styles.Highlight.Render("ON") + "  "
		} else {
			statusText = "  Reload Spigot on build: " + "OFF" + " "
		}
	}

	if m.cmdInput.Focused() {
		statusText += " " + m.cmdInput.View()
	}

	statusBar = statusStyle.Render(statusText)

	bottom := lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().PaddingLeft(3).Render(m.help.View(m.keys)), statusBar)

	logsContainer := lipgloss.NewStyle().Height(m.height-lipgloss.Height(bottom)-4).Padding(1, 3).Border(lipgloss.NormalBorder()).Width(m.width - 2)
	// logs := logsContainer.Render(m.viewport.View())
	m.viewport.Height = m.height - lipgloss.Height(bottom) - 6
	m.viewport.Width = m.width - 8
	m.viewport.GotoBottom()
	logs := logsContainer.Render(
		m.viewport.View(),
	)

	doc.WriteString(logs + "\n" + "\n")
	doc.WriteString(bottom)

	return doc.String()
}
