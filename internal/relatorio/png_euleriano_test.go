package relatorio

import (
	"strings"
	"testing"

	"github.com/PauloFH/grafos-2026/internal/grafo"
)

// TestGerarDOTTrilhaEuler_Grafo verifica a numeracao das arestas e o destaque
// dos vertices inicial/final para um grafo nao-direcionado.
func TestGerarDOTTrilhaEuler_Grafo(t *testing.T) {
	g := grafo.NovoGrafo(false, "T")
	g.AdicionarAresta("a", "b")
	g.AdicionarAresta("b", "c")
	g.AdicionarAresta("c", "a")

	// Circuito euleriano: a -> b -> c -> a
	dot := gerarDOTTrilhaEuler(g, []string{"a", "b", "c", "a"}, "TRILHA")

	casos := []string{
		"graph TRILHA {",
		`"a" [fillcolor=lightgreen, style=filled];`, // inicio (e fim) destacado
		`"a" -- "b" [label="1"`,
		`"b" -- "c" [label="2"`,
		`"c" -- "a" [label="3"`,
	}
	for _, c := range casos {
		if !strings.Contains(dot, c) {
			t.Errorf("DOT nao contem %q\nDOT:\n%s", c, dot)
		}
	}
}

// TestGerarDOTTrilhaEuler_Digrafo verifica o operador direcionado (->) e o
// destaque de inicio (verde) e fim (azul) num caminho.
func TestGerarDOTTrilhaEuler_Digrafo(t *testing.T) {
	g := grafo.NovoGrafo(true, "TD")
	g.AdicionarAresta("1", "2")
	g.AdicionarAresta("2", "3")

	dot := gerarDOTTrilhaEuler(g, []string{"1", "2", "3"}, "TRILHA_D")

	casos := []string{
		"digraph TRILHA_D {",
		`"1" [fillcolor=lightgreen, style=filled];`, // inicio
		`"3" [fillcolor=lightblue, style=filled];`,  // fim
		`"1" -> "2" [label="1"`,
		`"2" -> "3" [label="2"`,
	}
	for _, c := range casos {
		if !strings.Contains(dot, c) {
			t.Errorf("DOT nao contem %q\nDOT:\n%s", c, dot)
		}
	}
}
