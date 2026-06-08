package euleriano

import (
	"errors"
	"fmt"

	"github.com/PauloFH/grafos-2026/internal/algoritmos"
	"github.com/PauloFH/grafos-2026/internal/grafo"
)

var (
	errGrafoDirecionado = errors.New("euleriano: esperado grafo nao-direcionado")
	errSemTrilhaEuler   = errors.New("euleriano: grafo nao possui trilha euleriana")
	errInicioInvalido   = errors.New("euleriano: vertice inicial invalido")
)

// HierholzerGrafo extrai a cadeia euleriana de um grafo NAO-direcionado,
// partindo de `inicio`. Cada aresta deve ser percorrida exatamente uma vez.
func HierholzerGrafo(g *grafo.Grafo, inicio string) (TrilhaEuler, error) {
	if g == nil {
		return TrilhaEuler{}, fmt.Errorf("validar grafo: %w", errInicioInvalido)
	}
	if g.Direcionado {
		return TrilhaEuler{}, errGrafoDirecionado
	}

	res := Classifica(g)
	if res.Classe == NaoEuleriano {
		return TrilhaEuler{}, errSemTrilhaEuler
	}

	inicioEscolhido := escolheInicio(g, inicio, res.VerticeInicial)
	if inicioEscolhido == "" && algoritmos.TotalArestas(g) > 0 {
		return TrilhaEuler{}, errInicioInvalido
	}
	if !inicioValidoGrafo(g, inicioEscolhido, res) {
		return TrilhaEuler{}, errInicioInvalido
	}

	sequencia, err := executarHierholzer(g.Clone(), inicioEscolhido)
	if err != nil {
		return TrilhaEuler{}, fmt.Errorf("executar hierholzer em grafo: %w", err)
	}
	if arestasNaSequencia(sequencia) != algoritmos.TotalArestas(g) {
		return TrilhaEuler{}, fmt.Errorf("validar trilha extraida: %w", errSemTrilhaEuler)
	}

	return TrilhaEuler{
		Sequencia: sequencia,
		Circuito:  verificaCircuito(sequencia),
	}, nil
}
