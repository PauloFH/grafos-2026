// Package pcv implementa a Unidade 3: heuristicas construtivas, buscas locais
// e metaheuristicas (algoritmo genetico e memetico) para o Problema do
// Caixeiro Viajante.
package pcv

import "fmt"

// Instancia representa uma instancia do PCV carregada de inputs_u3/*.txt.
// A matriz e SIMETRICA e indexada por POSICAO (0..N-1); os ids originais
// (1..48) so aparecem em Cidades, Nomes e nas saidas para humanos.
type Instancia struct {
	Nome    string      // vem da linha 1 do arquivo, nao do nome do arquivo
	Medida  string      // "KM" ou "MIN"
	N       int         // numero de cidades
	Cidades []int       // ids originais; Cidades[0] = deposito (cidade 1)
	Nomes   []string    // mesma ordem de Cidades
	Matriz  [][]float64 // N x N, simetrica, indexada por posicao
}

// Custo devolve o custo entre duas POSICOES (0..N-1), nao ids originais.
func (in *Instancia) Custo(i, j int) float64 {
	return in.Matriz[i][j]
}

// Rota e uma solucao do PCV: permutacao das posicoes 0..N-1 com Ordem[0]==0
// (deposito fixo). O retorno ao deposito NAO e repetido no slice - esta
// implicito no custo.
type Rota struct {
	Ordem []int
	Custo float64 // custo do ciclo completo, incluindo a volta ao deposito
}

// CustoRota calcula o custo do ciclo completo definido por ordem, incluindo
// a volta ao deposito.
func CustoRota(ordem []int, in *Instancia) float64 {
	total := 0.0
	for k := 0; k < len(ordem)-1; k++ {
		total += in.Custo(ordem[k], ordem[k+1])
	}
	if len(ordem) > 1 {
		total += in.Custo(ordem[len(ordem)-1], ordem[0])
	}
	return total
}

// ValidaRota devolve nil se ordem tem tamanho N, e uma permutacao de 0..N-1
// e comeca no deposito; caso contrario, erro com a primeira violacao.
func ValidaRota(ordem []int, in *Instancia) error {
	if len(ordem) != in.N {
		return fmt.Errorf("rota com tamanho %d, esperado %d", len(ordem), in.N)
	}
	if ordem[0] != 0 {
		return fmt.Errorf("rota nao inicia no deposito (posicao 0)")
	}
	visto := make([]bool, in.N)
	for _, p := range ordem {
		if p < 0 || p >= in.N {
			return fmt.Errorf("posicao %d fora do intervalo [0,%d)", p, in.N)
		}
		if visto[p] {
			return fmt.Errorf("posicao %d repetida na rota", p)
		}
		visto[p] = true
	}
	return nil
}

// NovaRota monta uma Rota com copia propria de ordem (sem aliasing com o
// slice do chamador) e o custo ja calculado.
func NovaRota(ordem []int, in *Instancia) Rota {
	copia := make([]int, len(ordem))
	copy(copia, ordem)
	return Rota{Ordem: copia, Custo: CustoRota(copia, in)}
}

// Clona devolve uma copia profunda da rota. Obrigatoria ao guardar ou
// derivar rotas: atribuicao simples compartilha o slice Ordem e corrompe a
// populacao silenciosamente.
func (r Rota) Clona() Rota {
	copia := make([]int, len(r.Ordem))
	copy(copia, r.Ordem)
	return Rota{Ordem: copia, Custo: r.Custo}
}

// IdsOriginais converte a rota de posicoes para ids originais de cidade
// (1..48, sem repetir o deposito no fim).
func (r Rota) IdsOriginais(in *Instancia) []int {
	ids := make([]int, len(r.Ordem))
	for k, p := range r.Ordem {
		ids[k] = in.Cidades[p]
	}
	return ids
}

// BuscaLocal e o contrato das buscas locais (2-opt, Or-opt, Swap): Aplica
// devolve uma rota valida com custo menor ou igual, sem mutar a entrada;
// deterministica, roda ate o otimo local do movimento.
type BuscaLocal interface {
	Nome() string
	Aplica(r Rota, in *Instancia) Rota
}

// Construtivo e o contrato das heuristicas construtivas (Vizinho mais
// Proximo, Insercao mais Barata): Constroi devolve uma rota valida com
// custo preenchido; deterministico.
type Construtivo interface {
	Nome() string
	Constroi(in *Instancia) Rota
}

// Estocastico e o contrato das metaheuristicas (AG, memetico): Executa faz
// uma execucao completa com a semente dada e devolve a melhor rota (clonada).
// Rng canonico: rand.New(rand.NewPCG(uint64(semente), uint64(semente))).
type Estocastico interface {
	Nome() string
	Executa(in *Instancia, semente int64) Rota
}
