package relatorio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/algoritmos"
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
		algoritmos.TotalVertices(g), strings.Join(g.Vertices, ", "))
}

// FormataArestas gera info básica das arestas
func FormataArestas(g *grafo.Grafo) string {
	return fmt.Sprintf("  Total de arestas: %d\n", algoritmos.TotalArestas(g))
}

const errMatrizInvalida = "  (matriz com dimensoes invalidas)\n"

// FormataMatriz gera o texto da matriz de adjacência
func FormataMatriz(g *grafo.Grafo, m [][]int) string {
	n := len(g.Vertices)
	if len(m) != n {
		return "  (matriz nao gerada)\n"
	}
	for i := range m {
		if len(m[i]) != n {
			return errMatrizInvalida
		}
	}
	var sb strings.Builder

	fmt.Fprintf(&sb, "%5s", "")
	for _, v := range g.Vertices {
		fmt.Fprintf(&sb, "%4s", v)
	}
	sb.WriteByte('\n')

	for i, v := range g.Vertices {
		fmt.Fprintf(&sb, "%4s ", v)
		for _, val := range m[i] {
			fmt.Fprintf(&sb, "%4d", val)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// FormataAdjacentes lista todos os pares de vértices adjacentes
func FormataAdjacentes(g *grafo.Grafo) string {
	var sb strings.Builder
	sep := " - "
	if g.Direcionado {
		sep = " -> "
	}
	for _, par := range algoritmos.ParesAdjacentes(g) {
		fmt.Fprintf(&sb, "  %s%s%s\n", par[0], sep, par[1])
	}
	return sb.String()
}

// FormataConexo indica se o grafo é conexo
func FormataConexo(g *grafo.Grafo) string {
	if g.Direcionado {
		if algoritmos.EhConexo(g) {
			return "  O digrafo e fracamente conexo.\n"
		}
		return "  O digrafo NAO e fracamente conexo.\n"
	}
	if algoritmos.EhConexo(g) {
		return "  O grafo e conexo.\n"
	}
	return "  O grafo NAO e conexo.\n"
}

// FormataContagem exibe total de vértices e arestas
func FormataContagem(g *grafo.Grafo) string {
	return fmt.Sprintf("  Total de vertices: %d\n  Total de arestas: %d\n",
		algoritmos.TotalVertices(g), algoritmos.TotalArestas(g))
}

// FormataMatrizIncidencia exibe a matriz de incidência com rótulos de colunas
func FormataMatrizIncidencia(g *grafo.Grafo, m [][]int, arestas [][2]string) string {
	if len(arestas) == 0 {
		return "  (sem arestas)\n"
	}
	if len(m) != len(g.Vertices) {
		return errMatrizInvalida
	}
	for i := range m {
		if len(m[i]) != len(arestas) {
			return errMatrizInvalida
		}
	}

	sep := "-"
	if g.Direcionado {
		sep = "->"
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "%5s", "")
	for _, a := range arestas {
		fmt.Fprintf(&sb, "%7s", a[0]+sep+a[1])
	}
	sb.WriteByte('\n')

	for i, v := range g.Vertices {
		fmt.Fprintf(&sb, "%4s ", v)
		for _, val := range m[i] {
			fmt.Fprintf(&sb, "%7d", val)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
