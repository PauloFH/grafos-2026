package relatorio

import (
	"fmt"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/euleriano"
)

// FormataEulerianidade gera o texto da secao EULERIANIDADE a partir do
// resultado da validacao.
func FormataEulerianidade(res euleriano.ResultadoEuler) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "  Classificacao: %s\n", res.Classe.String())
	fmt.Fprintf(&sb, "  Conexo: %t\n", res.Conexo)

	if len(res.GrausImpares) > 0 {
		fmt.Fprintf(&sb, "  Vertices de grau impar: %s\n", strings.Join(res.GrausImpares, ", "))
	}
	if len(res.Desbalanceados) > 0 {
		sb.WriteString("  Vertices desbalanceados (saida-entrada):\n")
		for v, d := range res.Desbalanceados {
			fmt.Fprintf(&sb, "    %s: %+d\n", v, d)
		}
	}
	if res.VerticeInicial != "" {
		fmt.Fprintf(&sb, "  Vertice inicial sugerido: %s\n", res.VerticeInicial)
	}

	return sb.String()
}

// FormataTrilhaEuler gera o texto da secao TRILHA_EULERIANA.
func FormataTrilhaEuler(t euleriano.TrilhaEuler) string {
	if len(t.Sequencia) == 0 {
		return "  (trilha vazia)\n"
	}

	tipo := "Caminho"
	if t.Circuito {
		tipo = "Circuito"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "  Tipo: %s euleriano\n", tipo)
	fmt.Fprintf(&sb, "  Arestas percorridas: %d\n", len(t.Sequencia)-1)
	fmt.Fprintf(&sb, "  Sequencia: %s\n", strings.Join(t.Sequencia, " -> "))
	return sb.String()
}
