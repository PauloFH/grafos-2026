package algoritmos

import (
	"sort"

	"github.com/PauloFH/grafos-2026/internal/grafo"
)

type aresta struct{ u, v string }

type ResultadoBiconectividade struct {
	Articulacoes []string
	Blocos       [][]string
	Lowpt        map[string]int
	Num          map[string]int
}

func Biconectividade(g *grafo.Grafo) ResultadoBiconectividade {
	num := make(map[string]int)
	lowpt := make(map[string]int)
	articulacoes := make(map[string]bool)
	blocos := [][]string{}
	pilha := []aresta{}
	cnt := 0

	var dfs func(u, pai string)
	dfs = func(u, pai string) {
		num[u] = cnt
		lowpt[u] = cnt
		cnt++
		filhos := 0

		for _, w := range g.GetVizinhos(u) {
			if _, visitado := num[w]; !visitado {
				pilha = append(pilha, aresta{u, w})
				filhos++
				dfs(w, u)

				if lowpt[w] < lowpt[u] {
					lowpt[u] = lowpt[w]
				}

				if lowpt[w] >= num[u] {
					blocos = append(blocos, extrairBloco(&pilha, u, w))
					if pai != "" {
						articulacoes[u] = true
					}
				}
			} else if w != pai && num[w] < num[u] {
				pilha = append(pilha, aresta{u, w})
				if num[w] < lowpt[u] {
					lowpt[u] = num[w]
				}
			}
		}

		if pai == "" && filhos > 1 {
			articulacoes[u] = true
		}
	}

	for _, v := range g.Vertices {
		if _, ok := num[v]; !ok {
			dfs(v, "")
		}
	}

	arts := make([]string, 0, len(articulacoes))
	for v := range articulacoes {
		arts = append(arts, v)
	}
	sort.Strings(arts)

	return ResultadoBiconectividade{
		Articulacoes: arts,
		Blocos:       blocos,
		Lowpt:        lowpt,
		Num:          num,
	}
}

func extrairBloco(pilha *[]aresta, u, w string) []string {
	vertices := map[string]bool{}
	for len(*pilha) > 0 {
		idx := len(*pilha) - 1
		a := (*pilha)[idx]
		*pilha = (*pilha)[:idx]
		vertices[a.u] = true
		vertices[a.v] = true
		if a.u == u && a.v == w {
			break
		}
	}
	result := make([]string, 0, len(vertices))
	for v := range vertices {
		result = append(result, v)
	}
	sort.Strings(result)
	return result
}
