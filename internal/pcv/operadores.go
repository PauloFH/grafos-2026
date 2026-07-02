package pcv

import "math/rand/v2"

// ParametrosAG agrupa os parametros dos algoritmos genetico e memetico.
// Os valores canonicos do trabalho estao em ParametrosPadrao.
type ParametrosAG struct {
	TamPopulacao   int
	Geracoes       int
	TaxaCruzamento float64 // probabilidade de aplicar o OX a um par de pais
	TaxaMutacao    float64 // probabilidade de mutacao POR INDIVIDUO
	TamTorneio     int     // k da selecao por torneio
	Elitismo       int     // melhores copiados intactos a cada geracao
	TaxaEducacao   float64 // probabilidade de educar um filho (so o memetico usa)
}

// ParametrosPadrao devolve os parametros canonicos do trabalho; o memetico
// (Parte 4) sobrescreve TamPopulacao=50 e Geracoes=200.
func ParametrosPadrao() ParametrosAG {
	return ParametrosAG{
		TamPopulacao:   100,
		Geracoes:       500,
		TaxaCruzamento: 0.9,
		TaxaMutacao:    0.10,
		TamTorneio:     3,
		Elitismo:       2,
		TaxaEducacao:   0.2,
	}
}

// PopulacaoInicial gera tam rotas aleatorias embaralhando as posicoes 1..N-1
// com Fisher-Yates; o deposito permanece fixo na posicao 0.
func PopulacaoInicial(tam int, in *Instancia, rng *rand.Rand) []Rota {
	panic("pcv: nao implementado - Parte 3 (genetico)")
}

// SelecaoTorneio sorteia k candidatos com reposicao e devolve o de menor
// custo. O retorno compartilha Ordem com a populacao, o que e seguro porque
// nenhum operador muta suas entradas.
func SelecaoTorneio(pop []Rota, k int, rng *rand.Rand) Rota {
	panic("pcv: nao implementado - Parte 3 (genetico)")
}

// CruzamentoOX aplica o Order Crossover ao cromossomo efetivo Ordem[1:]
// (deposito preservado na posicao 0) e produz UM filho valido com custo
// calculado. Os pais nao sao modificados.
func CruzamentoOX(pai1, pai2 Rota, in *Instancia, rng *rand.Rand) Rota {
	panic("pcv: nao implementado - Parte 3 (genetico)")
}

// MutacaoTroca troca duas posicoes distintas de Ordem[1:] com probabilidade
// taxa (por individuo). A entrada nunca e modificada.
func MutacaoTroca(r Rota, taxa float64, in *Instancia, rng *rand.Rand) Rota {
	panic("pcv: nao implementado - Parte 3 (genetico)")
}
