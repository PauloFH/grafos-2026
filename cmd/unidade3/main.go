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

// ResultadoInstancia armazena os custos para o comparativo final.
type ResultadoInstancia struct {
	Nome      string
	Otimo     float64
	VMP_2opt  float64
	IMB_OrOpt float64
	AG        float64
	Memetico  float64
}

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

	resultados := make([]ResultadoInstancia, 0, len(instancias))

	for _, in := range instancias {
		res := processaInstancia(in, saidas)
		resultados = append(resultados, res)
	}

	// Parte 4: gerar outputs_u3/COMPARATIVO_GERAL.txt
	gerarComparativoGeral(resultados, saidas)

	fmt.Println("\nConcluido. Saidas em:", saidas)
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

func processaInstancia(in *pcv.Instancia, saidas string) ResultadoInstancia {
	fmt.Printf("[%s] processando (%s, N=%d)...\n", in.Nome, in.Medida, in.N)
	otimo := pcv.ValoresOtimos[in.Nome]

	// 1. VMP + 2-opt
	vmp := pcv.VizinhoMaisProximo{}.Constroi(in)
	depoisVMP := pcv.DoisOpt{}.Aplica(vmp, in)
	fmt.Printf("   -> VMP+2opt: custo=%.2f gap=%.2f%%\n", depoisVMP.Custo, pcv.Gap(depoisVMP.Custo, otimo))

	// 2. IMB + Or-opt
	imb := pcv.InsercaoMaisBarata{}.Constroi(in)
	depoisIMB := pcv.OrOpt{}.Aplica(imb, in)
	fmt.Printf("   -> IMB+OrOpt: custo=%.2f gap=%.2f%%\n", depoisIMB.Custo, pcv.Gap(depoisIMB.Custo, otimo))

	// 3. Algoritmo Genético
	ag := pcv.AlgoritmoGenetico{Par: pcv.ParametrosPadrao()}
	resumoAG := pcv.ExecutaExperimento(ag, in, execucoesEstocasticas, sementeBase)
	caminhoAG := filepath.Join(saidas, "RESUMO_AG_"+in.Nome+".txt")
	if err := os.WriteFile(caminhoAG, []byte(resumoAG.TextoResumo(in)), 0o644); err != nil {
		fmt.Println("Erro ao escrever resumo AG:", err)
		os.Exit(1)
	}
	fmt.Printf("   -> AG: menor=%.2f gap=%.2f%% media=%.2f\n", resumoAG.Melhor, pcv.Gap(resumoAG.Melhor, otimo), resumoAG.Media)

	// 4. Algoritmo Memético
	memetico := pcv.AlgoritmoMemetico{Par: pcv.ParametrosMemeticoPadrao()}
	resumoMemetico := pcv.ExecutaExperimento(memetico, in, execucoesEstocasticas, sementeBase)
	caminhoMemetico := filepath.Join(saidas, "RESUMO_MEMETICO_"+in.Nome+".txt")
	if err := os.WriteFile(caminhoMemetico, []byte(resumoMemetico.TextoResumo(in)), 0o644); err != nil {
		fmt.Println("Erro ao escrever resumo Memetico:", err)
		os.Exit(1)
	}
	fmt.Printf("   -> Memetico: menor=%.2f gap=%.2f%% media=%.2f\n", resumoMemetico.Melhor, pcv.Gap(resumoMemetico.Melhor, otimo), resumoMemetico.Media)
	fmt.Println(strings.Repeat("-", 50))

	return ResultadoInstancia{
		Nome:      in.Nome,
		Otimo:     otimo,
		VMP_2opt:  depoisVMP.Custo,
		IMB_OrOpt: depoisIMB.Custo,
		AG:        resumoAG.Melhor,
		Memetico:  resumoMemetico.Melhor,
	}
}

func gerarComparativoGeral(resultados []ResultadoInstancia, saidas string) {
	var sb strings.Builder
	_, _ = sb.WriteString("==========================================================================================\n")
	_, _ = sb.WriteString("                              COMPARATIVO GERAL DOS ALGORITMOS\n")
	_, _ = sb.WriteString("==========================================================================================\n\n")

	_, _ = sb.WriteString(fmt.Sprintf("%-12s | %10s | %10s | %10s | %10s | %10s\n",
		"Instancia", "Otimo", "VMP+2opt", "IMB+OrOpt", "AG(Menor)", "Memetico"))
	_, _ = sb.WriteString(strings.Repeat("-", 90) + "\n")

	for _, r := range resultados {
		// Formata os custos e gaps
		sOtimo := fmt.Sprintf("%.2f", r.Otimo)
		sVMP := fmt.Sprintf("%.2f (%.1f%%)", r.VMP_2opt, pcv.Gap(r.VMP_2opt, r.Otimo))
		sIMB := fmt.Sprintf("%.2f (%.1f%%)", r.IMB_OrOpt, pcv.Gap(r.IMB_OrOpt, r.Otimo))
		sAG := fmt.Sprintf("%.2f (%.1f%%)", r.AG, pcv.Gap(r.AG, r.Otimo))
		sMem := fmt.Sprintf("%.2f (%.1f%%)", r.Memetico, pcv.Gap(r.Memetico, r.Otimo))

		_, _ = sb.WriteString(fmt.Sprintf("%-12s | %10s | %14s | %14s | %14s | %14s\n",
			r.Nome, sOtimo, sVMP, sIMB, sAG, sMem))
	}

	_, _ = sb.WriteString(strings.Repeat("-", 90) + "\n")
	
	caminho := filepath.Join(saidas, "COMPARATIVO_GERAL.txt")
	if err := os.WriteFile(caminho, []byte(sb.String()), 0o644); err != nil {
		fmt.Println("Erro ao escrever COMPARATIVO_GERAL:", err)
		os.Exit(1)
	}
	fmt.Println("Arquivo COMPARATIVO_GERAL.txt gerado com sucesso!")
}
