package pcv

import (
	"math/rand/v2"
	"testing"
)

func carregaInstanciaTeste(t *testing.T, arquivo string) *Instancia {
	t.Helper()
	in, err := CarregaInstancia("../../inputs_u3/" + arquivo)
	if err != nil {
		t.Fatalf("CarregaInstancia(%s): %v", arquivo, err)
	}
	return in
}

func rngTeste(semente int64) *rand.Rand {
	return rand.New(rand.NewPCG(uint64(semente), uint64(semente)))
}

func TestPopulacaoInicial(t *testing.T) {
	casos := []struct {
		arquivo string
		tam     int
	}{
		{"PROBLEMA_11_KM_06.txt", 30},
		{"PROBLEMA_07_KM_12.txt", 50},
		{"PROBLEMA_01_KM_48.txt", 20},
	}
	for _, c := range casos {
		t.Run(c.arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, c.arquivo)
			pop := PopulacaoInicial(c.tam, in, rngTeste(42))
			if len(pop) != c.tam {
				t.Fatalf("len(pop) = %d, esperado %d", len(pop), c.tam)
			}
			for i, r := range pop {
				if err := ValidaRota(r.Ordem, in); err != nil {
					t.Fatalf("rota %d invalida: %v", i, err)
				}
				if got := CustoRota(r.Ordem, in); got != r.Custo {
					t.Fatalf("rota %d: custo %f inconsistente com ordem (%f)", i, r.Custo, got)
				}
			}
		})
	}
}

func TestSelecaoTorneio(t *testing.T) {
	in := carregaInstanciaTeste(t, "PROBLEMA_07_KM_12.txt")
	rng := rngTeste(7)
	pop := PopulacaoInicial(10, in, rng)
	custos := make(map[float64]bool, len(pop))
	for _, r := range pop {
		custos[r.Custo] = true
	}
	for rodada := 0; rodada < 50; rodada++ {
		venc := SelecaoTorneio(pop, 3, rng)
		if !custos[venc.Custo] {
			t.Fatalf("rodada %d: custo %f do vencedor nao pertence a populacao", rodada, venc.Custo)
		}
		// mutacao no clone nao pode vazar para a populacao
		original := venc.Ordem[1]
		venc.Ordem[1] = -999
		for i, r := range pop {
			if r.Ordem[1] == -999 {
				t.Fatalf("rodada %d: retorno compartilha Ordem com pop[%d]", rodada, i)
			}
		}
		venc.Ordem[1] = original
	}
}

func TestCruzamentoOXPermutacaoValida(t *testing.T) {
	casos := []struct {
		arquivo string
		semente int64
	}{
		{"PROBLEMA_11_KM_06.txt", 1},
		{"PROBLEMA_11_KM_06.txt", 42},
		{"PROBLEMA_07_KM_12.txt", 7},
		{"PROBLEMA_01_KM_48.txt", 42},
	}
	for _, c := range casos {
		t.Run(c.arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, c.arquivo)
			rng := rngTeste(c.semente)
			pop := PopulacaoInicial(20, in, rng)
			for rodada := 0; rodada < 200; rodada++ {
				pai1 := pop[rng.IntN(len(pop))]
				pai2 := pop[rng.IntN(len(pop))]
				copia1 := pai1.Clona()
				copia2 := pai2.Clona()

				filho := CruzamentoOX(pai1, pai2, in, rng)
				if filho.Ordem[0] != 0 {
					t.Fatalf("rodada %d: filho.Ordem[0] = %d, esperado 0", rodada, filho.Ordem[0])
				}
				if err := ValidaRota(filho.Ordem, in); err != nil {
					t.Fatalf("rodada %d: filho invalido: %v", rodada, err)
				}
				if got := CustoRota(filho.Ordem, in); got != filho.Custo {
					t.Fatalf("rodada %d: custo do filho inconsistente", rodada)
				}
				for k := range pai1.Ordem {
					if pai1.Ordem[k] != copia1.Ordem[k] || pai2.Ordem[k] != copia2.Ordem[k] {
						t.Fatalf("rodada %d: CruzamentoOX modificou um dos pais", rodada)
					}
				}
			}
		})
	}
}

func TestMutacaoTrocaNaoModificaEntrada(t *testing.T) {
	casos := []struct {
		nome string
		taxa float64
	}{
		{"taxa zero", 0.0},
		{"taxa um", 1.0},
		{"taxa canonica", 0.10},
	}
	in := carregaInstanciaTeste(t, "PROBLEMA_07_KM_12.txt")
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			rng := rngTeste(42)
			r := PopulacaoInicial(1, in, rng)[0]
			copia := r.Clona()
			for rodada := 0; rodada < 100; rodada++ {
				filho := MutacaoTroca(r, c.taxa, in, rng)
				for k := range r.Ordem {
					if r.Ordem[k] != copia.Ordem[k] {
						t.Fatalf("rodada %d: MutacaoTroca modificou a entrada", rodada)
					}
				}
				if err := ValidaRota(filho.Ordem, in); err != nil {
					t.Fatalf("rodada %d: saida invalida: %v", rodada, err)
				}
				if got := CustoRota(filho.Ordem, in); got != filho.Custo {
					t.Fatalf("rodada %d: custo da saida inconsistente", rodada)
				}
				if c.taxa == 0.0 && filho.Custo != r.Custo {
					t.Fatalf("rodada %d: taxa 0 alterou a rota", rodada)
				}
			}
		})
	}
}
