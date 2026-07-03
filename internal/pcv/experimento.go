package pcv

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ResumoExperimento struct {
	Algoritmo  string        
	Instancia  string        
	Execucoes  int          
	Melhor     float64       
	MelhorRota Rota         
	Media      float64       
	TempoMedio time.Duration 
	Custos     []float64     
}


func ExecutaExperimento(alg Estocastico, in *Instancia, execucoes int, sementeBase int64) ResumoExperimento {
	res := ResumoExperimento{
		Algoritmo: alg.Nome(),
		Instancia: in.Nome,
		Execucoes: execucoes,
		Custos:    make([]float64, 0, execucoes),
	}
	somaCusto := 0.0
	somaTempo := time.Duration(0)

	//O tempo medio deve refletir o custo real de uma execucao e o determinismo por semente nao pode depender de goroutines
	for i := 0; i < execucoes; i++ {
		semente := sementeBase + int64(i) 
		inicio := time.Now()
		rota := alg.Executa(in, semente)
		dur := time.Since(inicio)

		res.Custos = append(res.Custos, rota.Custo)
		somaCusto += rota.Custo
		somaTempo += dur

		// empate de custo mantem a PRIMEIRA rota encontrada
		if i == 0 || rota.Custo < res.Melhor {
			res.Melhor = rota.Custo
			res.MelhorRota = rota.Clona()
		}
	}

	res.Media = somaCusto / float64(execucoes)
	res.TempoMedio = somaTempo / time.Duration(execucoes)
	return res
}

func (r ResumoExperimento) TextoResumo(in *Instancia) string {
	var b strings.Builder
	igual := strings.Repeat("=", 50)
	traco := strings.Repeat("-", 50)

	fmt.Fprintf(&b, "%s\n", igual)
	fmt.Fprintf(&b, "RESUMO DO EXPERIMENTO - %s\n", strings.ToUpper(r.Algoritmo))
	fmt.Fprintf(&b, "%s\n", igual)

	fmt.Fprintf(&b, "%-17s: %s (%s, N=%d)\n", "Instancia", r.Instancia, in.Medida, in.N)
	fmt.Fprintf(&b, "%-17s: %d\n", "Execucoes", r.Execucoes)
	fmt.Fprintf(&b, "%-17s: %.2f\n", "Menor valor", r.Melhor)
	fmt.Fprintf(&b, "%-17s: %.2f\n", "Valor medio", r.Media)
	fmt.Fprintf(&b, "%-17s: %.3fs\n", "Tempo medio", r.TempoMedio.Seconds())
	otimo := ValoresOtimos[r.Instancia] // 0.00 se a instancia nao esta no mapa
	fmt.Fprintf(&b, "%-17s: %.2f\n", "Melhor conhecido", otimo)
	fmt.Fprintf(&b, "%-17s: %.2f%%\n", "Gap do menor", Gap(r.Melhor, otimo))

	// melhor rota em ids ORIGINAIS, com o deposito repetido no fim para explicitar o fechamento do ciclo
	fmt.Fprintf(&b, "%s\n", traco)
	b.WriteString("Melhor rota (ids das cidades):\n")
	ids := r.MelhorRota.IdsOriginais(in)
	partes := make([]string, 0, len(ids)+1)
	for _, id := range ids {
		partes = append(partes, strconv.Itoa(id))
	}
	if len(ids) > 0 {
		partes = append(partes, strconv.Itoa(ids[0]))
	}
	fmt.Fprintf(&b, "%s\n", strings.Join(partes, " -> "))

	fmt.Fprintf(&b, "%s\n", traco)
	fmt.Fprintf(&b, "Custos das %d execucoes:\n", r.Execucoes)
	for i, c := range r.Custos {
		fmt.Fprintf(&b, "%4d: %.2f\n", i+1, c)
	}
	return b.String()
}
