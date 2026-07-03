package pcv

import (
	"math"
	"strings"
	"testing"
)

func TestExecutaReprodutibilidade(t *testing.T) {
	casos := []struct {
		arquivo string
		semente int64
	}{
		{"PROBLEMA_11_KM_06.txt", 42},
		{"PROBLEMA_07_KM_12.txt", 42},
		{"PROBLEMA_07_KM_12.txt", 7},
	}
	alg := AlgoritmoGenetico{Par: ParametrosPadrao()}
	for _, c := range casos {
		t.Run(c.arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, c.arquivo)
			r1 := alg.Executa(in, c.semente)
			r2 := alg.Executa(in, c.semente)
			if r1.Custo != r2.Custo {
				t.Fatalf("custos diferentes: %f != %f", r1.Custo, r2.Custo)
			}
			if len(r1.Ordem) != len(r2.Ordem) {
				t.Fatalf("tamanhos diferentes: %d != %d", len(r1.Ordem), len(r2.Ordem))
			}
			for k := range r1.Ordem {
				if r1.Ordem[k] != r2.Ordem[k] {
					t.Fatalf("Ordem difere na posicao %d: %d != %d", k, r1.Ordem[k], r2.Ordem[k])
				}
			}
		})
	}
}

func TestExecutaRotaValida(t *testing.T) {
	casos := []string{
		"PROBLEMA_11_KM_06.txt",
		"PROBLEMA_09_KM_07.txt", // ids nao contiguos
		"PROBLEMA_07_KM_12.txt",
	}
	alg := AlgoritmoGenetico{Par: ParametrosPadrao()}
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

func TestDoDOtimosP11P12(t *testing.T) {
	casos := []struct {
		arquivo string
		otimo   float64
	}{
		{"PROBLEMA_11_KM_06.txt", 344.90},
		{"PROBLEMA_12_MIN_06.txt", 305.00},
	}
	alg := AlgoritmoGenetico{Par: ParametrosPadrao()}
	for _, c := range casos {
		t.Run(c.arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, c.arquivo)
			res := ExecutaExperimento(alg, in, 20, 42)

			if len(res.Custos) != 20 {
				t.Fatalf("len(Custos) = %d, esperado 20", len(res.Custos))
			}
			for i, custo := range res.Custos {
				if math.Abs(custo-c.otimo) > 1e-6 {
					t.Errorf("execucao %d: custo %.2f, esperado otimo %.2f", i+1, custo, c.otimo)
				}
			}
			if math.Abs(res.Melhor-c.otimo) > 1e-6 {
				t.Errorf("Melhor = %.2f, esperado %.2f", res.Melhor, c.otimo)
			}
			if math.Abs(res.Media-c.otimo) > 1e-6 {
				t.Errorf("Media = %.2f, esperado %.2f", res.Media, c.otimo)
			}
			menor := res.Custos[0]
			for _, custo := range res.Custos {
				if custo < menor {
					menor = custo
				}
			}
			if res.Melhor != menor {
				t.Errorf("Melhor = %f difere de min(Custos) = %f", res.Melhor, menor)
			}
			if res.MelhorRota.Custo != res.Melhor {
				t.Errorf("MelhorRota.Custo = %f difere de Melhor = %f", res.MelhorRota.Custo, res.Melhor)
			}
			if err := ValidaRota(res.MelhorRota.Ordem, in); err != nil {
				t.Errorf("MelhorRota invalida: %v", err)
			}
			if res.TempoMedio <= 0 {
				t.Errorf("TempoMedio = %v, esperado > 0", res.TempoMedio)
			}
		})
	}
}

func TestDoDGapP01(t *testing.T) {
	in := carregaInstanciaTeste(t, "PROBLEMA_01_KM_48.txt")
	alg := AlgoritmoGenetico{Par: ParametrosPadrao()}
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
func TestTextoResumo(t *testing.T) {
	in := carregaInstanciaTeste(t, "PROBLEMA_11_KM_06.txt")
	alg := AlgoritmoGenetico{Par: ParametrosPadrao()}
	res := ExecutaExperimento(alg, in, 3, 42)
	texto := res.TextoResumo(in)

	linhas := []string{
		strings.Repeat("=", 50) + "\n",
		"RESUMO DO EXPERIMENTO - ALGORITMO GENETICO\n",
		"Instancia        : PROBLEMA_11 (KM, N=6)\n",
		"Execucoes        : 3\n",
		"Menor valor      : 344.90\n",
		"Valor medio      : 344.90\n",
		"Melhor conhecido : 344.90\n",
		"Gap do menor     : 0.00%\n",
		strings.Repeat("-", 50) + "\n",
		"Melhor rota (ids das cidades):\n",
		"Custos das 3 execucoes:\n",
		"   1: 344.90\n",
		"   2: 344.90\n",
		"   3: 344.90\n",
	}
	for _, l := range linhas {
		if !strings.Contains(texto, l) {
			t.Errorf("texto do resumo nao contem a linha canonica %q", l)
		}
	}
	if !strings.Contains(texto, "Tempo medio      : 0.") || !strings.Contains(texto, "s\n") {
		t.Errorf("linha de tempo fora do formato %%-17s / %%.3fs")
	}
	if !strings.Contains(texto, "\n1 -> ") {
		t.Errorf("melhor rota nao inicia no deposito (id 1)")
	}
	if !strings.Contains(texto, " -> 1\n") {
		t.Errorf("melhor rota nao fecha no deposito (id 1)")
	}
	if !strings.HasSuffix(texto, "\n") {
		t.Errorf("texto nao termina em \\n")
	}
}
