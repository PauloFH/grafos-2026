package algoritmos

import "github.com/PauloFH/grafos-2026/internal/grafo"

type ResultadoBipartido struct {
	Bipartido bool
	Nivel     map[string]int
}

func Bipartido(g *grafo.Grafo) ResultadoBipartido {
	nivel := make(map[string]int)
	bipartido := true

	var visitar func(u, pai string)
	visitar = func(u, pai string) {
		for _, w := range g.GetVizinhos(u) {
			if _, visitado := nivel[w]; !visitado {
				nivel[w] = nivel[u] + 1
				visitar(w, u)
			} else if w != pai {
				ciclo := nivel[u] - nivel[w] + 1
				if ciclo%2 != 0 {
					bipartido = false
				}
			}
		}
	}

	for _, v := range g.Vertices {
		if _, ok := nivel[v]; !ok {
			nivel[v] = 0
			visitar(v, "")
		}
	}

	return ResultadoBipartido{
		Bipartido: bipartido,
		Nivel:     nivel,
	}
}
