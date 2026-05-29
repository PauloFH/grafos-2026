package main

import (
	"fmt"
	"os"

	"github.com/PauloFH/grafos-2026/internal/algoritmos"
	"github.com/PauloFH/grafos-2026/internal/euleriano"
	"github.com/PauloFH/grafos-2026/internal/grafo"
	"github.com/PauloFH/grafos-2026/internal/leitor"
	"github.com/PauloFH/grafos-2026/internal/relatorio"
)

func main() {
	entradas := "inputs_u2"
	saidas := "outputs_u2"

	fmt.Println("========================================")
	fmt.Println("  TRABALHO DE GRAFOS - 2026 (UNIDADE 2)")
	fmt.Println("  Projeto 3 - Grafos Eulerianos")
	fmt.Println("========================================")
	fmt.Println()

	grafos, err := leitor.LerDiretorio(entradas)
	if err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}

	fmt.Println("Grafos encontrados:", len(grafos))

	for nome, g := range grafos {
		tipo := "GRAFO"
		if g.Direcionado {
			tipo = "DIGRAFO"
		}
		fmt.Printf("[%s] %s - %d vertices, %d arestas\n",
			tipo, nome, algoritmos.TotalVertices(g), algoritmos.TotalArestas(g))

		r := relatorio.Novo(nome)
		processarEuleriano(g, r, nome, saidas)
		r.Salva(saidas)
		r.SalvaPNG(saidas, g)
	}

	fmt.Println("Concluido. Saidas em:", saidas)
}

// processarEuleriano roda a validacao + Hierholzer e monta o relatorio da U2.
func processarEuleriano(g *grafo.Grafo, r *relatorio.Relatorio, nome, saidas string) {
	res := euleriano.Classifica(g)
	r.Adiciona("EULERIANIDADE", relatorio.FormataEulerianidade(res))

	if res.Classe == euleriano.NaoEuleriano {
		r.Adiciona("TRILHA_EULERIANA", "  Grafo nao-euleriano: nao ha trilha a extrair.\n")
		return
	}

	var trilha euleriano.TrilhaEuler
	var err error
	if g.Direcionado {
		trilha, err = euleriano.HierholzerDigrafo(g, res.VerticeInicial)
	} else {
		trilha, err = euleriano.HierholzerGrafo(g, res.VerticeInicial)
	}
	if err != nil {
		r.Adiciona("TRILHA_EULERIANA", "  Erro ao extrair trilha: "+err.Error()+"\n")
		return
	}

	r.Adiciona("TRILHA_EULERIANA", relatorio.FormataTrilhaEuler(trilha))

	if err := relatorio.GerarPNGTrilhaEuler(g, trilha.Sequencia, nome+"_TRILHA", saidas); err != nil {
		fmt.Println("Aviso: erro ao gerar PNG da trilha para", nome, ":", err)
	}
}
