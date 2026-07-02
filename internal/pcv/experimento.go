package pcv

import "time"

// ResumoExperimento agrega o resultado de varias execucoes independentes de
// um algoritmo estocastico sobre uma instancia; e a fonte de dados dos
// arquivos outputs_u3/RESUMO_*.txt.
type ResumoExperimento struct {
	Algoritmo  string        // alg.Nome()
	Instancia  string        // Instancia.Nome
	Execucoes  int           // numero de execucoes (20 no trabalho)
	Melhor     float64       // menor custo entre as execucoes
	MelhorRota Rota          // rota da melhor execucao (clonada)
	Media      float64       // media dos custos
	TempoMedio time.Duration // tempo medio de uma chamada de Executa
	Custos     []float64     // custo de cada execucao, na ordem das sementes
}

// ExecutaExperimento roda alg.Executa execucoes vezes sobre in, usando a
// semente sementeBase+int64(i) na execucao i, e agrega os resultados; em
// empate de custo vale a primeira melhor rota encontrada.
func ExecutaExperimento(alg Estocastico, in *Instancia, execucoes int, sementeBase int64) ResumoExperimento {
	panic("pcv: nao implementado - Parte 3 (genetico)")
}

// TextoResumo gera o texto canonico dos arquivos outputs_u3/RESUMO_*.txt
// (menor valor, media, tempo medio, gap e melhor rota em ids originais).
func (r ResumoExperimento) TextoResumo(in *Instancia) string {
	panic("pcv: nao implementado - Parte 3 (genetico)")
}
