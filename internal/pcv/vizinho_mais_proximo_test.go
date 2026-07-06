package pcv

import (
	"math"
	"testing"
)

// instanciasConstrutivas sao os arquivos reais usados nos testes dos dois
// construtivos: uma pequena (P11), uma media (P07) e a maior (P01).
var instanciasConstrutivas = []string{
	"PROBLEMA_11_KM_06.txt",
	"PROBLEMA_12_MIN_06.txt",
	"PROBLEMA_07_KM_12.txt",
	"PROBLEMA_01_KM_48.txt",
}

// verificaConstrutivo roda um construtivo sobre a instancia e confere as
// invariantes comuns a qualquer heuristica construtiva: rota valida iniciando
// no deposito, custo coerente com CustoRota e nao inferior ao otimo conhecido
// (que e cota inferior). Compartilhado com o teste da Insercao mais Barata.
func verificaConstrutivo(t *testing.T, c Construtivo, in *Instancia) Rota {
	t.Helper()
	r := c.Constroi(in)

	if err := ValidaRota(r.Ordem, in); err != nil {
		t.Fatalf("%s/%s: rota invalida: %v", c.Nome(), in.Nome, err)
	}
	if r.Ordem[0] != 0 {
		t.Errorf("%s/%s: rota nao inicia no deposito, Ordem[0]=%d", c.Nome(), in.Nome, r.Ordem[0])
	}
	if esperado := CustoRota(r.Ordem, in); math.Abs(r.Custo-esperado) > 1e-9 {
		t.Errorf("%s/%s: Custo=%.4f, esperado %.4f", c.Nome(), in.Nome, r.Custo, esperado)
	}
	if otimo, ok := ValoresOtimos[in.Nome]; ok && r.Custo < otimo-1e-6 {
		t.Errorf("%s/%s: custo %.2f abaixo do otimo %.2f (impossivel)", c.Nome(), in.Nome, r.Custo, otimo)
	}
	return r
}

// TestVizinhoMaisProximoRotaValida confere as invariantes do construtivo em
// todas as instancias de referencia.
func TestVizinhoMaisProximoRotaValida(t *testing.T) {
	for _, arquivo := range instanciasConstrutivas {
		t.Run(arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, arquivo)
			verificaConstrutivo(t, VizinhoMaisProximo{}, in)
		})
	}
}

// TestVizinhoMaisProximoDeterministico confere que duas construcoes produzem
// exatamente a mesma rota (nao ha aleatoriedade envolvida).
func TestVizinhoMaisProximoDeterministico(t *testing.T) {
	for _, arquivo := range instanciasConstrutivas {
		t.Run(arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, arquivo)
			r1 := VizinhoMaisProximo{}.Constroi(in)
			r2 := VizinhoMaisProximo{}.Constroi(in)
			if r1.Custo != r2.Custo {
				t.Fatalf("custos divergem: %.4f vs %.4f", r1.Custo, r2.Custo)
			}
			for i := range r1.Ordem {
				if r1.Ordem[i] != r2.Ordem[i] {
					t.Fatalf("Ordem diverge na posicao %d: %d vs %d", i, r1.Ordem[i], r2.Ordem[i])
				}
			}
		})
	}
}
