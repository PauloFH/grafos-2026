// TODOs os pontos de integracao das Partes 1 a 4.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/pcv"
)

// semente da primeira execucao dos metodos estocasticos;
const sementeBase int64 = 42

const execucoesEstocasticas = 20

func main() {
	entradas := "inputs_u3"
	saidas := "outputs_u3"

	fmt.Println("========================================")
	fmt.Println("  TRABALHO DE GRAFOS - 2026 (UNIDADE 3)")
	fmt.Println("  Problema do Caixeiro Viajante (PCV)")
	fmt.Println("========================================")
	fmt.Println()

	instancias, err := pcv.CarregaDiretorio(entradas)
	if err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(saidas, 0o755); err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}

	fmt.Println("Instancias encontradas:", len(instancias))
	fmt.Println()
	imprimeTabelaInstancias(instancias)
	fmt.Println()

	for _, in := range instancias {
		processaInstancia(in, saidas)
	}

	// TODO(Parte 4): gerar outputs_u3/COMPARATIVO_GERAL.txt
	// (Instancia | Otimo | VMP+2opt | IMB+OrOpt | AG | Memetico, com gaps;
	// AG e Memetico entram com o MENOR das 20 execucoes).

	fmt.Println("Concluido. Saidas em:", saidas)
}

func imprimeTabelaInstancias(instancias []*pcv.Instancia) {
	fmt.Printf("%-12s | %-6s | %3s | %s\n", "Instancia", "Medida", "N", "Melhor conhecido")
	fmt.Println(strings.Repeat("-", 50))
	for _, in := range instancias {
		otimo := "-"
		if v, ok := pcv.ValoresOtimos[in.Nome]; ok {
			otimo = fmt.Sprintf("%.2f", v)
		}
		fmt.Printf("%-12s | %-6s | %3d | %s\n", in.Nome, in.Medida, in.N, otimo)
	}
}

func processaInstancia(in *pcv.Instancia, saidas string) {
	fmt.Printf("[%s] processando (%s, N=%d)...\n", in.Nome, in.Medida, in.N)

	// Heuristicas construtivas (deterministicas, 1 execucao cada).
	otimo := pcv.ValoresOtimos[in.Nome]

	vmp := pcv.VizinhoMaisProximo{}.Constroi(in)
	// TODO: depois := pcv.DoisOpt{}.Aplica(vmp, in)
	fmt.Printf("[%s] VMP: custo=%.2f gap=%.2f%%\n", in.Nome, vmp.Custo, pcv.Gap(vmp.Custo, otimo))

	imb := pcv.InsercaoMaisBarata{}.Constroi(in)
	// TODO: depois := pcv.OrOpt{}.Aplica(imb, in)
	fmt.Printf("[%s] IMB: custo=%.2f gap=%.2f%%\n", in.Nome, imb.Custo, pcv.Gap(imb.Custo, otimo))

	ag := pcv.AlgoritmoGenetico{Par: pcv.ParametrosPadrao()}
	resumoAG := pcv.ExecutaExperimento(ag, in, execucoesEstocasticas, sementeBase)
	caminhoAG := filepath.Join(saidas, "RESUMO_AG_"+in.Nome+".txt")
	if err := os.WriteFile(caminhoAG, []byte(resumoAG.TextoResumo(in)), 0o644); err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}
	fmt.Printf("[%s] AG: menor=%.2f media=%.2f -> %s\n",
		in.Nome, resumoAG.Melhor, resumoAG.Media, caminhoAG)

	// TODO(Parte 4): Memetico (pop=50, ger=200, Buscas={DoisOpt, OrOpt,
	// Swap}), 20 execucoes; gravar outputs_u3/RESUMO_MEMETICO_<in.Nome>.txt.

	// TODO(Parte 4): relatorio por instancia (outputs_u3/<in.Nome>.txt) com
	// as secoes INSTANCIA, VIZINHO_MAIS_PROXIMO, INSERCAO_MAIS_BARATA,
	// ALGORITMO_GENETICO e ALGORITMO_MEMETICO.
}
