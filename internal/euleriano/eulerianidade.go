// Package euleriano implementa a Unidade 2 (Projeto 3): validacao de
// eulerianidade e extracao de trilhas eulerianas (Hierholzer) para grafos e
// digrafos. Reusa o nucleo compartilhado em internal/{grafo,algoritmos}.
package euleriano

import (
	"github.com/PauloFH/grafos-2026/internal/algoritmos"
	"github.com/PauloFH/grafos-2026/internal/grafo"
)

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

func Classifica(g *grafo.Grafo) ResultadoEuler {

	if len(g.Vertices) == 0 {
		return ResultadoEuler{
			Classe:         CircuitoEuleriano,
			Conexo:         true,
			GrausImpares:   []string{},
			Desbalanceados: map[string]int{},
		}
	}

	if g.Direcionado {
		return classificaDigrafo(g)
	}

	return classificaGrafo(g)
}

func classificaGrafo(g *grafo.Grafo) ResultadoEuler {

	//verificando se é conexo
	conexo := algoritmos.EhConexo(g)

	// se não for conexo, é não-euleriano
	if !conexo {
		return ResultadoEuler{
			Classe:         NaoEuleriano,
			Conexo:         false,
			Desbalanceados: map[string]int{},
		}
	}

	var impares []string

	//verificando e contando quantos impares existem
	for _, v := range g.Vertices {

		if g.GrauVertice(v)%2 != 0 {
			impares = append(impares, v)
		}
	}

	//se não tiver vertice com grau impares é um circuito euleriano
	if len(impares) == 0 {

		return ResultadoEuler{
			Classe:         CircuitoEuleriano,
			Conexo:         true,
			GrausImpares:   impares,
			Desbalanceados: map[string]int{},
			VerticeInicial: g.Vertices[0],
		}
	}

	//se tiver dois vertices com graus impares, existe um caminho euleriano
	if len(impares) == 2 {

		return ResultadoEuler{
			Classe:         CaminhoEuleriano,
			Conexo:         true,
			GrausImpares:   impares,
			Desbalanceados: map[string]int{},
			VerticeInicial: impares[0],
		}
	}

	//se tiver mais dois vertices com graus impares, não existem nem caminho nem circuito
	return ResultadoEuler{
		Classe:         NaoEuleriano,
		Conexo:         true,
		GrausImpares:   impares,
		Desbalanceados: map[string]int{},
	}
}

func classificaDigrafo(g *grafo.Grafo) ResultadoEuler {

	conexo := algoritmos.EhConexo(g)

	if !conexo {
		return ResultadoEuler{
			Classe:         NaoEuleriano,
			Conexo:         false,
			Desbalanceados: map[string]int{},
		}
	}

	// Armazena o grau de entrada e saída de cada vértice.
	entrada := make(map[string]int)
	saida := make(map[string]int)

	// Calcula grau de saída e grau de entrada
	for _, v := range g.Vertices {

		saida[v] = len(g.ListaAdj[v])

		for _, destino := range g.ListaAdj[v] {
			entrada[destino]++
		}
	}

	desbalanceados := make(map[string]int)

	var inicio string
	inicioCount := 0
	fimCount := 0

	for _, v := range g.Vertices {

		//diferença de entrada e saida
		diff := saida[v] - entrada[v]

		if diff != 0 {
			desbalanceados[v] = diff
		}

		if diff == 1 {
			inicioCount++
			inicio = v
		}

		if diff == -1 {
			fimCount++
		}
	}

	// Circuito Euleriano:
	// todos os vértices possuem entrada == saída
	if len(desbalanceados) == 0 {
		return ResultadoEuler{
			Classe:         CircuitoEuleriano,
			Conexo:         true,
			Desbalanceados: desbalanceados,
			VerticeInicial: g.Vertices[0],
		}
	}

	// Caminho Euleriano:
	// um vértice com +1 e um vértice com -1
	if inicioCount == 1 && fimCount == 1 && len(desbalanceados) == 2 {
		return ResultadoEuler{
			Classe:         CaminhoEuleriano,
			Conexo:         true,
			Desbalanceados: desbalanceados,
			VerticeInicial: inicio,
		}
	}

	return ResultadoEuler{
		Classe:         NaoEuleriano,
		Conexo:         true,
		Desbalanceados: desbalanceados,
	}
}
