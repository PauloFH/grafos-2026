package web

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"sort"
	"strconv"

	"github.com/PauloFH/grafos-2026/internal/algoritmos"
	"github.com/PauloFH/grafos-2026/internal/euleriano"
	"github.com/PauloFH/grafos-2026/internal/grafo"
	"github.com/PauloFH/grafos-2026/internal/pcv"
)

// Server serve as APIs da U2 (EcoUrbano, /api/u2/*) e da U3 (Caixeiro
// Viajante, /api/u3/*) + os arquivos estaticos dos fronts. Os dados sao
// carregados uma vez na inicializacao e ficam imutaveis aqui: o Hierholzer
// clona internamente e os algoritmos do PCV trabalham sobre copias de Rota,
// entao chamadas concorrentes sao seguras.
type Server struct {
	datasets         map[string]*grafo.Grafo
	instancias       []*pcv.Instancia
	instanciaPorNome map[string]*pcv.Instancia
	static           fs.FS
}

// NovoServer cria o servidor a partir dos datasets da U2, das instancias do
// PCV (na ordem de pcv.CarregaDiretorio) e do FS estatico; devolve o Server
// pronto para servir via Rotas().
func NovoServer(datasets map[string]*grafo.Grafo, instancias []*pcv.Instancia, static fs.FS) *Server {
	porNome := make(map[string]*pcv.Instancia, len(instancias))
	for _, in := range instancias {
		porNome[in.Nome] = in
	}
	return &Server{
		datasets:         datasets,
		instancias:       instancias,
		instanciaPorNome: porNome,
		static:           static,
	}
}

// Rotas devolve o handler HTTP com a API e o servidor de estaticos.
func (s *Server) Rotas() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/u2/datasets", s.handleDatasets)
	mux.HandleFunc("GET /api/u2/graph", s.handleGraph)
	mux.HandleFunc("GET /api/u2/rota", s.handleRota)
	mux.HandleFunc("GET /api/u3/instancias", s.handleInstanciasU3)
	mux.HandleFunc("GET /api/u3/instancia", s.handleInstanciaU3)
	mux.HandleFunc("GET /api/u3/resolve", s.handleResolveU3)
	mux.Handle("/", http.FileServerFS(s.static))
	return mux
}

// GET /api/datasets -> lista resumida dos grafos disponiveis.
func (s *Server) handleDatasets(w http.ResponseWriter, _ *http.Request) {
	infos := make([]DatasetInfo, 0, len(s.datasets))
	for id, g := range s.datasets {
		infos = append(infos, DatasetInfo{
			ID:          id,
			Nome:        g.NomeArquivo,
			Direcionado: g.Direcionado,
			Vertices:    algoritmos.TotalVertices(g),
			Arestas:     algoritmos.TotalArestas(g),
		})
	}
	// Grafos antes de digrafos, depois por id, para uma ordem estavel na UI.
	sort.Slice(infos, func(i, j int) bool {
		if infos[i].Direcionado != infos[j].Direcionado {
			return !infos[i].Direcionado
		}
		return infos[i].ID < infos[j].ID
	})
	escreveJSON(w, http.StatusOK, infos)
}

// GET /api/graph?dataset=ID -> malha viaria desenhavel.
func (s *Server) handleGraph(w http.ResponseWriter, r *http.Request) {
	g, ok := s.datasets[r.URL.Query().Get("dataset")]
	if !ok {
		escreveErro(w, http.StatusNotFound, "dataset nao encontrado")
		return
	}
	escreveJSON(w, http.StatusOK, grafoParaDTO(r.URL.Query().Get("dataset"), g))
}

// GET /api/rota?dataset=ID&inicio=V -> classificacao + trilha euleriana.
func (s *Server) handleRota(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("dataset")
	g, ok := s.datasets[id]
	if !ok {
		escreveErro(w, http.StatusNotFound, "dataset nao encontrado")
		return
	}

	res := euleriano.Classifica(g)
	classe := classificacaoParaDTO(res)
	resp := RouteResponse{Classificacao: &classe}

	if res.Classe == euleriano.NaoEuleriano {
		resp.OK = false
		resp.Erro = "Grafo nao-euleriano: nao existe trilha que percorra cada rua exatamente uma vez."
		escreveJSON(w, http.StatusOK, resp)
		return
	}

	inicio := r.URL.Query().Get("inicio")
	var trilha euleriano.TrilhaEuler
	var err error
	if g.Direcionado {
		trilha, err = euleriano.HierholzerDigrafo(g, inicio)
	} else {
		trilha, err = euleriano.HierholzerGrafo(g, inicio)
	}
	if err != nil {
		resp.OK = false
		resp.Erro = err.Error()
		escreveJSON(w, http.StatusOK, resp)
		return
	}

	t := trilhaParaDTO(trilha)
	resp.OK = true
	resp.Trilha = &t
	escreveJSON(w, http.StatusOK, resp)
}

// --- mapeadores dominio -> DTO ---

func grafoParaDTO(id string, g *grafo.Grafo) GraphResponse {
	return GraphResponse{
		ID:          id,
		Nome:        g.NomeArquivo,
		Direcionado: g.Direcionado,
		Vertices:    verticesOrdenados(g.Vertices),
		Arestas:     arestasDTO(g),
		Graus:       g.GrausVertices(),
	}
}

// arestasDTO extrai as ruas. Para grafo nao-direcionado a lista de adjacencia
// guarda cada aresta nos dois sentidos, entao deduplicamos pelo par canonico.
func arestasDTO(g *grafo.Grafo) []ArestaDTO {
	arestas := make([]ArestaDTO, 0)
	vistos := make(map[string]bool)
	for _, v := range verticesOrdenados(g.Vertices) {
		for _, w := range g.ListaAdj[v] {
			if g.Direcionado {
				arestas = append(arestas, ArestaDTO{De: v, Para: w})
				continue
			}
			lo, hi := parCanonico(v, w)
			chave := lo + "\x00" + hi
			if vistos[chave] {
				continue
			}
			vistos[chave] = true
			arestas = append(arestas, ArestaDTO{De: lo, Para: hi})
		}
	}
	return arestas
}

func classificacaoParaDTO(res euleriano.ResultadoEuler) ClassificationDTO {
	impares := res.GrausImpares
	if impares == nil {
		impares = []string{}
	}
	desb := res.Desbalanceados
	if desb == nil {
		desb = map[string]int{}
	}
	return ClassificationDTO{
		Classe:         classeChave(res.Classe),
		Texto:          res.Classe.String(),
		Conexo:         res.Conexo,
		GrausImpares:   impares,
		Desbalanceados: desb,
		VerticeInicial: res.VerticeInicial,
	}
}

func trilhaParaDTO(t euleriano.TrilhaEuler) TrailDTO {
	seq := t.Sequencia
	if seq == nil {
		seq = []string{}
	}
	arestas := 0
	if len(seq) > 0 {
		arestas = len(seq) - 1
	}
	return TrailDTO{Sequencia: seq, Circuito: t.Circuito, Arestas: arestas}
}

// classeChave devolve uma chave estavel para o front (sem espacos/acentos).
func classeChave(c euleriano.Classe) string {
	switch c {
	case euleriano.CircuitoEuleriano:
		return "CircuitoEuleriano"
	case euleriano.CaminhoEuleriano:
		return "CaminhoEuleriano"
	default:
		return "NaoEuleriano"
	}
}

// verticesOrdenados ordena por valor inteiro quando possivel (1,2,...,10,11),
// caindo para ordem lexicografica se algum rotulo nao for numerico.
func verticesOrdenados(vs []string) []string {
	out := make([]string, len(vs))
	copy(out, vs)
	sort.Slice(out, func(i, j int) bool {
		return menorRotulo(out[i], out[j])
	})
	return out
}

func parCanonico(a, b string) (string, string) {
	if menorRotulo(a, b) {
		return a, b
	}
	return b, a
}

func menorRotulo(a, b string) bool {
	na, erra := strconv.Atoi(a)
	nb, errb := strconv.Atoi(b)
	if erra == nil && errb == nil {
		return na < nb
	}
	return a < b
}

// --- helpers HTTP ---

func escreveJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func escreveErro(w http.ResponseWriter, status int, msg string) {
	escreveJSON(w, status, map[string]string{"erro": msg})
}
