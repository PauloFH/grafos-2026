package algoritmos

import "github.com/PauloFH/grafos-2026/internal/grafo"

type ResultadoDFS struct {
	Visitados   []string
	Predecessor map[string]string
	Entrada     map[string]int
	Saida       map[string]int
}

func DFS(g *grafo.Grafo, inicio string) ResultadoDFS {
	visitado := make(map[string]bool)
	predecessor := make(map[string]string)
	entrada := make(map[string]int)
	saida := make(map[string]int)
	visitados := []string{}
	tempo := 0

	var visitar func(u string)
	visitar = func(u string) {
		visitado[u] = true
		visitados = append(visitados, u)
		entrada[u] = tempo
		tempo++
		for _, w := range g.GetVizinhos(u) {
			if !visitado[w] {
				predecessor[w] = u
				visitar(w)
			}
		}
		saida[u] = tempo
		tempo++
	}

	visitar(inicio)

	return ResultadoDFS{
		Visitados:   visitados,
		Predecessor: predecessor,
		Entrada:     entrada,
		Saida:       saida,
	}
}
