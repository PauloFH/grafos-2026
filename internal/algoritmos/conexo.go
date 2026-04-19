package algoritmos

import "github.com/PauloFH/grafos-2026/internal/grafo"

// bfsConexo percorre o grafo em largura a partir de um vértice inicial.
func bfsConexo(adj map[string][]string, inicio string) map[string]bool {
	visitados := make(map[string]bool, len(adj))
	fila := make([]string, 1, len(adj))
	fila[0] = inicio
	visitados[inicio] = true

	for i := 0; i < len(fila); i++ {
		v := fila[i]
		for _, viz := range adj[v] {
			if !visitados[viz] {
				visitados[viz] = true
				fila = append(fila, viz)
			}
		}
	}

	return visitados
}

// EhConexo verifica se o grafo é conexo.
// Para grafos não-direcionados: conexo se todos os vértices são alcançáveis a partir de qualquer um.
// Para dígrafos: verifica conectividade fraca.
func EhConexo(g *grafo.Grafo) bool {
	if len(g.Vertices) == 0 {
		return true
	}

	adj := make(map[string][]string, len(g.Vertices))
	for _, v := range g.Vertices {
		adj[v] = append(adj[v], g.ListaAdj[v]...)
	}
	if g.Direcionado {
		for _, v := range g.Vertices {
			for _, viz := range g.ListaAdj[v] {
				adj[viz] = append(adj[viz], v)
			}
		}
	}

	visitados := bfsConexo(adj, g.Vertices[0])
	return len(visitados) == len(g.Vertices)
}
