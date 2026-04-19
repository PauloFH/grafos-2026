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

const dotRankdir = "  rankdir=TB;\n"
const dotNodeShape = "  node [shape=ellipse];\n"

func GerarPNGBFS(g *grafo.Grafo, res algoritmos.ResultadoBFS, inicio, nome, caminho string) error {
	return executarDot(gerarDOTBFS(g, res, inicio, nome), nome, caminho)
}

func GerarPNGDFS(g *grafo.Grafo, res algoritmos.ResultadoDFS, inicio, nome, caminho string) error {
	return executarDot(gerarDOTDFS(g, res, inicio, nome), nome, caminho)
}

func GerarPNGDFSDigrafo(g *grafo.Grafo, res *algoritmos.DFSResult, nome, caminho string) error {
	return executarDot(gerarDOTDFSDigrafo(g, res, nome), nome, caminho)
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

func dotHeader(graphType, nome string) string {
	clean := strings.ReplaceAll(nome, " ", "_")
	return graphType + " " + clean + " {\n" + dotRankdir + dotNodeShape
}

func dotOperador(direcionado bool) (string, string) {
	if direcionado {
		return "digraph", "->"
	}
	return "graph", "--"
}

// escreverVerticesBFS emite nós com destaque no vértice inicial
func escreverVerticesBFS(sb *strings.Builder, g *grafo.Grafo, inicio string) {
	for _, v := range g.Vertices {
		if v == inicio {
			fmt.Fprintf(sb, "  \"%s\" [fillcolor=lightgreen, style=filled];\n", v)
		} else {
			fmt.Fprintf(sb, "  \"%s\";\n", v)
		}
	}
}

// buildTreeKeysBFS constrói o conjunto de chaves das arestas de árvore e emite as linhas DOT
func buildTreeKeysBFS(sb *strings.Builder, res algoritmos.ResultadoBFS, op string, direcionado bool) map[string]bool {
	keys := map[string]bool{}
	for _, v := range res.Visitados {
		pred, ok := res.Predecessor[v]
		if !ok {
			continue
		}
		keys[pred+op+v] = true
		if !direcionado {
			keys[v+op+pred] = true
		}
		fmt.Fprintf(sb, "  \"%s\" %s \"%s\";\n", pred, op, v)
	}
	return keys
}

func arestaPular(chave, inv string, direcionado bool, treeKeys, vistas map[string]bool) bool {
	if treeKeys[chave] || vistas[chave] {
		return true
	}
	return !direcionado && (treeKeys[inv] || vistas[inv])
}

func dotArestaDashed(sb *strings.Builder, u, op, w, label string) {
	if label != "" {
		fmt.Fprintf(sb, "  \"%s\" %s \"%s\" [style=dashed, color=red, label=\"%s\"];\n", u, op, w, label)
	} else {
		fmt.Fprintf(sb, "  \"%s\" %s \"%s\" [style=dashed];\n", u, op, w)
	}
}

// escreverNaoArvore emite arestas não-árvore tracejadas (genérico BFS/DFS)
func escreverNaoArvore(sb *strings.Builder, g *grafo.Grafo, op string, treeKeys map[string]bool, label string) {
	vistas := map[string]bool{}
	for _, u := range g.Vertices {
		for _, w := range g.GetVizinhos(u) {
			chave := u + op + w
			if arestaPular(chave, w+op+u, g.Direcionado, treeKeys, vistas) {
				continue
			}
			vistas[chave] = true
			dotArestaDashed(sb, u, op, w, label)
		}
	}
}

// escreverNiveis emite blocos rank=same por nível BFS
func escreverNiveis(sb *strings.Builder, res algoritmos.ResultadoBFS) {
	levelMap := map[int][]string{}
	for _, v := range res.Visitados {
		lvl := res.Nivel[v]
		levelMap[lvl] = append(levelMap[lvl], v)
	}
	levels := make([]int, 0, len(levelMap))
	for lvl := range levelMap {
		levels = append(levels, lvl)
	}
	sort.Ints(levels)
	for _, lvl := range levels {
		sb.WriteString("  { rank=same;")
		for _, v := range levelMap[lvl] {
			fmt.Fprintf(sb, " \"%s\";", v)
		}
		sb.WriteString(" }\n")
	}
}

func gerarDOTBFS(g *grafo.Grafo, res algoritmos.ResultadoBFS, inicio, nome string) string {
	gType, op := dotOperador(g.Direcionado)
	var sb strings.Builder
	sb.WriteString(dotHeader(gType, nome))
	escreverVerticesBFS(&sb, g, inicio)
	treeKeys := buildTreeKeysBFS(&sb, res, op, g.Direcionado)
	escreverNaoArvore(&sb, g, op, treeKeys, "")
	escreverNiveis(&sb, res)
	sb.WriteString("}\n")
	return sb.String()
}

// buildTreeKeysDFS constrói chaves de árvore e emite arestas DFS
func buildTreeKeysDFS(sb *strings.Builder, res algoritmos.ResultadoDFS, op string, direcionado bool) map[string]bool {
	keys := map[string]bool{}
	for v, pred := range res.Predecessor {
		keys[pred+op+v] = true
		if !direcionado {
			keys[v+op+pred] = true
		}
	}
	for _, v := range res.Visitados {
		if pred, ok := res.Predecessor[v]; ok {
			fmt.Fprintf(sb, "  \"%s\" %s \"%s\";\n", pred, op, v)
		}
	}
	return keys
}

// escreverRetornoDFS emite arestas de retorno filtradas por tempo de entrada
func escreverRetornoDFS(sb *strings.Builder, g *grafo.Grafo, res algoritmos.ResultadoDFS, op string, treeKeys map[string]bool) {
	vistas := map[string]bool{}
	for _, u := range res.Visitados {
		for _, w := range g.GetVizinhos(u) {
			chave := u + op + w
			inv := w + op + u
			if treeKeys[chave] || vistas[chave] {
				continue
			}
			if !g.Direcionado && (treeKeys[inv] || vistas[inv]) {
				continue
			}
			if res.Entrada[w] < res.Entrada[u] {
				vistas[chave] = true
				fmt.Fprintf(sb, "  \"%s\" %s \"%s\" [style=dashed, color=red, label=\"R\"];\n", u, op, w)
			}
		}
	}
}

func gerarDOTDFS(g *grafo.Grafo, res algoritmos.ResultadoDFS, inicio, nome string) string {
	gType, op := dotOperador(g.Direcionado)
	var sb strings.Builder
	sb.WriteString(dotHeader(gType, nome))
	escreverVerticesBFS(&sb, g, inicio) // mesmo comportamento de highlight
	treeKeys := buildTreeKeysDFS(&sb, res, op, g.Direcionado)
	escreverRetornoDFS(&sb, g, res, op, treeKeys)
	sb.WriteString("}\n")
	return sb.String()
}

var dotEdgeColors = map[string]string{
	"Árvore":     "black",
	"Retorno":    "red",
	"Avanço":     "blue",
	"Cruzamento": "green",
}

func gerarDOTDFSDigrafo(g *grafo.Grafo, res *algoritmos.DFSResult, nome string) string {
	var sb strings.Builder
	sb.WriteString(dotHeader("digraph", nome))

	for _, v := range g.Vertices {
		label := fmt.Sprintf("%s\\n[%d/%d]", v, res.Discovery[v], res.Finish[v])
		fmt.Fprintf(&sb, "  \"%s\" [label=\"%s\"];\n", v, label)
	}

	for _, e := range res.Edges {
		from, to, edgeType, ok := parseEdgeEntry(e)
		if !ok {
			continue
		}
		cor := dotEdgeColors[edgeType]
		if cor == "" {
			cor = "black"
		}
		fmt.Fprintf(&sb, "  \"%s\" -> \"%s\" [color=%s, label=\"%s\"];\n", from, to, cor, edgeType)
	}

	sb.WriteString("}\n")
	return sb.String()
}

// parseEdgeEntry parseia strings no formato "(u, v): Tipo"
func parseEdgeEntry(s string) (from, to, edgeType string, ok bool) {
	start := strings.Index(s, "(")
	end := strings.Index(s, ")")
	if start < 0 || end <= start {
		return
	}
	parts := strings.SplitN(s[start+1:end], ", ", 2)
	if len(parts) != 2 {
		return
	}
	colonIdx := strings.Index(s, ": ")
	if colonIdx < 0 {
		return
	}
	return parts[0], parts[1], s[colonIdx+2:], true
}
