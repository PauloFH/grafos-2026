package pcv

import "math/rand/v2"

// ParametrosAG agrupa os parametros dos algoritmos genetico e memetico.
type ParametrosAG struct {
	TamPopulacao   int
	Geracoes       int
	TaxaCruzamento float64 // probabilidade de aplicar o OX a um par de pais
	TaxaMutacao    float64 // probabilidade de mutacao POR INDIVIDUO
	TamTorneio     int     // k da selecao por torneio
	Elitismo       int     // melhores copiados intactos a cada geracao
	TaxaEducacao   float64
}

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

// PopulacaoInicial gera tam rotas aleatorias embaralhando as posicoes 1..N-1 com Fisher-Yates; o deposito permanece fixo na posicao 0.
func PopulacaoInicial(tam int, in *Instancia, rng *rand.Rand) []Rota {
	pop := make([]Rota, 0, tam)
	for t := 0; t < tam; t++ {
		// comeca da identidade [0, 1, ..., N-1]
		ordem := make([]int, in.N)
		for i := range ordem {
			ordem[i] = i
		}
		// Fisher-Yates restrito ao trecho 1..N-1 
		// posicao 0 = deposito, intocada
		for i := in.N - 1; i >= 2; i-- {
			j := 1 + rng.IntN(i)
			ordem[i], ordem[j] = ordem[j], ordem[i]
		}
		pop = append(pop, NovaRota(ordem, in))
	}
	return pop
}

// SelecaoTorneio seleciona um individuo por torneio de k com reposicao.
func SelecaoTorneio(pop []Rota, k int, rng *rand.Rand) Rota {
	melhor := pop[rng.IntN(len(pop))]
	for t := 1; t < k; t++ {
		cand := pop[rng.IntN(len(pop))]
		if cand.Custo < melhor.Custo {
			melhor = cand
		}
	}
	return melhor.Clona()
}

// CruzamentoOX aplica o Order Crossover ao cromossomo efetivo Ordem[1:] e produz UM filho.
func CruzamentoOX(pai1, pai2 Rota, in *Instancia, rng *rand.Rand) Rota {
	m := in.N - 1
	if m < 3 {
		// sem espaco util para recorte: devolve copia do pai1
		return pai1.Clona()
	}

	// cromossomos efetivos
	c1 := pai1.Ordem[1:]
	c2 := pai2.Ordem[1:]

	// Sorteia o segmento de recorte [a, b] dentro de 0..M-1 (a == b e valido)
	a := rng.IntN(m)
	b := rng.IntN(m)
	if a > b {
		a, b = b, a
	}

	//Copia o segmento do pai1 para o filho, nas mesmas posicoes
	filho := make([]int, m)
	usado := make([]bool, in.N)
	for k := a; k <= b; k++ {
		filho[k] = c1[k]
		usado[c1[k]] = true
	}

	// Percorre o pai2 a partir de b+1 escrevendo nas posicoes livres do filho tambem a partir de b+1
	pos := (b + 1) % m
	for t := 0; t < m; t++ {
		g := c2[(b+1+t)%m]
		if usado[g] {
			continue // gene ja veio do pai1 - pula
		}
		filho[pos] = g
		usado[g] = true
		pos = (pos + 1) % m
	}

	// 4. reconstroi a rota completa com o deposito de volta na frente
	ordem := make([]int, 0, in.N)
	ordem = append(ordem, 0)
	ordem = append(ordem, filho...)
	return NovaRota(ordem, in)
}

// MutacaoTroca troca duas posicoes distintas de Ordem[1:] com probabilidade taxa, sorteando duas posicoes distintas i e j em [1, N-1] e trocando-as.
func MutacaoTroca(r Rota, taxa float64, in *Instancia, rng *rand.Rand) Rota {
	filho := r.Clona() // nunca mutar a entrada
	if in.N <= 2 {
		return filho // menos de 2 genes moveis - nada a fazer
	}
	// a moeda e sorteada ANTES das posicoes: a ordem dos sorteios faz parte
	// do determinismo canonico da execucao
	if rng.Float64() >= taxa {
		return filho
	}

	// duas posicoes distintas em [1, N-1]; j e re-sorteado ate diferir de i
	i := 1 + rng.IntN(in.N-1)
	j := 1 + rng.IntN(in.N-1)
	for j == i {
		j = 1 + rng.IntN(in.N-1)
	}

	filho.Ordem[i], filho.Ordem[j] = filho.Ordem[j], filho.Ordem[i]
	filho.Custo = CustoRota(filho.Ordem, in)
	return filho
}
