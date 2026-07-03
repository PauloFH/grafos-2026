// Package web expoe a Unidade 2 (Projeto 3 - EcoUrbano) via HTTP/JSON.
// Envelopa os pacotes internos (euleriano, grafo, leitor) sem altera-los:
// os tipos de dominio nao possuem tags JSON, entao mapeamos para os DTOs
// abaixo antes de serializar.
package web

// DatasetInfo resume um grafo carregado, para a lista /api/datasets.
type DatasetInfo struct {
	ID          string `json:"id"`          // nome do arquivo sem extensao (ex.: GRAFO_EULER)
	Nome        string `json:"nome"`        // nome do arquivo de origem
	Direcionado bool   `json:"direcionado"` // true = digrafo (mao unica)
	Vertices    int    `json:"vertices"`    // total de vertices
	Arestas     int    `json:"arestas"`     // total de arestas/arcos
}

// ArestaDTO e uma rua do setor (par origem/destino).
type ArestaDTO struct {
	De   string `json:"de"`
	Para string `json:"para"`
}

// GraphResponse e a malha viaria desenhavel: vertices, arestas e graus.
// Para grafos nao-direcionados as arestas vem deduplicadas (uma rua = um par).
type GraphResponse struct {
	ID          string         `json:"id"`
	Nome        string         `json:"nome"`
	Direcionado bool           `json:"direcionado"`
	Vertices    []string       `json:"vertices"`
	Arestas     []ArestaDTO    `json:"arestas"`
	Graus       map[string]int `json:"graus"`
}

// ClassificationDTO espelha euleriano.ResultadoEuler com tags JSON.
type ClassificationDTO struct {
	Classe         string         `json:"classe"`         // CircuitoEuleriano | CaminhoEuleriano | NaoEuleriano
	Texto          string         `json:"texto"`          // descricao legivel (PT-BR)
	Conexo         bool           `json:"conexo"`         // conexo (grafo) / fracamente conexo (digrafo)
	GrausImpares   []string       `json:"grausImpares"`   // grafos: vertices de grau impar
	Desbalanceados map[string]int `json:"desbalanceados"` // digrafos: vertice -> (saida-entrada)
	VerticeInicial string         `json:"verticeInicial"` // ponto de partida sugerido
}

// TrailDTO espelha euleriano.TrilhaEuler + a contagem de arestas percorridas.
type TrailDTO struct {
	Sequencia []string `json:"sequencia"` // cadeia ordenada de vertices (len = arestas + 1)
	Circuito  bool     `json:"circuito"`  // true se volta ao ponto de origem
	Arestas   int      `json:"arestas"`   // ruas percorridas (len(sequencia) - 1)
}

// RouteResponse e a resposta de /api/rota: classificacao + trilha (ou erro).
type RouteResponse struct {
	OK            bool               `json:"ok"`
	Erro          string             `json:"erro"`
	Classificacao *ClassificationDTO `json:"classificacao"`
	Trilha        *TrailDTO          `json:"trilha"`
}
