package conversoes

import "github.com/PauloFH/grafos-2026/internal/grafo"


func ListaParaMatriz(g *grafo.Grafo) [][]int {
	n := len(g.Vertices)
	idx := make(map[string]int, n)
	for i, v := range g.Vertices {
		idx[v] = i
	}

	m := make([][]int, n)
	for i := range m {
		m[i] = make([]int, n)
	}

	for _, v := range g.Vertices {
		i := idx[v]
		for _, viz := range g.ListaAdj[v] {
			j, ok := idx[viz]
			if !ok {
				continue
			}
			m[i][j] = 1
		}
	}

	return m
}


func MatrizParaLista(g *grafo.Grafo, m [][]int) {
	n := len(g.Vertices)
	if len(m) != n {
		return
	}
	for i := range m {
		if len(m[i]) != n {
			return
		}
	}

	for _, v := range g.Vertices {
		g.ListaAdj[v] = make([]string, 0)
	}

	for i, v := range g.Vertices {
		for j, viz := range g.Vertices {
			if m[i][j] == 1 {
				g.ListaAdj[v] = append(g.ListaAdj[v], viz)
			}
		}
	}
}
