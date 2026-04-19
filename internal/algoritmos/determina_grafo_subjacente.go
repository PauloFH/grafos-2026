package algoritmos

import (
	"github.com/PauloFH/grafos-2026/internal/grafo"
)

// Crie um grafo subjacente a partir de um digrafo
func DeterminaGrafoSubjacente(g *grafo.Grafo) *grafo.Grafo {
	subjacente := grafo.NovoGrafo(false, g.NomeArquivo) //Cria um novo grafo não direcionado com o mesmo nome
	arestasAdicionadas := make(map[[2]string]struct{})

	for _, i := range g.Vertices { //Itera sobre cada vértice na ordem de leitura
		for _, j := range g.ListaAdj[i] { //Para cada vizinho, adiciona uma aresta no novo grafo subjacente
			u, v := i, j
			if u > v {
				u, v = v, u
			}

			aresta := [2]string{u, v}
			if _, existe := arestasAdicionadas[aresta]; existe {
				continue
			}

			arestasAdicionadas[aresta] = struct{}{}
			subjacente.AdicionarAresta(i, j) //Adiciona a aresta (i, j) apenas uma vez por par não direcionado
		}
	}
	return subjacente
}
