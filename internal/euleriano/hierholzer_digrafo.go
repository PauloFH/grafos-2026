package euleriano

import "github.com/PauloFH/grafos-2026/internal/grafo"

// HierholzerDigrafo extrai a cadeia euleriana de um DIGRAFO, partindo de
// `inicio`. Cada arco (origem->destino) deve ser percorrido exatamente uma vez.
//
// TODO: implementar o algoritmo de Hierholzer para digrafos.
//   - Validar antes com Classifica(g); retornar erro se nao for euleriano/semi.
//   - Trabalhar sobre g.Clone(); em digrafo, RemoverAresta ja remove so o arco
//     origem->destino. Consumir apenas arcos de saida.
//   - Conferir: len(Sequencia)-1 == algoritmos.TotalArestas(g).
func HierholzerDigrafo(g *grafo.Grafo, inicio string) (TrilhaEuler, error) {
	return TrilhaEuler{}, nil
}
