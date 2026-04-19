package relatorio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/grafo"
)

// GeradorPNG converte um grafo para imagem PNG usando Graphviz 
type GeradorPNG struct{}

// Gera salva um PNG do grafo em caminho/nome.png
func (gp GeradorPNG) Gera(g *grafo.Grafo, nome, caminho string) error {
	dot := gerarDOT(g, nome)

	dotFile := filepath.Join(caminho, nome+".dot")
	pngFile := filepath.Join(caminho, nome+".png")

	if err := os.WriteFile(dotFile, []byte(dot), 0644); err != nil {
		return err
	}
	defer os.Remove(dotFile)

	cmd := exec.Command("dot", "-Tpng", dotFile, "-o", pngFile)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("dot: %w\n%s", err, out)
	}
	return nil
}

func gerarDOT(g *grafo.Grafo, nome string) string {
	var sb strings.Builder

	nome = strings.ReplaceAll(nome, " ", "_")
	if nome == "" {
		nome = "G"
	}

	if g.Direcionado {
		sb.WriteString("digraph " + nome + " {\n")
	} else {
		sb.WriteString("graph " + nome + " {\n")
	}

	sb.WriteString("  node [shape=circle];\n")

	for _, v := range g.Vertices {
		sb.WriteString(fmt.Sprintf("  \"%s\";\n", v))
	}

	op := "--"
	if g.Direcionado {
		op = "->"
	}

	visitados := map[string]bool{}
	for _, origem := range g.Vertices {
		for _, dest := range g.ListaAdj[origem] {
			chave := origem + op + dest
			chaveInv := dest + op + origem
			if !g.Direcionado && visitados[chaveInv] {
				continue
			}
			visitados[chave] = true
			sb.WriteString(fmt.Sprintf("  \"%s\" %s \"%s\";\n", origem, op, dest))
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}
