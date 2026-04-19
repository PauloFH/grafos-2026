package relatorio

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/algoritmos"
	"github.com/PauloFH/grafos-2026/internal/grafo"
)

func FormataBFS(res algoritmos.ResultadoBFS, inicio string) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Vértice inicial: %s\n", inicio)
	fmt.Fprintf(&sb, "Ordem de visita: %s\n\n", strings.Join(res.Visitados, " -> "))
	fmt.Fprintf(&sb, "%-10s %-8s %s\n", "Vértice", "Nível", "Predecessor")
	sb.WriteString(strings.Repeat("-", 30) + "\n")
	for _, v := range res.Visitados {
		pred := res.Predecessor[v]
		if pred == "" {
			pred = "-"
		}
		fmt.Fprintf(&sb, "%-10s %-8d %s\n", v, res.Nivel[v], pred)
	}

	return sb.String()
}

func FormataDFS(res algoritmos.ResultadoDFS, inicio string) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Vértice inicial: %s\n", inicio)
	fmt.Fprintf(&sb, "Ordem de visita: %s\n\n", strings.Join(res.Visitados, " -> "))
	fmt.Fprintf(&sb, "%-10s %-8s %-8s %s\n", "Vértice", "Entrada", "Saída", "Predecessor")
	sb.WriteString(strings.Repeat("-", 38) + "\n")
	for _, v := range res.Visitados {
		pred := res.Predecessor[v]
		if pred == "" {
			pred = "-"
		}
		fmt.Fprintf(&sb, "%-10s %-8d %-8d %s\n", v, res.Entrada[v], res.Saida[v], pred)
	}

	return sb.String()
}

func FormataBiconectividade(g *grafo.Grafo) string {
	res := algoritmos.Biconectividade(g)
	var sb strings.Builder

	fmt.Fprintf(&sb, "%-10s %-6s %s\n", "Vértice", "num", "lowpt")
	sb.WriteString(strings.Repeat("-", 26) + "\n")
	for _, v := range g.Vertices {
		fmt.Fprintf(&sb, "%-10s %-6d %d\n", v, res.Num[v], res.Lowpt[v])
	}

	sb.WriteString("\n")
	if len(res.Articulacoes) == 0 {
		sb.WriteString("Articulações: nenhuma\n")
	} else {
		fmt.Fprintf(&sb, "Articulações: %s\n", strings.Join(res.Articulacoes, ", "))
	}

	fmt.Fprintf(&sb, "\nBlocos (%d):\n", len(res.Blocos))
	for i, bloco := range res.Blocos {
		fmt.Fprintf(&sb, "  Bloco %d: {%s}\n", i+1, strings.Join(bloco, ", "))
	}

	return sb.String()
}

func FormataBipartido(g *grafo.Grafo) string {
	res := algoritmos.Bipartido(g)
	var sb strings.Builder

	if !res.Bipartido {
		sb.WriteString("É bipartido: NÃO\n")
		return sb.String()
	}

	sb.WriteString("É bipartido: SIM\n\n")

	grupoA := []string{}
	grupoB := []string{}
	for _, v := range g.Vertices {
		if res.Nivel[v]%2 == 0 {
			grupoA = append(grupoA, v)
		} else {
			grupoB = append(grupoB, v)
		}
	}
	sort.Strings(grupoA)
	sort.Strings(grupoB)

	fmt.Fprintf(&sb, "Grupo A (nível par):   {%s}\n", strings.Join(grupoA, ", "))
	fmt.Fprintf(&sb, "Grupo B (nível ímpar): {%s}\n", strings.Join(grupoB, ", "))

	return sb.String()
}
