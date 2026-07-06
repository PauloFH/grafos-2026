package pcv

import "testing"

// TestInsercaoMaisBarataRotaValida confere as invariantes do construtivo em
// todas as instancias de referencia (helper e casos compartilhados com o
// teste do Vizinho mais Proximo).
func TestInsercaoMaisBarataRotaValida(t *testing.T) {
	for _, arquivo := range instanciasConstrutivas {
		t.Run(arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, arquivo)
			verificaConstrutivo(t, InsercaoMaisBarata{}, in)
		})
	}
}

// TestInsercaoMaisBarataDeterministico confere que duas construcoes produzem
// exatamente a mesma rota (nao ha aleatoriedade envolvida).
func TestInsercaoMaisBarataDeterministico(t *testing.T) {
	for _, arquivo := range instanciasConstrutivas {
		t.Run(arquivo, func(t *testing.T) {
			in := carregaInstanciaTeste(t, arquivo)
			r1 := InsercaoMaisBarata{}.Constroi(in)
			r2 := InsercaoMaisBarata{}.Constroi(in)
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
