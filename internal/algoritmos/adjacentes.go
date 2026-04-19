package algoritmos

import "github.com/PauloFH/grafos-2026/internal/grafo"

// SaoAdjacentes
func SaoAdjacentes(g *grafo.Grafo, a, b string) bool {
	vizinhos, ok := g.ListaAdj[a]
	if !ok {
		return false
	}
	for _, v := range vizinhos {
		if v == b {
			return true
		}
	}
	return false
}

// ParesAdjacentes retorna todos os pares adjacentes do grafo.
func ParesAdjacentes(g *grafo.Grafo) [][2]string {
	type par = [2]string
	visitados := make(map[par]bool)
	var resultado [][2]string

	for _, v := range g.Vertices {
		for _, viz := range g.ListaAdj[v] {
			var chave par
			if !g.Direcionado && viz < v {
				chave = par{viz, v}
			} else {
				chave = par{v, viz}
			}
			if visitados[chave] {
				continue
			}
			visitados[chave] = true
			resultado = append(resultado, par{v, viz})
		}
	}
	return resultado
}
