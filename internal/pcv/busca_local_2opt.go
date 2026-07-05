package pcv

//O algoritmo tenta melhorar uma rota invertendo trechos do percurso.
type DoisOpt struct{}

func (DoisOpt) Nome() string {
	return "2-opt"
}

//inverterTrecho inverte os elementos do slice entre as posições
func inverterTrecho(ordem []int, inicio, fim int) {
	for inicio < fim {
		ordem[inicio], ordem[fim] = ordem[fim], ordem[inicio]
		inicio++
		fim--
	}
}

// Aplica executa a busca local 2-opt sobre a rota recebida.
// A rota original não é modificada, sempre é utilizada uma cópia.
// (Best Improvement). O algoritmo termina quando não há mais melhorias.
func (DoisOpt) Aplica(r Rota, in *Instancia) Rota {
	//Cria uma cópia da rota
	atual := r.Clona()

	for {
		melhorou := false

		// Guarda o melhor vizinho encontrado nesta iteração.
		melhorVizinho := atual
		melhorCusto := atual.Custo

		// Testa todas as combinações possíveis de inversão.
		for i := 1; i < len(atual.Ordem)-2; i++ {
			for j := i + 1; j < len(atual.Ordem)-1; j++ {

				// Cria uma nova rota candidata.
				candidato := atual.Clona()

				// Inverte o trecho selecionado
				inverterTrecho(candidato.Ordem, i, j)

				//calcula o custo
				candidato.Custo = CustoRota(candidato.Ordem, in)

				// Caso encontre uma solução melhor, guarda como melhor vizinho da iteração.
				if candidato.Custo < melhorCusto {
					melhorVizinho = candidato
					melhorCusto = candidato.Custo
					melhorou = true
				}
			}
		}

		// Se nenhuma melhoria foi encontrada,
		// a busca chegou a um ótimo local.
		if !melhorou {
			break
		}

		//continua a busca a partir da melhor solução
		atual = melhorVizinho
	}

	return atual
}
