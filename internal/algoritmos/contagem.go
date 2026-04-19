package algoritmos

import "github.com/PauloFH/grafos-2026/internal/grafo"

// TotalVertices retorna o número de vértices do grafo.
func TotalVertices(g *grafo.Grafo) int {
	return len(g.Vertices)
}

// TotalArestas retorna o número de arestas do grafo.
func TotalArestas(g *grafo.Grafo) int {
	total := 0
	for _, vizinhos := range g.ListaAdj {
		total += len(vizinhos)
	}
	if !g.Direcionado {
		total /= 2
	}
	return total
}
