package relatorio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/pcv"
)

type PontoPlano struct {
	X, Y float64
	Id   int    // id original da cidade (1..48)
	Nome string // nome canonico
}

const escalaPNG = 760.0 // coordenadas [0,1] -> pontos

func GraphvizDisponivel() bool {
	_, err := exec.LookPath("neato")
	return err == nil
}

func GerarPNGRotaPCV(in *pcv.Instancia, pontos []PontoPlano, rota pcv.Rota, titulo, nome, caminho string) error {
	return executarNeato(gerarDOTRotaPCV(in, pontos, rota, titulo), nome, caminho)
}

func GerarPNGGrafoPCV(in *pcv.Instancia, pontos []PontoPlano, titulo, nome, caminho string) error {
	return executarNeato(gerarDOTGrafoPCV(in, pontos, titulo), nome, caminho)
}

func gerarDOTRotaPCV(in *pcv.Instancia, pontos []PontoPlano, rota pcv.Rota, titulo string) string {
	var sb strings.Builder
	sb.WriteString("graph rota {\n")
	sb.WriteString("  bgcolor=\"#111845\";\n")
	sb.WriteString("  outputorder=edgesfirst;\n")
	fmt.Fprintf(&sb, "  label=%q;\n", titulo)
	sb.WriteString("  labelloc=t; fontcolor=\"#ee4c9b\"; fontsize=20; fontname=\"sans-serif\";\n")
	sb.WriteString("  node [shape=circle, style=filled, fontname=\"sans-serif\", fontsize=11, fixedsize=true, width=0.3];\n")

	for pos, p := range pontos { // deposito (pos 0) em ambar
		fill, cor, borda := "#232c68", "#f4f6ff", "#ee4c9b"
		if pos == 0 {
			fill, borda = "#4a3407", "#f59e0b"
		}
		fmt.Fprintf(&sb, "  n%d [label=%q, pos=\"%.3f,%.3f!\", fillcolor=%q, fontcolor=%q, color=%q];\n",
			pos, fmt.Sprintf("%d", p.Id), p.X*escalaPNG, (1-p.Y)*escalaPNG, fill, cor, borda)
	}

	sb.WriteString("  edge [color=\"#ee4c9b\", penwidth=2.2];\n")
	n := len(rota.Ordem)
	for i := 0; i < n; i++ { // ciclo fechado
		fmt.Fprintf(&sb, "  n%d -- n%d;\n", rota.Ordem[i], rota.Ordem[(i+1)%n])
	}

	sb.WriteString("}\n")
	return sb.String()
}

func gerarDOTGrafoPCV(in *pcv.Instancia, pontos []PontoPlano, titulo string) string {
	var sb strings.Builder
	sb.WriteString("graph instancia {\n")
	sb.WriteString("  bgcolor=\"#111845\";\n")
	sb.WriteString("  outputorder=edgesfirst;\n")
	fmt.Fprintf(&sb, "  label=%q;\n", titulo)
	sb.WriteString("  labelloc=t; fontcolor=\"#ee4c9b\"; fontsize=20; fontname=\"sans-serif\";\n")
	sb.WriteString("  node [shape=circle, style=filled, fontname=\"sans-serif\", fontsize=11, fixedsize=true, width=0.3];\n")

	for pos, p := range pontos { // deposito (pos 0) em ambar
		fill, cor, borda := "#232c68", "#f4f6ff", "#ee4c9b"
		if pos == 0 {
			fill, borda = "#4a3407", "#f59e0b"
		}
		fmt.Fprintf(&sb, "  n%d [label=%q, pos=\"%.3f,%.3f!\", fillcolor=%q, fontcolor=%q, color=%q];\n",
			pos, fmt.Sprintf("%d", p.Id), p.X*escalaPNG, (1-p.Y)*escalaPNG, fill, cor, borda)
	}
	sb.WriteString("  edge [color=\"#8fa4b5\", penwidth=0.35];\n")
	for i := 0; i < in.N; i++ { // grafo completo
		for j := i + 1; j < in.N; j++ {
			fmt.Fprintf(&sb, "  n%d -- n%d;\n", i, j)
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

func executarNeato(dot, nome, caminho string) error {
	dotFile := filepath.Join(caminho, nome+".dot")
	pngFile := filepath.Join(caminho, nome+".png")
	if err := os.WriteFile(dotFile, []byte(dot), 0o644); err != nil {
		return err
	}
	defer os.Remove(dotFile)
	cmd := exec.Command("neato", "-n", "-Tpng", dotFile, "-o", pngFile)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("neato: %w\n%s", err, out)
	}
	return nil
}
