package algoritmos

import (
	"fmt"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/grafo"
)

// Guarda os dados brutos da busca
type DFSResult struct {
	Discovery map[string]int
	Finish    map[string]int
	Edges     []string
}

// DFSDigrafo realiza DFS para dígrafos e retorna tempos e tipos de arestas
func DFSDigrafo(g *grafo.Grafo) *DFSResult {
	res := &DFSResult{
		Discovery: make(map[string]int),
		Finish:    make(map[string]int),
		Edges:     make([]string, 0),
	}

	cor := make(map[string]int) // 0:Branco; 1:Cinza; 2:Preto
	tempo := 0

	// Função recursiva usada para catalogar a visita aos vértices
	var visitar func(u string)
	visitar = func(u string) {
		tempo++
		res.Discovery[u] = tempo
		cor[u] = 1 // Cinza

		for _, v := range g.GetVizinhos(u) {
			switch cor[v] {
			case 0: // Branco
				res.Edges = append(res.Edges, fmt.Sprintf("(%s, %s): Árvore", u, v))
				visitar(v)
			case 1: // Cinza
				res.Edges = append(res.Edges, fmt.Sprintf("(%s, %s): Retorno", u, v))
			case 2: // Preto
				if res.Discovery[u] < res.Discovery[v] {
					res.Edges = append(res.Edges, fmt.Sprintf("(%s, %s): Avanço", u, v))
				} else {
					res.Edges = append(res.Edges, fmt.Sprintf("(%s, %s): Cruzamento", u, v))
				}
			}
		}

		cor[u] = 2 // Preto
		tempo++
		res.Finish[u] = tempo
	}

	// Processa os vértices na ordem original
	for _, v := range g.Vertices {
		if cor[v] == 0 {
			visitar(v)
		}
	}

	return res
}

// FormatarDFS converte o resultado em strings amigáveis para o relatório
func FormatarDFS(g *grafo.Grafo, res *DFSResult) (tempos string, arestas string) {
	var sbTempos strings.Builder
	var sbArestas strings.Builder

	for _, v := range g.Vertices {
		sbTempos.WriteString(fmt.Sprintf("  %s: [%d, %d]\n", v, res.Discovery[v], res.Finish[v]))
	}

	for _, e := range res.Edges {
		sbArestas.WriteString(fmt.Sprintf("  %s\n", e))
	}

	return sbTempos.String(), sbArestas.String()
}
