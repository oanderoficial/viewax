package ui

import (
	"fmt"
	"strings"

	"viewax/internal/system"

	"github.com/charmbracelet/lipgloss"
)

var (
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	dangerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cardStyle   = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
)

func RenderSimple(s *system.Snapshot) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Viewax") + "\n")
	b.WriteString(mutedStyle.Render("Visão simples do sistema") + "\n\n")

	b.WriteString(metricLine("CPU", s.CPUPercent) + "\n")
	b.WriteString(metricLine("Memória", s.MemoryPercent) + "\n")
	b.WriteString(metricLine("Disco", s.DiskPercent) + "\n\n")

	b.WriteString(scoreLine(healthScore(s)) + "\n\n")
	b.WriteString(explain(s) + "\n\n")
	b.WriteString(mutedStyle.Render("Ctrl+C para sair"))

	return cardStyle.Render(b.String())
}

func RenderPro(s *system.Snapshot) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Viewax Pro") + "\n")
	b.WriteString(mutedStyle.Render("Visão completa do sistema") + "\n\n")

	b.WriteString(metricLine("CPU", s.CPUPercent) + "\n")
	b.WriteString(metricLine("Memória", s.MemoryPercent) + "\n")
	b.WriteString(metricLine("Swap", s.SwapPercent) + "\n")
	b.WriteString(metricLine("Disco", s.DiskPercent) + "\n")
	b.WriteString(fmt.Sprintf("%-8s %.2f\n", "Load", s.Load1))
	b.WriteString(fmt.Sprintf("%-8s %s\n", "Uptime", s.UptimeHuman))

	b.WriteString("\n")
	b.WriteString(scoreLine(healthScore(s)) + "\n\n")
	b.WriteString(renderTopProcesses(s.TopCPUProcesses, 5))
	b.WriteString("\n\n")
	b.WriteString(explain(s))

	return cardStyle.Render(b.String())
}

func RenderProcesses(s *system.Snapshot) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Viewax Processos") + "\n")
	b.WriteString(mutedStyle.Render("Processos consumindo mais CPU") + "\n\n")
	b.WriteString(renderTopProcesses(s.TopCPUProcesses, 8))
	b.WriteString("\n\n")
	b.WriteString(mutedStyle.Render("Ctrl+C para sair"))

	return cardStyle.Render(b.String())
}
func RenderMemory(s *system.Snapshot) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Viewax Memória") + "\n")
	b.WriteString(mutedStyle.Render("Processos consumindo mais memória") + "\n\n")
	b.WriteString(metricLine("Memória", s.MemoryPercent) + "\n")
	b.WriteString(metricLine("Swap", s.SwapPercent) + "\n\n")
	b.WriteString(renderTopMemoryProcesses(s.TopMemoryProcesses, 8))
	b.WriteString("\n\n")
	b.WriteString(mutedStyle.Render("Ctrl+C para sair"))

	return cardStyle.Render(b.String())
}

func renderTopProcesses(processes []system.ProcessInfo, limit int) string {
	var b strings.Builder

	b.WriteString("Top processos:\n")

	if len(processes) == 0 {
		b.WriteString(mutedStyle.Render("Nenhum processo encontrado."))
		return b.String()
	}

	for i, p := range processes {
		if i >= limit {
			break
		}

		name := p.Name
		if len(name) > 18 {
			name = name[:18]
		}

		b.WriteString(fmt.Sprintf(
			"%d. %-18s CPU %6.1f%%   RAM %7.0f MB\n",
			i+1,
			name,
			p.CPUPercent,
			p.MemoryMB,
		))
	}

	return strings.TrimRight(b.String(), "\n")
}

func renderTopMemoryProcesses(processes []system.ProcessInfo, limit int) string {
	var b strings.Builder

	b.WriteString("Top processos por memória:\n")

	if len(processes) == 0 {
		b.WriteString(mutedStyle.Render("Nenhum processo encontrado."))
		return b.String()
	}

	for i, p := range processes {
		if i >= limit {
			break
		}

		name := p.Name
		if len(name) > 18 {
			name = name[:18]
		}

		b.WriteString(fmt.Sprintf(
			"%d. %-18s RAM %7.0f MB   %5.1f%%\n",
			i+1,
			name,
			p.MemoryMB,
			p.MemoryUsage,
		))
	}

	return strings.TrimRight(b.String(), "\n")
}
func metricLine(name string, value float64) string {
	return fmt.Sprintf(
		"%-8s %s %5.0f%%  %s",
		name,
		progressBar(value, 18),
		value,
		statusLabel(value),
	)
}

func progressBar(value float64, width int) string {
	filled := int((value / 100) * float64(width))

	if filled > width {
		filled = width
	}

	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	switch {
	case value >= 90:
		return dangerStyle.Render(bar)
	case value >= 75:
		return warnStyle.Render(bar)
	default:
		return okStyle.Render(bar)
	}
}

func statusLabel(value float64) string {
	switch {
	case value >= 90:
		return dangerStyle.Render("Crítico")
	case value >= 75:
		return warnStyle.Render("Atenção")
	default:
		return okStyle.Render("Normal")
	}
}

func healthScore(s *system.Snapshot) int {
	score := 100

	score -= penalty(s.CPUPercent, 70, 95, 25)
	score -= penalty(s.MemoryPercent, 70, 95, 30)
	score -= penalty(s.SwapPercent, 40, 90, 20)
	score -= penalty(s.DiskPercent, 75, 95, 35)

	if score < 0 {
		return 0
	}

	return score
}

func penalty(value, warningAt, criticalAt float64, maxPenalty int) int {
	if value <= warningAt {
		return 0
	}

	if value >= criticalAt {
		return maxPenalty
	}

	ratio := (value - warningAt) / (criticalAt - warningAt)
	return int(ratio * float64(maxPenalty))
}

func scoreLine(score int) string {
	label := "Excelente"
	style := okStyle

	switch {
	case score < 50:
		label = "Ruim"
		style = dangerStyle
	case score < 75:
		label = "Atenção"
		style = warnStyle
	}

	return fmt.Sprintf("Saúde do sistema: %s", style.Render(fmt.Sprintf("%d/100 • %s", score, label)))
}

func explain(s *system.Snapshot) string {
	var messages []string

	if s.CPUPercent >= 90 {
		messages = append(messages, dangerStyle.Render("CPU muito alta. Algum app pode estar exigindo muito processamento."))
	} else if s.CPUPercent >= 75 {
		messages = append(messages, warnStyle.Render("CPU alta. O sistema pode ficar menos responsivo."))
	}

	if s.MemoryPercent >= 90 {
		messages = append(messages, dangerStyle.Render("Memória quase cheia. Apps podem travar ou ficar lentos."))
	} else if s.MemoryPercent >= 75 {
		messages = append(messages, warnStyle.Render("Memória em uso elevado. Fechar apps pode ajudar."))
	}

	if s.SwapPercent >= 70 {
		messages = append(messages, warnStyle.Render("Swap em uso alto. Pode indicar falta de memória RAM."))
	}

	if s.DiskPercent >= 90 {
		messages = append(messages, dangerStyle.Render("Disco quase cheio. Libere espaço o quanto antes."))
	} else if s.DiskPercent >= 75 {
		messages = append(messages, warnStyle.Render("Disco com uso alto. Vale revisar arquivos grandes."))
	}

	if len(messages) == 0 {
		return okStyle.Render("Tudo parece saudável agora.")
	}

	return strings.Join(messages, "\n")
}
