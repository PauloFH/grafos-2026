package pcv

import (
	"math/rand/v2"
	"sort"
)

// AlgoritmoGenetico implementa a interface Estocastico com um AG geracional
type AlgoritmoGenetico struct{ Par ParametrosAG }

// Nome devolve o nome canonico do algoritmo usado nos relatorios.
func (a AlgoritmoGenetico) Nome() string { return "Algoritmo Genetico" }

// Executa roda uma execucao completa do AG sobre a instancia com a semente dada.
func (a AlgoritmoGenetico) Executa(in *Instancia, semente int64) Rota {
	// rng canonico da U3 - unico ponto de aleatoriedade de toda a execucao
	rng := rand.New(rand.NewPCG(uint64(semente), uint64(semente)))

	pop := PopulacaoInicial(a.Par.TamPopulacao, in, rng)
	melhor := pop[indiceMenorCusto(pop)].Clona()
	semMelhora := 0

	for geracao := 1; geracao <= a.Par.Geracoes; geracao++ {
		//Os Par.Elitismo melhores passam intactos.
		//Empates de custo preservam o indice menor, garantindo determinismo.
		idx := make([]int, len(pop))
		for i := range idx {
			idx[i] = i
		}
		sort.SliceStable(idx, func(x, y int) bool {
			return pop[idx[x]].Custo < pop[idx[y]].Custo
		})
		nova := make([]Rota, 0, a.Par.TamPopulacao)
		for e := 0; e < a.Par.Elitismo && e < len(idx); e++ {
			nova = append(nova, pop[idx[e]].Clona())
		}

		// Torneio sempre sobre a populacao antiga
		for len(nova) < a.Par.TamPopulacao {
			pai1 := SelecaoTorneio(pop, a.Par.TamTorneio, rng)
			pai2 := SelecaoTorneio(pop, a.Par.TamTorneio, rng)
			var filho Rota
			if rng.Float64() < a.Par.TaxaCruzamento {
				filho = CruzamentoOX(pai1, pai2, in, rng)
			} else {
				filho = pai1.Clona()
			}
			filho = mutacaoInversao(filho, a.Par.TaxaMutacao, in, rng)
			nova = append(nova, filho)
		}

		pop = nova

		// atualizacao do melhor global;
		// Melhora estrita para nao zerar o contador de estagnacao em empates
		melhorGer := pop[indiceMenorCusto(pop)]
		if melhorGer.Custo < melhor.Custo-1e-9 {
			melhor = melhorGer.Clona()
			semMelhora = 0
		} else {
			semMelhora++
		}
		if semMelhora >= 100 {
			break // parada antecipada por estagnacao
		}
	}

	return melhor.Clona()
}

// mutacaoInversao inverte um segmento de Ordem[1:] com probabilidade taxa
// (POR INDIVIDUO), sorteando duas posicoes distintas i e j em [1, N-1] e
// revertendo o trecho entre elas.
func mutacaoInversao(r Rota, taxa float64, in *Instancia, rng *rand.Rand) Rota {
	filho := r.Clona() // nunca mutar a entrada
	if in.N <= 2 {
		return filho // menos de 2 genes moveis - nada a fazer
	}
	// A ordem dos sorteios faz parte do determinismo canonico da execucao
	if rng.Float64() >= taxa {
		return filho
	}

	// duas posicoes distintas em [1, N-1]; j e re-sorteado ate diferir de i
	i := 1 + rng.IntN(in.N-1)
	j := 1 + rng.IntN(in.N-1)
	for j == i {
		j = 1 + rng.IntN(in.N-1)
	}
	if i > j {
		i, j = j, i
	}

	// inverte o segmento [i, j] do cromossomo (posicao 0 = deposito, intocada)
	for i < j {
		filho.Ordem[i], filho.Ordem[j] = filho.Ordem[j], filho.Ordem[i]
		i++
		j--
	}
	filho.Custo = CustoRota(filho.Ordem, in)
	return filho
}

// indiceMenorCusto devolve o indice do individuo de menor custo da populacao
func indiceMenorCusto(pop []Rota) int {
	menor := 0
	for i := 1; i < len(pop); i++ {
		if pop[i].Custo < pop[menor].Custo {
			menor = i
		}
	}
	return menor
}
