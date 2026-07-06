package pcv

import "math"

// InsercaoMaisBarata implementa a interface Construtivo com a heuristica da
// insercao mais barata (cheapest insertion).
type InsercaoMaisBarata struct{}

// Nome devolve o nome canonico da heuristica usado nos relatorios.
func (InsercaoMaisBarata) Nome() string { return "Insercao mais Barata" }

// Constroi comeca com um subtour deposito -> cidade mais proxima do deposito
// e, a cada passo, insere a cidade nao visitada na aresta que gera o MENOR
// acrescimo de custo delta = c(i,k) + c(k,j) - c(i,j). A varredura e por
// cidade crescente e depois por posicao crescente, com comparacao estrita,
// de modo que empates sejam resolvidos de forma deterministica. O deposito
// permanece na posicao 0 (nunca se insere antes dela). Devolve uma rota
// valida com o custo do ciclo completo ja preenchido.
func (InsercaoMaisBarata) Constroi(in *Instancia) Rota {
	if in.N == 1 {
		return NovaRota([]int{0}, in)
	}

	visitado := make([]bool, in.N)
	visitado[0] = true

	// semente: cidade mais proxima do deposito (menor indice em empates)
	c0 := -1
	menor := math.Inf(1)
	for c := 1; c < in.N; c++ {
		if d := in.Custo(0, c); d < menor {
			menor = d
			c0 = c
		}
	}
	visitado[c0] = true
	tour := []int{0, c0}

	for len(tour) < in.N {
		melhorK, melhorPos := -1, -1
		melhorDelta := math.Inf(1)
		for k := 0; k < in.N; k++ {
			if visitado[k] {
				continue
			}
			for i := 0; i < len(tour); i++ {
				j := (i + 1) % len(tour)
				delta := in.Custo(tour[i], k) + in.Custo(k, tour[j]) - in.Custo(tour[i], tour[j])
				// < estrito preserva o primeiro (k, posicao) em empates
				if delta < melhorDelta {
					melhorDelta = delta
					melhorK = k
					melhorPos = i
				}
			}
		}
		// insere melhorK logo apos a posicao melhorPos
		tour = append(tour, 0)
		copy(tour[melhorPos+2:], tour[melhorPos+1:])
		tour[melhorPos+1] = melhorK
		visitado[melhorK] = true
	}

	return NovaRota(tour, in)
}
