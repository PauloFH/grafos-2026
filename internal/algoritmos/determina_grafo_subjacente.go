package algoritmos

import (
	"github.com/PauloFH/grafos-2026/internal/grafo"
)

// Crie um grafo subjacente a partir de um digrafo
// Item 18 OPC - Responsável: João Victor
func DeterminaGrafoSubjacente(g *grafo.Grafo) *grafo.Grafo {
	subjacente := grafo.NovoGrafo(false, g.NomeArquivo) //Cria um novo grafo não direcionado com o mesmo nome
	arestasAdicionadas := make(map[[2]int]struct{})

	for i, vizinhos := range g.ListaAdj { //Itera sobre cada vértice e seus vizinhos
		for _, j := range vizinhos { //Para cada vizinho, adiciona uma aresta no novo grafo subjacente
			u, v := i, j
			if u > v {
				u, v = v, u
			}

			aresta := [2]int{u, v}
			if _, existe := arestasAdicionadas[aresta]; existe {
				continue
			}

			arestasAdicionadas[aresta] = struct{}{}
			subjacente.AdicionarAresta(i, j) //Adiciona a aresta (i, j) apenas uma vez por par não direcionado
		}
	}
	return subjacente
}
