package relatorio

import (
	"fmt"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/grafo"
)

// GerarPNGTrilhaEuler gera um PNG destacando a ordem em que as arestas sao
// percorridas na trilha euleriana (`sequencia` = vertices na ordem da trilha).
// O vertice inicial fica em verde, o final em azul; cada aresta percorrida
// recebe um rotulo com a sua posicao na trilha (1, 2, 3, ...).
func GerarPNGTrilhaEuler(g *grafo.Grafo, sequencia []string, nome, caminho string) error {
	return executarDot(gerarDOTTrilhaEuler(g, sequencia, nome), nome, caminho)
}

func gerarDOTTrilhaEuler(g *grafo.Grafo, sequencia []string, nome string) string {
	gType, op := dotOperador(g.Direcionado)

	var inicio, fim string
	if len(sequencia) > 0 {
		inicio = sequencia[0]
		fim = sequencia[len(sequencia)-1]
	}

	var sb strings.Builder
	sb.WriteString(dotHeader(gType, nome))

	// Vertices: destaca inicio (verde) e fim (azul) da trilha.
	for _, v := range g.Vertices {
		switch v {
		case inicio:
			fmt.Fprintf(&sb, "  \"%s\" [fillcolor=lightgreen, style=filled];\n", v)
		case fim:
			fmt.Fprintf(&sb, "  \"%s\" [fillcolor=lightblue, style=filled];\n", v)
		default:
			fmt.Fprintf(&sb, "  \"%s\";\n", v)
		}
	}

	// Arestas da trilha, numeradas na ordem de percurso.
	for i := 0; i+1 < len(sequencia); i++ {
		u := sequencia[i]
		w := sequencia[i+1]
		fmt.Fprintf(&sb,
			"  \"%s\" %s \"%s\" [label=\"%d\", color=red, fontcolor=red, penwidth=2.0];\n",
			u, op, w, i+1)
	}

	sb.WriteString("}\n")
	return sb.String()
}
