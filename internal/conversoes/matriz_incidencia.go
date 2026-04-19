package conversoes

import "github.com/PauloFH/grafos-2026/internal/grafo"

type par = [2]string

func enumerarArestas(g *grafo.Grafo) [][2]string {
	visitados := make(map[par]bool)
	var arestas [][2]string

	for _, v := range g.Vertices {
		for _, viz := range g.ListaAdj[v] {
			if g.Direcionado {
				arestas = append(arestas, par{v, viz})
				continue
			}
			chave := par{v, viz}
			if viz < v {
				chave = par{viz, v}
			}
			if !visitados[chave] {
				visitados[chave] = true
				arestas = append(arestas, chave)
			}
		}
	}

	return arestas
}

// MatrizIncidencia retorna a matriz de incidência (n×m) e a lista de arestas.
// Não-direcionado: célula = 1 se vértice é extremo da aresta.
// Direcionado: +1 para origem, -1 para destino.
func MatrizIncidencia(g *grafo.Grafo) ([][]int, [][2]string) {
	arestas := enumerarArestas(g)

	n := len(g.Vertices)
	m := len(arestas)

	idx := make(map[string]int, n)
	for i, v := range g.Vertices {
		idx[v] = i
	}

	matriz := make([][]int, n)
	for i := range matriz {
		matriz[i] = make([]int, m)
	}

	for j, aresta := range arestas {
		i := idx[aresta[0]]
		k := idx[aresta[1]]
		if i == k {
			// Laço: convenção — marca 2 no vértice (não-direcionado) ou +1 (direcionado)
			if g.Direcionado {
				matriz[i][j] = +1
			} else {
				matriz[i][j] = 2
			}
			continue
		}
		if g.Direcionado {
			matriz[i][j] = +1
			matriz[k][j] = -1
		} else {
			matriz[i][j] = 1
			matriz[k][j] = 1
		}
	}

	return matriz, arestas
}
