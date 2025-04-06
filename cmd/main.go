package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type model struct {
	mesActual string     // Ej: "Enero 2024"
	dias      [][]string // Matriz 6x7 de días
	anio      int        // Año actual (2024)
	mes       time.Month // Mes actual (time.January, etc.)
}

// 2. Función para generar la matriz de días
func generarDiasDelMes(anio int, mes time.Month) [][]string {
	primerDia := time.Date(anio, mes, 1, 0, 0, 0, 0, time.UTC)

	matriz := make([][]string, 6)
	for i := range matriz {
		matriz[i] = make([]string, 7)
		for j := range matriz[i] {
			matriz[i][j] = "  " // Espacios vacíos inicialmente
		}
	}

	fila := 0
	for dia := primerDia; dia.Month() == mes; dia = dia.AddDate(0, 0, 1) {
		columna := int(dia.Weekday())
		if columna == 0 { // Domingo
			columna = 6 // Lo movemos al final (semana empieza en Lunes)
		} else {
			columna--
		}
		matriz[fila][columna] = fmt.Sprintf("%2d", dia.Day())
		if columna == 6 {
			fila++
		}
	}

	return matriz
}

func initialModel() model {
	now := time.Now()
	anio, mes := now.Year(), now.Month()
	return model{
		mesActual: fmt.Sprintf("%s %d", mes.String(), anio),
		dias:      generarDiasDelMes(anio, mes),
		anio:      anio,
		mes:       mes,
	}
}

func (m model) Init() tea.Cmd {
	return nil // No necesitamos comandos iniciales
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			// Mes anterior
			m.mes--
			if m.mes < time.January {
				m.mes = time.December
				m.anio--
			}

		case "right", "l":
			// Mes siguiente
			m.mes++
			if m.mes > time.December {
				m.mes = time.January
				m.anio++
			}

		case "q", "ctrl+c":
			return m, tea.Quit
		}

		// Regenerar días y actualizar título
		m.dias = generarDiasDelMes(m.anio, m.mes)
		m.mesActual = fmt.Sprintf("%s %d", m.mes, m.anio)
	}
	return m, nil
}

var (
	purple = lipgloss.Color("99")
	gray   = lipgloss.Color("245")

	headerStyle = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	cellStyle   = lipgloss.NewStyle().Width(4)
	rowStyle    = cellStyle.Foreground(gray).Align(lipgloss.Center)
)

func (m model) View() string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case table.HeaderRow: // Evaluación directa
				return headerStyle
			default:
				return rowStyle
			}
		}).
		Headers("Lun", "Mar", "Mié", "Jue", "Vie", "Sáb", "Dom").
		Rows(m.dias...)
	headerTitle := lipgloss.NewStyle().
		Foreground(purple).
		Bold(true).
		Padding(0, 1).
		Render(fmt.Sprintf("📅 %s", m.mesActual))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerTitle,
		t.Render(),
		"\n←/→: Cambiar mes | Q: Salir",
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error al ejecutar:", err)
		os.Exit(1)
	}
}
