package relatorio

import (
	"github.com/PauloFH/grafos-2026/internal/grafo"
)

// GerarPNGTrilhaEuler gera um PNG destacando a ordem em que as arestas sao
// percorridas na trilha euleriana (`sequencia` = vertices na ordem da trilha).
//
// TODO: gerar o DOT numerando/colorindo cada aresta da trilha.
//   - Reusar executarDot, dotHeader e dotOperador de png_algoritmos.go.
//   - Para cada par consecutivo (sequencia[i], sequencia[i+1]) emitir a aresta
//     com label "i+1" (a ordem de percurso).
//
// Stub atual: nao gera nada (retorna nil) para o esqueleto compilar.
func GerarPNGTrilhaEuler(g *grafo.Grafo, sequencia []string, nome, caminho string) error {
	_ = g
	_ = sequencia
	_ = nome
	_ = caminho
	return nil
}
