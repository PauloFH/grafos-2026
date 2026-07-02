package web

// DTOs da API /api/u3/* (Unidade 3 - Caixeiro Viajante). Os tipos de
// internal/pcv nao tem tags JSON; mapeamos para estes DTOs (camelCase).

// InstanciaInfo resume uma instancia para a lista /api/u3/instancias.
type InstanciaInfo struct {
	ID     string `json:"id"`     // ex.: PROBLEMA_01
	Medida string `json:"medida"` // KM ou MIN
	N      int    `json:"n"`
}

// PontoDTO e a posicao 2D de uma cidade (MDS sobre a matriz de custos,
// layout ilustrativo), com coordenadas normalizadas em [0.05, 0.95].
type PontoDTO struct {
	ID   int     `json:"id"` // id ORIGINAL da cidade (1..48)
	Nome string  `json:"nome"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

// InstanciaResponse e a resposta de /api/u3/instancia.
type InstanciaResponse struct {
	ID       string     `json:"id"`
	Medida   string     `json:"medida"`
	N        int        `json:"n"`
	Deposito int        `json:"deposito"` // id original do deposito (sempre 1)
	Pontos   []PontoDTO `json:"pontos"`   // nunca null
	Otimo    float64    `json:"otimo"`    // 0 se ausente de ValoresOtimos
}

// RotaDTO e uma rota serializavel: Ordem usa ids ORIGINAIS (1..48), comeca
// no deposito e NAO o repete no fim; rota vazia serializa como [].
type RotaDTO struct {
	Ordem []int   `json:"ordem"`
	Custo float64 `json:"custo"` // ciclo completo, incluindo a volta
}

// SolveResponse e a resposta de /api/u3/resolve. Erros de negocio (metodo
// desconhecido, semente invalida, nao implementado) respondem HTTP 200 com
// OK=false, como na U2.
type SolveResponse struct {
	OK         bool     `json:"ok"`
	Erro       string   `json:"erro"`   // "" em sucesso
	Metodo     string   `json:"metodo"` // vmp | imb | ag | memetico
	Antes      *RotaDTO `json:"antes"`  // construtivo antes da busca local; null para ag/memetico e em erro
	Rota       RotaDTO  `json:"rota"`
	GapPercent float64  `json:"gapPercent"`
	TempoMs    float64  `json:"tempoMs"`
}
