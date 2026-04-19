package algoritmos

import "github.com/PauloFH/grafos-2026/internal/grafo"

// EhConexo verifica se o grafo é conexo.
// Para grafos não-direcionados: conexo se todos os vértices são alcançáveis a partir de qualquer um.
// Para dígrafos: verifica conectividade fraca.
func EhConexo(g *grafo.Grafo) bool {
	if len(g.Vertices) == 0 {
		return true
	}

	alvo := g
	if g.Direcionado {
		alvo = digrafoToGrafo(g)
	}

	res := BFS(alvo, alvo.Vertices[0])
	return len(res.Visitados) == len(g.Vertices)
}

// digrafoToGrafo cria uma cópia não-direcionada do dígrafo para checar conectividade fraca.
func digrafoToGrafo(g *grafo.Grafo) *grafo.Grafo {
	sub := grafo.NovoGrafo(false, "")
	for _, v := range g.Vertices {
		sub.AdicionarVertice(v)
	}
	for _, v := range g.Vertices {
		for _, viz := range g.ListaAdj[v] {
			sub.AdicionarAresta(v, viz)
		}
	}
	return sub
}
