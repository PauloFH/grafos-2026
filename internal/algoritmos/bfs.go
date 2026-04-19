package algoritmos

import "github.com/PauloFH/grafos-2026/internal/grafo"

type ResultadoBFS struct {
	Visitados   []string
	Predecessor map[string]string
	Nivel       map[string]int
}

func BFS(g *grafo.Grafo, inicio string) ResultadoBFS {
	visitado := make(map[string]bool)
	predecessor := make(map[string]string)
	nivel := make(map[string]int)
	visitados := []string{}

	fila := &Fila{}
	visitado[inicio] = true
	nivel[inicio] = 0
	fila.Enfileira(inicio)

	for fila.Tamanho() > 0 {
		u, _ := fila.Desenfileira()
		visitados = append(visitados, u)
		for _, w := range g.GetVizinhos(u) {
			if !visitado[w] {
				visitado[w] = true
				predecessor[w] = u
				nivel[w] = nivel[u] + 1
				fila.Enfileira(w)
			}
		}
	}

	return ResultadoBFS{
		Visitados:   visitados,
		Predecessor: predecessor,
		Nivel:       nivel,
	}
}
