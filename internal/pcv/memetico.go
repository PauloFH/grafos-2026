package pcv

import (
	"math/rand/v2"
	"sort"
)

func ParametrosMemeticoPadrao() ParametrosAG {
	return ParametrosAG{
		TamPopulacao:   50,
		Geracoes:       200,
		TaxaCruzamento: 0.9,
		TaxaMutacao:    0.10,
		TamTorneio:     3,
		Elitismo:       2,
		// TaxaEducacao=0.1 melhor que 0.2 nas instancias testadas
		TaxaEducacao: 0.1,
	}
}

const limiarDiversificacao = 50

var buscasLocaisMemetico = []BuscaLocal{DoisOpt{}, OrOpt{}, Swap{}}

type AlgoritmoMemetico struct{ Par ParametrosAG }

func (a AlgoritmoMemetico) Nome() string { return "Algoritmo Memetico" }

// Executa segue o fluxo (1) Inicio -> (2) Fitness -> [(3) Nova Geracao ->
// (4) Busca Local -> (5) Renovar -> (6) Teste]*.
func (a AlgoritmoMemetico) Executa(in *Instancia, semente int64) Rota {
	rng := rand.New(rand.NewPCG(uint64(semente), uint64(semente)))

	// (1) Inicio
	pop := PopulacaoInicial(a.Par.TamPopulacao, in, rng)

	// (2) Fitness: custo de cada individuo ja calculado por PopulacaoInicial/NovaRota
	melhor := pop[indiceMenorCusto(pop)].Clona()
	semMelhora := 0
	taxaMutacao := a.Par.TaxaMutacao

	for geracao := 1; geracao <= a.Par.Geracoes; geracao++ {
		// (3) Nova Geracao: Selecao (torneio + elitismo), Cruzamento (OX), Mutacao (inversao)
		nova := novaGeracaoMemetico(pop, a.Par, taxaMutacao, in, rng)

		// (4) Busca Local: educa parte da nova populacao com 2-opt, Or-opt ou Swap
		nova = educaPopulacao(nova, a.Par.TaxaEducacao, in, rng)

		// Elite recebe as 3 buscas locais em cadeia, garantindo
		// o maximo resultado nas solucoes mais promissoras
		nova = educaElite(nova, a.Par.Elitismo, in)

		// (5) Renovar
		pop = nova

		// (6) Teste: melhor global e parada por geracoes sem melhoria
		melhorGer := pop[indiceMenorCusto(pop)]
		if melhorGer.Custo < melhor.Custo-1e-9 {
			melhor = melhorGer.Clona()
			semMelhora = 0
			taxaMutacao = a.Par.TaxaMutacao // encontrou melhora, volta a mutacao ao normal
		} else {
			semMelhora++
			if semMelhora == limiarDiversificacao {
				// Para diversificar a populacao, antes de desistir, se estagnar, dobra a mutação
				taxaMutacao = min(1.0, a.Par.TaxaMutacao*2)
			}
		}
		if semMelhora >= 100 {
			break
		}
	}

	return melhor.Clona()
}

// novaGeracaoMemetico monta a etapa (3) Nova Geracao: os Par.Elitismo
// melhores da populacao atual passam intactos, e o restante vem de torneio
func novaGeracaoMemetico(pop []Rota, par ParametrosAG, taxaMutacao float64, in *Instancia, rng *rand.Rand) []Rota {
	idx := make([]int, len(pop))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(x, y int) bool {
		return pop[idx[x]].Custo < pop[idx[y]].Custo
	})

	nova := make([]Rota, 0, par.TamPopulacao)
	for e := 0; e < par.Elitismo && e < len(idx); e++ {
		nova = append(nova, pop[idx[e]].Clona())
	}

	for len(nova) < par.TamPopulacao {
		pai1 := SelecaoTorneio(pop, par.TamTorneio, rng)
		pai2 := SelecaoTorneio(pop, par.TamTorneio, rng)
		var filho Rota
		if rng.Float64() < par.TaxaCruzamento {
			filho = CruzamentoOX(pai1, pai2, in, rng)
		} else {
			filho = pai1.Clona()
		}
		filho = mutacaoInversao(filho, taxaMutacao, in, rng)
		nova = append(nova, filho)
	}
	return nova
}

// educaElite aplica as 3 buscas locais em cadeia (2-opt -> Or-opt -> Swap)
// aos n primeiros individuos de pop - a elite
func educaElite(pop []Rota, n int, in *Instancia) []Rota {
	for i := 0; i < n && i < len(pop); i++ {
		r := pop[i]
		for _, busca := range buscasLocaisMemetico {
			r = busca.Aplica(r, in)
		}
		pop[i] = r
	}
	return pop
}

// educaPopulacao e a etapa (4) Busca Local: com probabilidade taxa,
// sorteia uma das buscas locais (2-opt, Or-opt, Swap) e a aplica
// ate o otimo local daquela vizinhanca
func educaPopulacao(pop []Rota, taxa float64, in *Instancia, rng *rand.Rand) []Rota {
	if taxa <= 0 {
		return pop
	}
	for i := range pop {
		if rng.Float64() >= taxa {
			continue
		}
		busca := buscasLocaisMemetico[rng.IntN(len(buscasLocaisMemetico))]
		pop[i] = busca.Aplica(pop[i], in)
	}
	return pop
}
