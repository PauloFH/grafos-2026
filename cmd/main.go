package main

import (
	"fmt"
	"os"

	"github.com/PauloFH/grafos-2026/internal/algoritmos"
	"github.com/PauloFH/grafos-2026/internal/leitor"
	"github.com/PauloFH/grafos-2026/internal/relatorio"
)

func main() {
	entradas := "inputs"
	saidas := "outputs"

	fmt.Println("========================================")
	fmt.Println("  TRABALHO DE GRAFOS - 2026")
	fmt.Println("========================================")
	fmt.Println()

	// Lê todos os grafos
	grafos, err := leitor.LerDiretorio(entradas)
	if err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}

	fmt.Println("Grafos encontrados:", len(grafos))

	// Para cada grafo, gera relatório
	for nome, g := range grafos {
		tipo := "GRAFO"
		if g.Direcionado {
			tipo = "DIGRAFO"
		}
		fmt.Printf("[%s] %s - %d vertices, %d arestas\n",
			tipo, nome, g.NumVertices(), g.NumArestas())

		r := relatorio.Novo(nome)

		if g.Direcionado { //Se for digrafo, gera o grafo subjacente e analises especificas
			subjacente := algoritmos.DeterminaGrafoSubjacente(g)
			r.Adiciona("GRAFO_SUBJACENTE", relatorio.FormataLista(subjacente))

			resultadoDFS := algoritmos.DFS(g)
			txtTempos, txtArestas := algoritmos.FormatarDFS(g, resultadoDFS)
			r.Adiciona("DFS_TEMPOS_ENTRADA_SAIDA", txtTempos)
			r.Adiciona("DFS_CLASSIFICACAO_ARESTAS", txtArestas)

			matriz, vertices := algoritmos.MatrizIncidencia(g)
			vetorA, vetorIP := algoritmos.EstrelaDireta(matriz, vertices)
			txtEstrela := algoritmos.FormatarEstrelaDireta(vetorA, vetorIP, vertices)
			r.Adiciona("ESTRELA_DIRETA", txtEstrela)
		}

		// Dados básicos
		r.Adiciona("VERTICES", relatorio.FormataVertices(g))
		r.Adiciona("ARESTAS", relatorio.FormataArestas(g))
		r.Adiciona("LISTA_DE_ADJACENCIA", relatorio.FormataLista(g))

		// -------------------------------------------------------
		// Veja o README para saber como fazer a adição de seções.
		// -------------------------------------------------------
		r.Salva(saidas)
	}

	fmt.Println("Concluido. Saidas em:", saidas)

}
