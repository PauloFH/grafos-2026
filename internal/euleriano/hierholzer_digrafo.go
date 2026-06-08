package euleriano

import (
	"errors"
	"fmt"

	"github.com/PauloFH/grafos-2026/internal/algoritmos"
	"github.com/PauloFH/grafos-2026/internal/grafo"
)

var errDigrafoNaoDirecionado = errors.New("euleriano: esperado digrafo")

// HierholzerDigrafo extrai a cadeia euleriana de um DIGRAFO, partindo de
// `inicio`. Cada arco (origem->destino) deve ser percorrido exatamente uma vez.
func HierholzerDigrafo(g *grafo.Grafo, inicio string) (TrilhaEuler, error) {
	if g == nil {
		return TrilhaEuler{}, fmt.Errorf("validar digrafo: %w", errInicioInvalido)
	}
	if !g.Direcionado {
		return TrilhaEuler{}, errDigrafoNaoDirecionado
	}

	res := Classifica(g)
	if res.Classe == NaoEuleriano {
		return TrilhaEuler{}, errSemTrilhaEuler
	}

	inicioEscolhido := escolheInicio(g, inicio, res.VerticeInicial)
	if inicioEscolhido == "" && algoritmos.TotalArestas(g) > 0 {
		return TrilhaEuler{}, errInicioInvalido
	}
	if !inicioValidoDigrafo(g, inicioEscolhido, res) {
		return TrilhaEuler{}, errInicioInvalido
	}

	sequencia, err := executarHierholzer(g.Clone(), inicioEscolhido)
	if err != nil {
		return TrilhaEuler{}, fmt.Errorf("executar hierholzer em digrafo: %w", err)
	}
	if arestasNaSequencia(sequencia) != algoritmos.TotalArestas(g) {
		return TrilhaEuler{}, fmt.Errorf("validar trilha extraida: %w", errSemTrilhaEuler)
	}

	return TrilhaEuler{
		Sequencia: sequencia,
		Circuito:  verificaCircuito(sequencia),
	}, nil
}
