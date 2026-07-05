package pcv

// OrOpt implementa a busca local Or-opt.
// O algoritmo tenta melhorar uma rota movendo blocos
// consecutivos de 1, 2 ou 3 cidades para outra posição da rota.
type OrOpt struct{}

// Nome devolve o nome da busca local.
func (OrOpt) Nome() string {
	return "Or-opt"
}

// moverBloco remove um bloco da rota e o insere em uma nova posição.
func moverBloco(ordem []int, inicio, tamanho, destino int) []int {
	// Copia o bloco que será movido.
	bloco := append([]int{}, ordem[inicio:inicio+tamanho]...)

	// Remove o bloco da rota.
	restante := make([]int, 0, len(ordem)-tamanho)
	restante = append(restante, ordem[:inicio]...)
	restante = append(restante, ordem[inicio+tamanho:]...)

	// Garante que o destino esteja dentro dos limites.
	if destino < 0 {
		destino = 0
	}
	if destino > len(restante) {
		destino = len(restante)
	}

	// Monta a nova rota.
	nova := make([]int, 0, len(ordem))
	nova = append(nova, restante[:destino]...)
	nova = append(nova, bloco...)
	nova = append(nova, restante[destino:]...)

	return nova
}

// Aplica executa a busca local Or-opt.
// Em cada iteração são testados movimentos de blocos
// de tamanho 1, 2 e 3, escolhendo sempre a melhor melhoria.
func (OrOpt) Aplica(r Rota, in *Instancia) Rota {
	atual := r.Clona()

	for {
		melhorou := false
		melhorVizinho := atual
		melhorCusto := atual.Custo

		// Testa blocos de tamanho 1, 2 e 3.
		for tamanho := 1; tamanho <= 3; tamanho++ {

			// Escolhe o início do bloco.
			for inicio := 1; inicio < len(atual.Ordem)-tamanho; inicio++ {

				// Testa todas as posições possíveis para inserir o bloco.
				for destino := 1; destino <= len(atual.Ordem)-tamanho; destino++ {

					// Ignora movimentos que não alteram a rota.
					if destino == inicio {
						continue
					}

					novaOrdem := moverBloco(
						atual.Ordem,
						inicio,
						tamanho,
						destino,
					)

					candidato := NovaRota(novaOrdem, in)

					if candidato.Custo < melhorCusto {
						melhorVizinho = candidato
						melhorCusto = candidato.Custo
						melhorou = true
					}
				}
			}
		}

		// Se nenhuma melhoria foi encontrada,
		// a busca chegou a um ótimo local.
		if !melhorou {
			break
		}

		// Continua a busca a partir da melhor solução encontrada.
		atual = melhorVizinho
	}

	return atual
}
