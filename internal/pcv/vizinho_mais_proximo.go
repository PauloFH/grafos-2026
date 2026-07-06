package pcv

import "math"

// VizinhoMaisProximo implementa a interface Construtivo com a heuristica do
// vizinho mais proximo (nearest neighbor).
type VizinhoMaisProximo struct{}

// Nome devolve o nome canonico da heuristica usado nos relatorios.
func (VizinhoMaisProximo) Nome() string { return "Vizinho mais Proximo" }

// Constroi parte do deposito (posicao 0) e, a cada passo, move-se para a
// cidade nao visitada de menor custo. Empates ficam com a de MENOR indice
// (comparacao estrita), garantindo determinismo. Devolve uma rota valida com
// o custo do ciclo completo ja preenchido.
func (VizinhoMaisProximo) Constroi(in *Instancia) Rota {
	visitado := make([]bool, in.N)
	ordem := make([]int, 0, in.N)

	// deposito fixo na posicao 0
	atual := 0
	visitado[0] = true
	ordem = append(ordem, 0)

	for len(ordem) < in.N {
		melhor := -1
		melhorCusto := math.Inf(1)
		for c := 0; c < in.N; c++ {
			if visitado[c] {
				continue
			}
			// < estrito preserva a cidade de menor indice em empates
			if d := in.Custo(atual, c); d < melhorCusto {
				melhorCusto = d
				melhor = c
			}
		}
		visitado[melhor] = true
		ordem = append(ordem, melhor)
		atual = melhor
	}

	return NovaRota(ordem, in)
}
