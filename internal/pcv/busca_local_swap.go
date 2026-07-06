package pcv

type Swap struct{}

func (Swap) Nome() string { return "Swap" }

func (Swap) Aplica(r Rota, in *Instancia) Rota {
	atual := r.Clona()
	n := len(atual.Ordem)

	for {
		melhorou := false
		melhorI, melhorJ := -1, -1
		melhorDelta := 0.0

		for i := 1; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				delta := deltaSwap(atual.Ordem, i, j, in)
				if delta < melhorDelta-1e-9 {
					melhorDelta = delta
					melhorI = i
					melhorJ = j
					melhorou = true
				}
			}
		}

		if melhorou {
			atual.Ordem[melhorI], atual.Ordem[melhorJ] = atual.Ordem[melhorJ], atual.Ordem[melhorI]
			atual.Custo += melhorDelta
		} else {
			break
		}
	}

	atual.Custo = CustoRota(atual.Ordem, in)
	return atual
}

func deltaSwap(ordem []int, i, j int, in *Instancia) float64 {
	n := len(ordem)
	antI, propI := ordem[i-1], ordem[(i+1)%n]
	antJ, propJ := ordem[j-1], ordem[(j+1)%n]
	ci, cj := ordem[i], ordem[j]

	if j == i+1 {
		return in.Custo(antI, cj) + in.Custo(ci, propJ) -
			in.Custo(antI, ci) - in.Custo(cj, propJ)
	}

	return in.Custo(antI, cj) + in.Custo(cj, propI) +
		in.Custo(antJ, ci) + in.Custo(ci, propJ) -
		in.Custo(antI, ci) - in.Custo(ci, propI) -
		in.Custo(antJ, cj) - in.Custo(cj, propJ)
}
