package relatorio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/algoritmos"
	"github.com/PauloFH/grafos-2026/internal/grafo"
)

func GerarPNGBFS(g *grafo.Grafo, res algoritmos.ResultadoBFS, inicio, nome, caminho string) error {
	dot := gerarDOTBFS(g, res, inicio, nome)
	return executarDot(dot, nome, caminho)
}

func GerarPNGDFS(g *grafo.Grafo, res algoritmos.ResultadoDFS, inicio, nome, caminho string) error {
	dot := gerarDOTDFS(g, res, inicio, nome)
	return executarDot(dot, nome, caminho)
}

func executarDot(dot, nome, caminho string) error {
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

func gerarDOTBFS(g *grafo.Grafo, res algoritmos.ResultadoBFS, inicio, nome string) string {
	var sb strings.Builder

	nomeClean := strings.ReplaceAll(nome, " ", "_")
	sb.WriteString("graph " + nomeClean + " {\n")
	sb.WriteString("  rankdir=TB;\n")
	sb.WriteString("  node [shape=ellipse];\n")

	// Vértices
	for _, v := range g.Vertices {
		if v == inicio {
			sb.WriteString(fmt.Sprintf("  \"%s\" [fillcolor=lightgreen, style=filled];\n", v))
		} else {
			sb.WriteString(fmt.Sprintf("  \"%s\";\n", v))
		}
	}

	// Detectar arestas de árvore
	treeEdges := map[string]bool{}
	for _, v := range res.Visitados {
		if pred, ok := res.Predecessor[v]; ok {
			treeEdges[pred+"--"+v] = true
			treeEdges[v+"--"+pred] = true
		}
	}

	// Arestas de árvore (sólidas)
	for _, v := range res.Visitados {
		if pred, ok := res.Predecessor[v]; ok {
			sb.WriteString(fmt.Sprintf("  \"%s\" -- \"%s\";\n", pred, v))
		}
	}

	// Arestas não-árvore (tracejadas)
	visitadas := map[string]bool{}
	for _, u := range g.Vertices {
		for _, w := range g.GetVizinhos(u) {
			chave := u + "--" + w
			chaveInv := w + "--" + u
			if treeEdges[chave] || treeEdges[chaveInv] {
				continue
			}
			if visitadas[chave] || visitadas[chaveInv] {
				continue
			}
			visitadas[chave] = true
			sb.WriteString(fmt.Sprintf("  \"%s\" -- \"%s\" [style=dashed];\n", u, w))
		}
	}

	// Agrupamento por nível (rank=same)
	levelMap := map[int][]string{}
	for _, v := range res.Visitados {
		lvl := res.Nivel[v]
		levelMap[lvl] = append(levelMap[lvl], v)
	}
	levels := []int{}
	for lvl := range levelMap {
		levels = append(levels, lvl)
	}
	sort.Ints(levels)
	for _, lvl := range levels {
		sb.WriteString("  { rank=same;")
		for _, v := range levelMap[lvl] {
			sb.WriteString(fmt.Sprintf(" \"%s\";", v))
		}
		sb.WriteString(" }\n")
	}

	sb.WriteString("}\n")
	return sb.String()
}

func gerarDOTDFS(g *grafo.Grafo, res algoritmos.ResultadoDFS, inicio, nome string) string {
	var sb strings.Builder

	nomeClean := strings.ReplaceAll(nome, " ", "_")
	sb.WriteString("graph " + nomeClean + " {\n")
	sb.WriteString("  rankdir=TB;\n")
	sb.WriteString("  node [shape=ellipse];\n")

	// Vértices
	for _, v := range g.Vertices {
		if v == inicio {
			sb.WriteString(fmt.Sprintf("  \"%s\" [fillcolor=lightgreen, style=filled];\n", v))
		} else {
			sb.WriteString(fmt.Sprintf("  \"%s\";\n", v))
		}
	}

	// Detectar arestas de árvore
	treeEdges := map[string]bool{}
	for v, pred := range res.Predecessor {
		treeEdges[pred+"--"+v] = true
		treeEdges[v+"--"+pred] = true
	}

	// Arestas de árvore (sólidas)
	for _, v := range res.Visitados {
		if pred, ok := res.Predecessor[v]; ok {
			sb.WriteString(fmt.Sprintf("  \"%s\" -- \"%s\";\n", pred, v))
		}
	}

	// Arestas de retorno (tracejadas vermelhas com "R")
	visitadas := map[string]bool{}
	for _, u := range res.Visitados {
		for _, w := range g.GetVizinhos(u) {
			chave := u + "--" + w
			chaveInv := w + "--" + u
			if treeEdges[chave] || treeEdges[chaveInv] {
				continue
			}
			if visitadas[chave] || visitadas[chaveInv] {
				continue
			}
			// aresta de retorno: w foi descoberto antes de u
			if res.Entrada[w] < res.Entrada[u] {
				visitadas[chave] = true
				sb.WriteString(fmt.Sprintf("  \"%s\" -- \"%s\" [style=dashed, color=red, label=\"R\"];\n", u, w))
			}
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}
