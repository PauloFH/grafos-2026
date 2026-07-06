package pcv

import (
	"math"
	"strings"
	"testing"
)

func TestMemeticoExecutaReprodutibilidade(t *testing.T) {
	casos := []struct {
		arquivo string
		semente int64
	}{
		{"PROBLEMA_11_KM_06.txt", 42},
		{"PROBLEMA_07_KM_12.txt", 42},
		{"PROBLEMA_07_KM_12.txt", 7},
	}
	alg := AlgoritmoMemetico{Par: ParametrosMemeticoPadrao()}
	for _, c := range casos {
		t.Run(c.arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, c.arquivo)
			r1 := alg.Executa(in, c.semente)
			r2 := alg.Executa(in, c.semente)
			if r1.Custo != r2.Custo {
				t.Fatalf("custos diferentes: %f != %f", r1.Custo, r2.Custo)
			}
			for k := range r1.Ordem {
				if r1.Ordem[k] != r2.Ordem[k] {
					t.Fatalf("Ordem difere na posicao %d: %d != %d", k, r1.Ordem[k], r2.Ordem[k])
				}
			}
		})
	}
}

func TestMemeticoExecutaRotaValida(t *testing.T) {
	casos := []string{
		"PROBLEMA_11_KM_06.txt",
		"PROBLEMA_09_KM_07.txt",
		"PROBLEMA_07_KM_12.txt",
	}
	alg := AlgoritmoMemetico{Par: ParametrosMemeticoPadrao()}
	for _, arquivo := range casos {
		t.Run(arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, arquivo)
			r := alg.Executa(in, 42)
			if err := ValidaRota(r.Ordem, in); err != nil {
				t.Fatalf("rota invalida: %v", err)
			}
			if got := CustoRota(r.Ordem, in); math.Abs(got-r.Custo) > 1e-9 {
				t.Fatalf("custo %f inconsistente com a ordem (%f)", r.Custo, got)
			}
		})
	}
}

func TestMemeticoAoMenosTaoBomQuantoAG(t *testing.T) {
	casos := []struct {
		arquivo string
		otimo   float64
	}{
		{"PROBLEMA_11_KM_06.txt", 344.90},
		{"PROBLEMA_12_MIN_06.txt", 305.00},
	}
	alg := AlgoritmoMemetico{Par: ParametrosMemeticoPadrao()}
	for _, c := range casos {
		t.Run(c.arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, c.arquivo)
			res := ExecutaExperimento(alg, in, 20, 42)
			if math.Abs(res.Melhor-c.otimo) > 1e-6 {
				t.Errorf("Melhor = %.2f, esperado o otimo %.2f", res.Melhor, c.otimo)
			}
		})
	}
}

func TestMemeticoDoGapP01(t *testing.T) {
	in := carregaInstanciaTeste(t, "PROBLEMA_01_KM_48.txt")
	alg := AlgoritmoMemetico{Par: ParametrosMemeticoPadrao()}
	res := ExecutaExperimento(alg, in, 20, 42)

	otimo := ValoresOtimos[in.Nome]
	if gap := Gap(res.Melhor, otimo); gap > 10.0 {
		t.Fatalf("gap do menor = %.2f%% (Melhor = %.2f), esperado <= 10%% do otimo %.2f",
			gap, res.Melhor, otimo)
	}
	if err := ValidaRota(res.MelhorRota.Ordem, in); err != nil {
		t.Fatalf("MelhorRota invalida: %v", err)
	}
}

func TestMemeticoTextoResumo(t *testing.T) {
	in := carregaInstanciaTeste(t, "PROBLEMA_11_KM_06.txt")
	alg := AlgoritmoMemetico{Par: ParametrosMemeticoPadrao()}
	res := ExecutaExperimento(alg, in, 3, 42)
	texto := res.TextoResumo(in)

	if !strings.Contains(texto, "RESUMO DO EXPERIMENTO - ALGORITMO MEMETICO\n") {
		t.Errorf("texto do resumo nao identifica o Algoritmo Memetico")
	}
	if !strings.Contains(texto, "Menor valor      : 344.90\n") {
		t.Errorf("texto do resumo nao converge ao otimo conhecido da instancia")
	}
}

func TestEducaPopulacaoNuncaPiora(t *testing.T) {
	in := carregaInstanciaTeste(t, "PROBLEMA_07_KM_12.txt")
	rng := rngTeste(42)
	pop := PopulacaoInicial(20, in, rng)
	custosAntes := make([]float64, len(pop))
	for i, r := range pop {
		custosAntes[i] = r.Custo
	}

	pop = educaPopulacao(pop, 1.0, in, rng)

	for i, r := range pop {
		if r.Custo > custosAntes[i]+1e-9 {
			t.Fatalf("individuo %d piorou: antes %.4f, depois %.4f", i, custosAntes[i], r.Custo)
		}
		if err := ValidaRota(r.Ordem, in); err != nil {
			t.Fatalf("individuo %d invalido apos educacao: %v", i, err)
		}
		if got := CustoRota(r.Ordem, in); math.Abs(got-r.Custo) > 1e-9 {
			t.Fatalf("individuo %d: custo %f inconsistente com a ordem (%f)", i, r.Custo, got)
		}
	}
}

func TestEducaPopulacaoTaxaZeroNaoAltera(t *testing.T) {
	in := carregaInstanciaTeste(t, "PROBLEMA_07_KM_12.txt")
	rng := rngTeste(42)
	pop := PopulacaoInicial(10, in, rng)
	custosAntes := make([]float64, len(pop))
	for i, r := range pop {
		custosAntes[i] = r.Custo
	}

	pop = educaPopulacao(pop, 0.0, in, rng)

	for i, r := range pop {
		if r.Custo != custosAntes[i] {
			t.Fatalf("individuo %d alterado com taxa 0.0", i)
		}
	}
}
