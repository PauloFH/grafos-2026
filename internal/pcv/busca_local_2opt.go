package pcv

// DoisOpt tenta melhorar uma rota invertendo trechos do percurso.
type DoisOpt struct{}

func (DoisOpt) Nome() string {
	return "2-opt"
}

// inverterTrecho inverte os elementos do slice entre as posições inicio e fim.
func inverterTrecho(ordem []int, inicio, fim int) {
	for inicio < fim {
		ordem[inicio], ordem[fim] = ordem[fim], ordem[inicio]
		inicio++
		fim--
	}
}

// Aplica executa a busca local 2-opt sobre a rota recebida.
// Utiliza avaliação por Delta em O(1) e trata o fechamento do ciclo.
func (DoisOpt) Aplica(r Rota, in *Instancia) Rota {
	atual := r.Clona()
	n := len(atual.Ordem)

	for {
		melhorou := false
		melhorI, melhorJ := -1, -1
		melhorDelta := 0.0 // Delta negativo significa que o custo diminuiu!

		// O 'i' vai até o penúltimo elemento
		for i := 1; i < n-1; i++ {

			for j := i + 1; j < n; j++ {

				// Evita o movimento degenerado (inverter a rota inteira não altera o custo no PCV simétrico)
				if i == 1 && j == n-1 {
					continue
				}

				// "Efeito Relógio": Se o j for o último, a próxima cidade (j+1) é o Depósito (índice 0)
				proxJ := (j + 1) % n

				// Cálculo do Delta: somamos os custos das arestas NOVAS e subtraímos as arestas VELHAS
				delta := in.Custo(atual.Ordem[i-1], atual.Ordem[j]) +
					in.Custo(atual.Ordem[i], atual.Ordem[proxJ]) -
					in.Custo(atual.Ordem[i-1], atual.Ordem[i]) -
					in.Custo(atual.Ordem[j], atual.Ordem[proxJ])

				if delta < melhorDelta-1e-9 {
					melhorDelta = delta
					melhorI = i
					melhorJ = j
					melhorou = true
				}
			}
		}

		if melhorou {
			// Aplica a inversão definitiva no mesmo array (in-place) sem alocar memória nova
			inverterTrecho(atual.Ordem, melhorI, melhorJ)
			atual.Custo += melhorDelta
		} else {
			break
		}
	}

	// Recalcula o custo final apenas por garantia contra pequenos erros de ponto flutuante
	atual.Custo = CustoRota(atual.Ordem, in)
	return atual
}
