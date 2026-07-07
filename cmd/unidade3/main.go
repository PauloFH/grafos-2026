package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/pcv"
	"github.com/PauloFH/grafos-2026/internal/relatorio"
	"github.com/PauloFH/grafos-2026/internal/web"
)

// sementeBase e a semente da primeira execucao dos metodos estocasticos;
const sementeBase int64 = 42

// execucoesEstocasticas e o numero de execucoes independentes de cada metaheuristica por instancia
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

	// imagens PNG das rotas sao um extra: exigem graphviz
	gerarImagens := relatorio.GraphvizDisponivel()
	if !gerarImagens {
		fmt.Println("Aviso: graphviz (neato) nao encontrado - imagens das rotas nao serao geradas.")
		fmt.Println()
	}

	// acumula uma linha por instancia para o comparativo geral
	linhas := make([]relatorio.LinhaComparativo, 0, len(instancias))
	for _, in := range instancias {
		linhas = append(linhas, processaInstancia(in, saidas, gerarImagens))
	}

	// COMPARATIVO_GERAL.txt: tabela final com o melhor custo de cada metodo
	// (AG e Memetico entram com o MENOR das 20 execucoes) e o gap vs otimo.
	caminhoComp := filepath.Join(saidas, "COMPARATIVO_GERAL.txt")
	if err := os.WriteFile(caminhoComp, []byte(relatorio.FormataComparativoPCV(linhas)), 0o644); err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}
	fmt.Println("Comparativo geral ->", caminhoComp)

	fmt.Println("Concluido. Saidas em:", saidas)
}

// imprimeTabelaInstancias imprime nome, medida, N e melhor valor conhecido
// de cada instancia
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

// processaInstancia roda os 4 metodos sobre uma instancia, grava os arquivos
// de saida em outputs_u3/ (RESUMO_AG, RESUMO_MEMETICO, o relatorio por
// instancia PROBLEMA_XX.txt e, se gerarImagens, os PNGs das rotas) e devolve
// a linha do comparativo geral.
func processaInstancia(in *pcv.Instancia, saidas string, gerarImagens bool) relatorio.LinhaComparativo {
	fmt.Printf("[%s] processando (%s, N=%d)...\n", in.Nome, in.Medida, in.N)
	otimo := pcv.ValoresOtimos[in.Nome]

	// (1) Vizinho mais Proximo + 2-opt (deterministico, 1 execucao)
	vmp := pcv.VizinhoMaisProximo{}.Constroi(in)
	vmp2opt := pcv.DoisOpt{}.Aplica(vmp, in)
	fmt.Printf("[%s] VMP+2opt: %.2f -> %.2f (gap %.2f%%)\n",
		in.Nome, vmp.Custo, vmp2opt.Custo, pcv.Gap(vmp2opt.Custo, otimo))

	// (2) Insercao mais Barata + Or-opt (deterministico, 1 execucao)
	imb := pcv.InsercaoMaisBarata{}.Constroi(in)
	imbOrOpt := pcv.OrOpt{}.Aplica(imb, in)
	fmt.Printf("[%s] IMB+OrOpt: %.2f -> %.2f (gap %.2f%%)\n",
		in.Nome, imb.Custo, imbOrOpt.Custo, pcv.Gap(imbOrOpt.Custo, otimo))

	// (3) Algoritmo Genetico: 20 execucoes com sementes sementeBase+i
	ag := pcv.AlgoritmoGenetico{Par: pcv.ParametrosPadrao()}
	resumoAG := pcv.ExecutaExperimento(ag, in, execucoesEstocasticas, sementeBase)
	gravaArquivo(filepath.Join(saidas, "RESUMO_AG_"+in.Nome+".txt"), resumoAG.TextoResumo(in))
	fmt.Printf("[%s] AG: menor=%.2f media=%.2f\n", in.Nome, resumoAG.Melhor, resumoAG.Media)

	// (4) Algoritmo Memetico: 20 execucoes com sementes sementeBase+i
	memetico := pcv.AlgoritmoMemetico{Par: pcv.ParametrosMemeticoPadrao()}
	resumoMem := pcv.ExecutaExperimento(memetico, in, execucoesEstocasticas, sementeBase)
	gravaArquivo(filepath.Join(saidas, "RESUMO_MEMETICO_"+in.Nome+".txt"), resumoMem.TextoResumo(in))
	fmt.Printf("[%s] Memetico: menor=%.2f media=%.2f\n", in.Nome, resumoMem.Melhor, resumoMem.Media)

	// Relatorio por instancia com as 5 secoes (instancia + 4 metodos)
	gravaArquivo(filepath.Join(saidas, in.Nome+".txt"),
		relatorioInstancia(in, vmp, vmp2opt, imb, imbOrOpt, resumoAG, resumoMem))

	// Imagens PNG das rotas resultantes de cada metodo (uma por metodo)
	if gerarImagens {
		geraImagensRota(in, saidas, vmp2opt, imbOrOpt, resumoAG.MelhorRota, resumoMem.MelhorRota)
	}

	return relatorio.LinhaComparativo{
		Instancia: in.Nome,
		Otimo:     otimo,
		VMP:       vmp2opt.Custo,
		IMB:       imbOrOpt.Custo,
		AG:        resumoAG.Melhor,
		Memetico:  resumoMem.Melhor,
	}
}

// relatorioInstancia monta o texto do relatorio por instancia com as secoes
// INSTANCIA, VIZINHO_MAIS_PROXIMO (antes/depois do 2-opt), INSERCAO_MAIS_BARATA
// (antes/depois do Or-opt), ALGORITMO_GENETICO e ALGORITMO_MEMETICO.
func relatorioInstancia(in *pcv.Instancia, vmp, vmp2opt, imb, imbOrOpt pcv.Rota, resumoAG, resumoMem pcv.ResumoExperimento) string {
	var sb strings.Builder
	sep := strings.Repeat("=", 50)
	secao := func(titulo, corpo string) {
		fmt.Fprintf(&sb, "[%s]\n%s\n", titulo, corpo)
	}

	fmt.Fprintf(&sb, "%s\n RELATORIO - %s\n%s\n\n", sep, in.Nome, sep)
	secao("INSTANCIA", relatorio.FormataInstanciaPCV(in))
	secao("VIZINHO_MAIS_PROXIMO", relatorio.FormataRotaPCV(in, vmp, vmp2opt, "2-opt"))
	secao("INSERCAO_MAIS_BARATA", relatorio.FormataRotaPCV(in, imb, imbOrOpt, "Or-opt"))
	secao("ALGORITMO_GENETICO", relatorio.FormataResumoPCV(resumoAG, in))
	secao("ALGORITMO_MEMETICO", relatorio.FormataResumoPCV(resumoMem, in))
	return sb.String()
}

// geraImagensRota grava um PNG por metodo com a rota resultante desenhada
// sobre as cidades, usando as posicoes 2D do layout MDS (mesmas do front).
func geraImagensRota(in *pcv.Instancia, saidas string, vmp2opt, imbOrOpt, ag, memetico pcv.Rota) {
	// posicoes 2D do MDS (normalizadas em [0.05, 0.95]), indexadas por posicao
	pontos := make([]relatorio.PontoPlano, in.N)
	for pos, dto := range web.PosicoesMDS(in) {
		pontos[pos] = relatorio.PontoPlano{X: dto.X, Y: dto.Y, Id: dto.ID, Nome: dto.Nome}
	}

	tituloGrafo := fmt.Sprintf("%s - grafo completo (%d cidades)", in.Nome, in.N)
	if err := relatorio.GerarPNGGrafoPCV(in, pontos, tituloGrafo, in.Nome+"_GRAFO", saidas); err != nil {
		fmt.Printf("[%s] aviso: falha ao gerar %s_GRAFO.png: %v\n", in.Nome, in.Nome, err)
	}

	imgs := []struct {
		rota   pcv.Rota
		metodo string
	}{
		{vmp2opt, "VMP_2OPT"},
		{imbOrOpt, "IMB_OROPT"},
		{ag, "AG"},
		{memetico, "MEMETICO"},
	}
	for _, img := range imgs {
		titulo := fmt.Sprintf("%s - %s - custo %.2f", in.Nome, img.metodo, img.rota.Custo)
		nome := in.Nome + "_" + img.metodo
		if err := relatorio.GerarPNGRotaPCV(in, pontos, img.rota, titulo, nome, saidas); err != nil {
			fmt.Printf("[%s] aviso: falha ao gerar %s.png: %v\n", in.Nome, nome, err)
		}
	}
}

// gravaArquivo escreve conteudo em caminho e aborta o programa em caso de
// erro de escrita (o make u3 nao deve seguir com saidas parciais).
func gravaArquivo(caminho, conteudo string) {
	if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}
}
