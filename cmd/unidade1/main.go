package main

import (
	"fmt"
	"os"

	"github.com/PauloFH/grafos-2026/internal/algoritmos"
	"github.com/PauloFH/grafos-2026/internal/conversoes"
	"github.com/PauloFH/grafos-2026/internal/grafo"
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
		processarGrafo(g, r, nome, saidas)
		r.Salva(saidas)
		r.SalvaPNG(saidas, g)
	}

	fmt.Println("Concluido. Saidas em:", saidas)
}

func processarGrafo(g *grafo.Grafo, r *relatorio.Relatorio, nome, saidas string) {
	if g.Direcionado {
		processarDigrafo(g, r, nome, saidas)
	}

	adicionarDadosBasicos(g, r)

	if nome == "GRAFO_1" || nome == "GRAFO_3" {
		processarBFSDFS(g, r, nome, saidas)
	}
	if nome == "GRAFO_3" {
		r.Adiciona("ARTICULACOES_E_BLOCOS", relatorio.FormataBiconectividade(g))
	}
	if nome == "GRAFO_1" || nome == "GRAFO_2" {
		r.Adiciona("BIPARTIDO", relatorio.FormataBipartido(g))
	}
	if nome == "GRAFO_1" {
		r.Adiciona("OPERACOES_SOBRE_VERTICES", relatorio.FormataOperacoesVertices(g))
		gerarPNGsModificacao(g, nome, saidas)
	}
}

func processarDigrafo(g *grafo.Grafo, r *relatorio.Relatorio, nome, saidas string) {
	subjacente := algoritmos.DeterminaGrafoSubjacente(g)
	r.Adiciona("GRAFO_SUBJACENTE", relatorio.FormataLista(subjacente))

	resDFS := algoritmos.DFSDigrafo(g)
	txtTempos, txtArestas := algoritmos.FormatarDFS(g, resDFS)
	r.Adiciona("DFS_TEMPOS_ENTRADA_SAIDA", txtTempos)
	r.Adiciona("DFS_CLASSIFICACAO_ARESTAS", txtArestas)

	matriz, vertices := algoritmos.MatrizIncidencia(g)
	vetorA, vetorIP := algoritmos.EstrelaDireta(matriz, vertices)
	r.Adiciona("ESTRELA_DIRETA", algoritmos.FormatarEstrelaDireta(vetorA, vetorIP, vertices))

	if nome == "DIGRAFO_DO_GRUPO" {
		processarBFSDFSDigrafo(g, r, resDFS, nome, saidas)
	}
}

func processarBFSDFSDigrafo(g *grafo.Grafo, r *relatorio.Relatorio, resDFS *algoritmos.DFSResult, nome, saidas string) {
	inicio := g.Vertices[0]
	resBFS := algoritmos.BFS(g, inicio)
	r.Adiciona("BFS", relatorio.FormataBFS(resBFS, inicio))

	if err := relatorio.GerarPNGBFS(g, resBFS, inicio, nome+"_BFS", saidas); err != nil {
		fmt.Println("Aviso: erro ao gerar PNG BFS para", nome, ":", err)
	}
	if err := relatorio.GerarPNGDFSDigrafo(g, resDFS, nome+"_DFS", saidas); err != nil {
		fmt.Println("Aviso: erro ao gerar PNG DFS para", nome, ":", err)
	}
}

func adicionarDadosBasicos(g *grafo.Grafo, r *relatorio.Relatorio) {
	r.Adiciona("VERTICES", relatorio.FormataVertices(g))
	r.Adiciona("ARESTAS", relatorio.FormataArestas(g))
	r.Adiciona("LISTA_DE_ADJACENCIA", relatorio.FormataLista(g))
	m := conversoes.ListaParaMatriz(g)
	r.Adiciona("MATRIZ_DE_ADJACENCIA", relatorio.FormataMatriz(g, m))
	conversoes.MatrizParaLista(g, m)
	r.Adiciona("LISTA_RECONVERTIDA_DA_MATRIZ", relatorio.FormataLista(g))
	r.Adiciona("SAO_ADJACENTES", relatorio.FormataAdjacentes(g))
	r.Adiciona("GRAU_DOS_VERTICES", relatorio.FormataGraus(g))
	r.Adiciona("CONEXO", relatorio.FormataConexo(g))
	r.Adiciona("CONTAGEM", relatorio.FormataContagem(g))
	mi, arestas := conversoes.MatrizIncidencia(g)
	r.Adiciona("MATRIZ_DE_INCIDENCIA", relatorio.FormataMatrizIncidencia(g, mi, arestas))
}

func processarBFSDFS(g *grafo.Grafo, r *relatorio.Relatorio, nome, saidas string) {
	inicio := g.Vertices[0]

	resBFS := algoritmos.BFS(g, inicio)
	r.Adiciona("BFS", relatorio.FormataBFS(resBFS, inicio))
	if err := relatorio.GerarPNGBFS(g, resBFS, inicio, nome+"_BFS", saidas); err != nil {
		fmt.Println("Aviso: erro ao gerar PNG BFS para", nome, ":", err)
	}

	resDFS := algoritmos.DFS(g, inicio)
	r.Adiciona("DFS", relatorio.FormataDFS(resDFS, inicio))
	if err := relatorio.GerarPNGDFS(g, resDFS, inicio, nome+"_DFS", saidas); err != nil {
		fmt.Println("Aviso: erro ao gerar PNG DFS para", nome, ":", err)
	}
}

func gerarPNGsModificacao(g *grafo.Grafo, nome, saidas string) {
	gen := relatorio.GeradorPNG{}

	gAdd := g.Clone()
	gAdd.AdicionarVertice("x")
	gAdd.AdicionarAresta(gAdd.Vertices[0], "x")
	if err := gen.Gera(gAdd, nome+"_ADD_VERTEX", saidas); err != nil {
		fmt.Println("Aviso: erro ao gerar PNG add vertex:", err)
	}

	gRem := gAdd.Clone()
	if len(g.Vertices) > 0 {
		gRem.RemoverVertice(g.Vertices[0])
	}
	if err := gen.Gera(gRem, nome+"_REMOVE_VERTEX", saidas); err != nil {
		fmt.Println("Aviso: erro ao gerar PNG remove vertex:", err)
	}
}
