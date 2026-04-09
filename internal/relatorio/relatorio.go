package relatorio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/grafo"
)

// Relatorio guarda as seções de saída de um grafo
type Relatorio struct {
	Nome   string
	Secoes []Secao
}

// Secao é uma parte do relatório
type Secao struct {
	Titulo   string
	Conteudo string
}

// Novo cria um relatório vazio
func Novo(nome string) *Relatorio {
	return &Relatorio{
		Nome:   nome,
		Secoes: make([]Secao, 0),
	}
}

// Adiciona uma seção ao relatório
func (r *Relatorio) Adiciona(titulo, conteudo string) {
	r.Secoes = append(r.Secoes, Secao{Titulo: titulo, Conteudo: conteudo})
}

// Texto gera o relatório completo como string
func (r *Relatorio) Texto() string {
	var sb strings.Builder

	sb.WriteString("==============================================\n")
	sb.WriteString("RELATORIO: " + r.Nome + "\n")
	sb.WriteString("==============================================\n\n")

	for _, s := range r.Secoes {
		sb.WriteString("--- " + s.Titulo + " ---\n")
		sb.WriteString(s.Conteudo)
		sb.WriteString("\n")
	}

	return sb.String()
}

// Salva escreve o relatório em um arquivo
func (r *Relatorio) Salva(caminho string) {
	os.MkdirAll(caminho, 0755)
	arquivo := filepath.Join(caminho, r.Nome+".txt")
	os.WriteFile(arquivo, []byte(r.Texto()), 0644)
}

// SalvaPNG gera um PNG do grafo em caminho/Nome.png
func (r *Relatorio) SalvaPNG(caminho string, g *grafo.Grafo) {
	os.MkdirAll(caminho, 0755)
	gen := GeradorPNG{}
	if err := gen.Gera(g, r.Nome, caminho); err != nil {
		fmt.Println("Aviso: nao foi possivel gerar PNG para", r.Nome, "-", err)
	}
}

// Imprime exibe no terminal
func (r *Relatorio) Imprime() {
	fmt.Print(r.Texto())
}

// --------------------------------------------------------
// Funções prontas de formatação
// --------------------------------------------------------

// FormataLista gera o texto da lista de adjacência
func FormataLista(g *grafo.Grafo) string {
	var sb strings.Builder
	for _, v := range g.Vertices {
		vizinhos := g.ListaAdj[v]
		sb.WriteString(v + " -> ")
		if len(vizinhos) > 0 {
			sb.WriteString(strings.Join(vizinhos, ", "))
		} else {
			sb.WriteString("(vazio)")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// FormataVertices gera info básica dos vértices
func FormataVertices(g *grafo.Grafo) string {
	return fmt.Sprintf("  Total de vertices: %d\n  Vertices: %s\n",
		g.NumVertices(), strings.Join(g.Vertices, ", "))
}

// FormataArestas gera info básica das arestas
func FormataArestas(g *grafo.Grafo) string {
	return fmt.Sprintf("  Total de arestas: %d\n", g.NumArestas())
}
