// Comando unidade3 roda a Unidade 3 (PCV): le inputs_u3/*.txt, roda os 4
// metodos (VMP+2opt, IMB+OrOpt, AG e Memetico) e grava relatorios, resumos e
// o comparativo geral em outputs_u3/. Este esqueleto do scaffolding marca nos
// TODOs os pontos de integracao das Partes 1 a 4.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/pcv"
)

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

// imprimeTabelaInstancias imprime nome, medida, N e melhor valor conhecido
// de cada instancia ("-" quando ausente de ValoresOtimos).
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

// processaInstancia roda os 4 metodos sobre uma instancia e grava os
// arquivos em outputs_u3/; enquanto as Partes 1 a 4 nao integrarem, apenas
// anuncia a instancia (chamadas canonicas nos TODOs abaixo).
func processaInstancia(in *pcv.Instancia, saidas string) {
	fmt.Printf("[%s] processando (%s, N=%d)...\n", in.Nome, in.Medida, in.N)
	_ = saidas // usado quando as Partes 1-4 integrarem as chamadas abaixo

	// TODO(Parte 1): VMP + 2-opt (deterministico, 1 execucao):
	//   antes := pcv.VizinhoMaisProximo{}.Constroi(in)
	//   depois := pcv.DoisOpt{}.Aplica(antes, in)

	// TODO(Parte 2): IMB + Or-opt (deterministico, 1 execucao):
	//   antes := pcv.InsercaoMaisBarata{}.Constroi(in)
	//   depois := pcv.OrOpt{}.Aplica(antes, in)

	// TODO(Parte 3): AG, 20 execucoes com sementes sementeBase+i:
	//   ag := pcv.AlgoritmoGenetico{Par: pcv.ParametrosPadrao()}
	//   resumoAG := pcv.ExecutaExperimento(ag, in, 20, sementeBase)
	//   gravar outputs_u3/RESUMO_AG_<in.Nome>.txt com resumoAG.TextoResumo(in)

	// TODO(Parte 4): Memetico (pop=50, ger=200, Buscas={DoisOpt, OrOpt,
	// Swap}), 20 execucoes; gravar outputs_u3/RESUMO_MEMETICO_<in.Nome>.txt.

	// TODO(Parte 4): relatorio por instancia (outputs_u3/<in.Nome>.txt) com
	// as secoes INSTANCIA, VIZINHO_MAIS_PROXIMO, INSERCAO_MAIS_BARATA,
	// ALGORITMO_GENETICO e ALGORITMO_MEMETICO.
}
