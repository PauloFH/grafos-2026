package relatorio

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/pcv"
)

// rotaFechadaPCV converte uma rota de posicoes para ids originais de cidade e devolve a sequencia unida por " -> ",
// com o id do deposito repetido no fim para explicitar o fechamento do ciclo (ex.: "1 -> 5 -> ... -> 1").
func rotaFechadaPCV(in *pcv.Instancia, r pcv.Rota) string {
	ids := r.IdsOriginais(in)
	partes := make([]string, 0, len(ids)+1)
	for _, id := range ids {
		partes = append(partes, strconv.Itoa(id))
	}
	if len(ids) > 0 {
		partes = append(partes, strconv.Itoa(ids[0]))
	}
	return strings.Join(partes, " -> ")
}


func FormataInstanciaPCV(in *pcv.Instancia) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "  Medida : %s\n", in.Medida)
	fmt.Fprintf(&sb, "  Cidades: %d\n", in.N)
	for k := 0; k < in.N; k++ {
		fmt.Fprintf(&sb, "    %2d %s\n", in.Cidades[k], in.Nomes[k])
	}

	return sb.String()
}

// FormataRotaPCV gera o texto de uma secao de heuristica construtiva,mostrando a rota e o custo ANTES e DEPOIS da busca local,
// e a melhora percentual obtida.
func FormataRotaPCV(in *pcv.Instancia, antes, depois pcv.Rota, nomeBusca string) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "  %-12s: %s\n", "Busca local", nomeBusca)
	fmt.Fprintf(&sb, "  %-12s: %s\n", "Rota antes", rotaFechadaPCV(in, antes))
	fmt.Fprintf(&sb, "  %-12s: %.2f\n", "Custo antes", antes.Custo)
	fmt.Fprintf(&sb, "  %-12s: %s\n", "Rota depois", rotaFechadaPCV(in, depois))
	fmt.Fprintf(&sb, "  %-12s: %.2f\n", "Custo depois", depois.Custo)

	// melhora percentual da busca local sobre o custo construido
	melhora := 0.0
	if antes.Custo > 0 {
		melhora = (antes.Custo - depois.Custo) / antes.Custo * 100
	}
	fmt.Fprintf(&sb, "  %-12s: %.2f%%\n", "Melhora", melhora)

	return sb.String()
}

// FormataResumoPCV gera o texto da secao de um experimento estocastico
func FormataResumoPCV(res pcv.ResumoExperimento, in *pcv.Instancia) string {
	var sb strings.Builder

	otimo := pcv.ValoresOtimos[res.Instancia] // 0.00 se ausente do mapa

	fmt.Fprintf(&sb, "  %-12s: %d\n", "Execucoes", res.Execucoes)
	fmt.Fprintf(&sb, "  %-12s: %.2f\n", "Menor valor", res.Melhor)
	fmt.Fprintf(&sb, "  %-12s: %.2f\n", "Valor medio", res.Media)
	fmt.Fprintf(&sb, "  %-12s: %.3fs\n", "Tempo medio", res.TempoMedio.Seconds())
	fmt.Fprintf(&sb, "  %-12s: %.2f%%\n", "Gap do menor", pcv.Gap(res.Melhor, otimo))
	fmt.Fprintf(&sb, "  %-12s: %s\n", "Melhor rota", rotaFechadaPCV(in, res.MelhorRota))

	return sb.String()
}

// LinhaComparativo agrega, para uma instancia, o melhor custo obtido por
// cada um dos 4 metodos, para a tabela COMPARATIVO_GERAL.
type LinhaComparativo struct {
	Instancia string  // ex. "PROBLEMA_01"
	Otimo     float64 // melhor valor conhecido
	VMP       float64 // Vizinho mais Proximo + 2-opt
	IMB       float64 // Insercao mais Barata + Or-opt
	AG        float64 // menor das 20 execucoes do AG
	Memetico  float64 // menor das 20 execucoes do memetico
}


func FormataComparativoPCV(linhas []LinhaComparativo) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%-12s| %8s | %8s | %6s | %9s | %6s | %8s | %6s | %8s | %6s\n",
		"Instancia", "Otimo", "VMP+2opt", "gap%", "IMB+OrOpt", "gap%", "AG", "gap%", "Memetico", "gap%")
	sb.WriteString(strings.Repeat("-", 12) + "+" + strings.Repeat("-", 10) + "+" +
		strings.Repeat("-", 10) + "+" + strings.Repeat("-", 8) + "+" +
		strings.Repeat("-", 11) + "+" + strings.Repeat("-", 8) + "+" +
		strings.Repeat("-", 10) + "+" + strings.Repeat("-", 8) + "+" +
		strings.Repeat("-", 10) + "+" + strings.Repeat("-", 7) + "\n")

	for _, l := range linhas {
		// gap de cada metodo em relacao ao melhor valor conhecido da linha
		gap := func(custo float64) string {
			return fmt.Sprintf("%.2f%%", pcv.Gap(custo, l.Otimo))
		}
		fmt.Fprintf(&sb, "%-12s| %8.2f | %8.2f | %6s | %9.2f | %6s | %8.2f | %6s | %8.2f | %6s\n",
			l.Instancia, l.Otimo,
			l.VMP, gap(l.VMP),
			l.IMB, gap(l.IMB),
			l.AG, gap(l.AG),
			l.Memetico, gap(l.Memetico))
	}

	return sb.String()
}
