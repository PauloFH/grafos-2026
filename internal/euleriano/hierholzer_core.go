package euleriano

import (
	"fmt"
	"slices"

	"github.com/PauloFH/grafos-2026/internal/grafo"
)

func executarHierholzer(g *grafo.Grafo, inicio string) ([]string, error) {
	if inicio == "" {
		return nil, nil
	}

	atual := inicio
	pilha := make([]string, 0)
	reversa := make([]string, 0)

	for len(pilha) > 0 || len(g.ListaAdj[atual]) > 0 {
		if len(g.ListaAdj[atual]) > 0 {
			proximo := g.ListaAdj[atual][0]
			pilha = append(pilha, atual)
			if !removeUmaAresta(g, atual, proximo) {
				return nil, fmt.Errorf("remover aresta %s -> %s: %w", atual, proximo, errSemTrilhaEuler)
			}
			atual = proximo
			continue
		}

		reversa = append(reversa, atual)
		atual = pilha[len(pilha)-1]
		pilha = pilha[:len(pilha)-1]
	}

	reversa = append(reversa, atual)
	inverte(reversa)
	return reversa, nil
}

func escolheInicio(g *grafo.Grafo, solicitado, sugerido string) string {
	if solicitado != "" {
		if _, existe := g.ListaAdj[solicitado]; existe {
			return solicitado
		}
		return ""
	}
	if sugerido != "" {
		return sugerido
	}
	if len(g.Vertices) == 0 {
		return ""
	}
	return g.Vertices[0]
}

func removeUmaAresta(g *grafo.Grafo, origem, destino string) bool {
	var removeuOrigem bool
	g.ListaAdj[origem], removeuOrigem = removePrimeiraOcorrencia(g.ListaAdj[origem], destino)
	if !removeuOrigem {
		return false
	}
	if g.Direcionado {
		return true
	}

	var removeuDestino bool
	g.ListaAdj[destino], removeuDestino = removePrimeiraOcorrencia(g.ListaAdj[destino], origem)
	return removeuDestino
}

func removePrimeiraOcorrencia(vizinhos []string, alvo string) ([]string, bool) {
	for i, v := range vizinhos {
		if v != alvo {
			continue
		}
		novos := make([]string, 0, len(vizinhos)-1)
		novos = append(novos, vizinhos[:i]...)
		novos = append(novos, vizinhos[i+1:]...)
		return novos, true
	}
	return vizinhos, false
}

func inverte(valores []string) {
	for i, j := 0, len(valores)-1; i < j; i, j = i+1, j-1 {
		valores[i], valores[j] = valores[j], valores[i]
	}
}

func verificaCircuito(sequencia []string) bool {
	if len(sequencia) == 0 {
		return true
	}
	return sequencia[0] == sequencia[len(sequencia)-1]
}

func arestasNaSequencia(sequencia []string) int {
	if len(sequencia) == 0 {
		return 0
	}
	return len(sequencia) - 1
}

func inicioValidoGrafo(g *grafo.Grafo, inicio string, res ResultadoEuler) bool {
	if inicio == "" {
		return len(g.Vertices) == 0
	}
	if _, existe := g.ListaAdj[inicio]; !existe {
		return false
	}
	if res.Classe == CircuitoEuleriano {
		return g.GrauVertice(inicio) > 0 || arestasDoGrafo(g) == 0
	}

	return slices.Contains(res.GrausImpares, inicio)
}

func inicioValidoDigrafo(g *grafo.Grafo, inicio string, res ResultadoEuler) bool {
	if inicio == "" {
		return len(g.Vertices) == 0
	}
	if _, existe := g.ListaAdj[inicio]; !existe {
		return false
	}
	if res.Classe == CircuitoEuleriano {
		return len(g.ListaAdj[inicio]) > 0 || arestasDoGrafo(g) == 0
	}
	return inicio == res.VerticeInicial
}

func arestasDoGrafo(g *grafo.Grafo) int {
	total := 0
	for _, vizinhos := range g.ListaAdj {
		total += len(vizinhos)
	}
	if !g.Direcionado {
		total /= 2
	}
	return total
}
