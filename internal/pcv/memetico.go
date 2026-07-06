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
		TaxaEducacao:   0.2,
	}
}

var buscasLocaisMemetico = []BuscaLocal{DoisOpt{}, OrOpt{}, Swap{}}

type AlgoritmoMemetico struct{ Par ParametrosAG }

func (a AlgoritmoMemetico) Nome() string { return "Algoritmo Memetico" }

func (a AlgoritmoMemetico) Executa(in *Instancia, semente int64) Rota {
	rng := rand.New(rand.NewPCG(uint64(semente), uint64(semente)))

	pop := PopulacaoInicial(a.Par.TamPopulacao, in, rng)
	pop = educaPopulacao(pop, a.Par.TaxaEducacao, in, rng)
	melhor := pop[indiceMenorCusto(pop)].Clona()
	semMelhora := 0

	for geracao := 1; geracao <= a.Par.Geracoes; geracao++ {
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

		nova = educaPopulacao(nova, a.Par.TaxaEducacao, in, rng)
		pop = nova

		melhorGer := pop[indiceMenorCusto(pop)]
		if melhorGer.Custo < melhor.Custo-1e-9 {
			melhor = melhorGer.Clona()
			semMelhora = 0
		} else {
			semMelhora++
		}
		if semMelhora >= 100 {
			break
		}
	}

	return melhor.Clona()
}

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
