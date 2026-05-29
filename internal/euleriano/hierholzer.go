package euleriano

import "github.com/PauloFH/grafos-2026/internal/grafo"

// HierholzerGrafo extrai a cadeia euleriana de um grafo NAO-direcionado,
// partindo de `inicio`. Cada aresta deve ser percorrida exatamente uma vez.
//
// TODO: implementar o algoritmo de Hierholzer.
//   - Validar antes com Classifica(g); retornar erro se nao for euleriano/semi.
//   - Trabalhar sobre g.Clone() consumindo arestas com RemoverAresta.
//   - Pode usar algoritmos.Pilha. Ao travar, faz backtrack adicionando o
//     vertice a trilha (ordem reversa).
//   - Conferir: len(Sequencia)-1 == algoritmos.TotalArestas(g).
func HierholzerGrafo(g *grafo.Grafo, inicio string) (TrilhaEuler, error) {
	return TrilhaEuler{}, nil
}
