package main

import (
	"fmt"
	"os"

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

		r.Adiciona("VERTICES", relatorio.FormataVertices(g))
		r.Adiciona("ARESTAS", relatorio.FormataArestas(g))
		r.Adiciona("LISTA_DE_ADJACENCIA", relatorio.FormataLista(g))

		// -------------------------------------------------------
		// Veja o README para saber como fazer a adição de seções.
		// -------------------------------------------------------
		r.Salva(saidas)
		r.SalvaPNG(saidas, g)
	}

	fmt.Println("Concluido. Saidas em:", saidas)
}
