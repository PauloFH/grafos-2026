// Package euleriano implementa a Unidade 2 (Projeto 3): validacao de
// eulerianidade e extracao de trilhas eulerianas (Hierholzer) para grafos e
// digrafos. Reusa o nucleo compartilhado em internal/{grafo,algoritmos}.
package euleriano

import "github.com/PauloFH/grafos-2026/internal/grafo"

// Classe indica o resultado da analise de eulerianidade de um grafo/digrafo.
type Classe int

const (
	// NaoEuleriano: nao possui caminho nem circuito euleriano.
	NaoEuleriano Classe = iota
	// CaminhoEuleriano (semi-euleriano): possui caminho euleriano, nao circuito.
	CaminhoEuleriano
	// CircuitoEuleriano (euleriano): possui circuito euleriano.
	CircuitoEuleriano
)

// String descreve a classe em texto para os relatorios.
func (c Classe) String() string {
	switch c {
	case CircuitoEuleriano:
		return "Euleriano (possui circuito euleriano)"
	case CaminhoEuleriano:
		return "Semi-euleriano (possui caminho euleriano)"
	default:
		return "Nao-euleriano"
	}
}

// ResultadoEuler e o contrato entre a validacao (Classifica) e o Hierholzer +
// relatorio. Preenchido por Classifica.
type ResultadoEuler struct {
	Classe         Classe
	Conexo         bool           // conexo (grafo) / fracamente conexo (digrafo) ignorando isolados
	GrausImpares   []string       // grafos: vertices de grau impar
	Desbalanceados map[string]int // digrafos: vertice -> (saida-entrada) quando != 0
	VerticeInicial string         // de onde o Hierholzer deve partir
}

// TrilhaEuler e a saida do Hierholzer: a sequencia ordenada de vertices.
type TrilhaEuler struct {
	Sequencia []string
	Circuito  bool
}

// Classifica decide os criterios de eulerianidade (grafo se !Direcionado,
// digrafo caso contrario).
//
// TODO: implementar os criterios.
//   - Grafo: conexo (ignorando isolados) + nº de vertices de grau impar:
//     0 => CircuitoEuleriano (inicio = g.Vertices[0]);
//     2 => CaminhoEuleriano (inicio = um dos impares);
//     senao => NaoEuleriano.
//   - Digrafo: fracamente conexo + para todo v entrada==saida => Circuito;
//     exatamente um v com saida-entrada=+1 (inicio) e um com +1 de entrada (fim)
//     => Caminho; senao => NaoEuleriano.
//     Reusar algoritmos.EhConexo / algoritmos.BFS e o calculo de grau de
//     entrada que ja existe em relatorio.FormataGraus.
func Classifica(g *grafo.Grafo) ResultadoEuler {
	return ResultadoEuler{
		Classe:         NaoEuleriano,
		Desbalanceados: map[string]int{},
	}
}
