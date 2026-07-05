package pcv

import (
	"math/rand/v2"
	"testing"
)

// testaBuscaLocal executa os testes comuns para qualquer busca local.
func testaBuscaLocal(t *testing.T, busca BuscaLocal) {
	t.Helper()

	in := carregaInstanciaTeste(t, "PROBLEMA_11_KM_06.txt")

	// Gera uma rota inicial determinística.
	rng := rand.New(rand.NewPCG(42, 42))
	rotaOriginal := PopulacaoInicial(1, in, rng)[0]

	// Guarda uma cópia da rota original para verificar
	// que ela não foi modificada pela busca.
	ordemOriginal := append([]int(nil), rotaOriginal.Ordem...)
	custoOriginal := rotaOriginal.Custo

	// Executa a busca local.
	rotaMelhorada := busca.Aplica(rotaOriginal, in)

	// A rota resultante deve continuar válida.
	if err := ValidaRota(rotaMelhorada.Ordem, in); err != nil {
		t.Fatalf("%s produziu uma rota inválida: %v", busca.Nome(), err)
	}

	// A busca local nunca pode piorar a solução.
	if rotaMelhorada.Custo > custoOriginal+1e-9 {
		t.Fatalf(
			"%s piorou a rota: custo inicial %.2f, custo final %.2f",
			busca.Nome(),
			custoOriginal,
			rotaMelhorada.Custo,
		)
	}

	// A rota original não pode ser alterada.
	if rotaOriginal.Custo != custoOriginal {
		t.Fatalf("%s alterou o custo da rota original", busca.Nome())
	}

	for i := range ordemOriginal {
		if rotaOriginal.Ordem[i] != ordemOriginal[i] {
			t.Fatalf("%s alterou a rota original", busca.Nome())
		}
	}
}

func TestDoisOptAplica(t *testing.T) {
	testaBuscaLocal(t, DoisOpt{})
}

func TestOrOptAplica(t *testing.T) {
	testaBuscaLocal(t, OrOpt{})
}
