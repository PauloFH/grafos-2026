package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PauloFH/grafos-2026/internal/pcv"
)

// sementePadrao e a semente usada em /api/u3/resolve quando o query param
const sementePadrao int64 = 42

func (s *Server) handleInstanciasU3(w http.ResponseWriter, _ *http.Request) {
	infos := make([]InstanciaInfo, 0, len(s.instancias))
	for _, in := range s.instancias {
		infos = append(infos, InstanciaInfo{ID: in.Nome, Medida: in.Medida, N: in.N})
	}
	escreveJSON(w, http.StatusOK, infos)
}

func (s *Server) handleInstanciaU3(w http.ResponseWriter, r *http.Request) {
	in, ok := s.instanciaPorNome[r.URL.Query().Get("id")]
	if !ok {
		escreveErro(w, http.StatusNotFound, "instancia nao encontrada")
		return
	}
	pontos := PosicoesMDS(in)
	if pontos == nil {
		pontos = []PontoDTO{}
	}
	escreveJSON(w, http.StatusOK, InstanciaResponse{
		ID:       in.Nome,
		Medida:   in.Medida,
		N:        in.N,
		Deposito: in.Cidades[0],
		Pontos:   pontos,
		Otimo:    pcv.ValoresOtimos[in.Nome],
	})
}

func (s *Server) handleResolveU3(w http.ResponseWriter, r *http.Request) {
	in, ok := s.instanciaPorNome[r.URL.Query().Get("id")]
	if !ok {
		escreveErro(w, http.StatusNotFound, "instancia nao encontrada")
		return
	}

	metodo := r.URL.Query().Get("metodo")

	semente := sementePadrao
	if bruta := r.URL.Query().Get("semente"); bruta != "" {
		v, err := strconv.ParseInt(bruta, 10, 64)
		if err != nil {
			escreveJSON(w, http.StatusOK, falhaResolveU3(metodo, "semente invalida"))
			return
		}
		semente = v
	}

	switch metodo {
	case "vmp":
		// deterministico (semente ignorada): Antes = construtivo (VMP),
		// Rota = apos a busca local 2-opt.
		inicio := time.Now()
		antes := pcv.VizinhoMaisProximo{}.Constroi(in)
		rota := pcv.DoisOpt{}.Aplica(antes, in)
		tempoMs := float64(time.Since(inicio).Nanoseconds()) / 1e6
		escreveJSON(w, http.StatusOK, SolveResponse{
			OK:         true,
			Metodo:     metodo,
			Antes:      &RotaDTO{Ordem: antes.IdsOriginais(in), Custo: antes.Custo},
			Rota:       RotaDTO{Ordem: rota.IdsOriginais(in), Custo: rota.Custo},
			GapPercent: pcv.Gap(rota.Custo, pcv.ValoresOtimos[in.Nome]),
			TempoMs:    tempoMs,
		})
	case "imb":
		// deterministico (semente ignorada): Antes = construtivo (IMB),
		// Rota = apos a busca local Or-opt.
		inicio := time.Now()
		antes := pcv.InsercaoMaisBarata{}.Constroi(in)
		rota := pcv.OrOpt{}.Aplica(antes, in)
		tempoMs := float64(time.Since(inicio).Nanoseconds()) / 1e6
		escreveJSON(w, http.StatusOK, SolveResponse{
			OK:         true,
			Metodo:     metodo,
			Antes:      &RotaDTO{Ordem: antes.IdsOriginais(in), Custo: antes.Custo},
			Rota:       RotaDTO{Ordem: rota.IdsOriginais(in), Custo: rota.Custo},
			GapPercent: pcv.Gap(rota.Custo, pcv.ValoresOtimos[in.Nome]),
			TempoMs:    tempoMs,
		})
	case "ag":
		// Uma execucao do AG com a semente pedida
		ag := pcv.AlgoritmoGenetico{Par: pcv.ParametrosPadrao()}
		inicio := time.Now()
		rota := ag.Executa(in, semente)
		tempoMs := float64(time.Since(inicio).Nanoseconds()) / 1e6
		escreveJSON(w, http.StatusOK, SolveResponse{
			OK:         true,
			Metodo:     metodo,
			Antes:      nil,
			Rota:       RotaDTO{Ordem: rota.IdsOriginais(in), Custo: rota.Custo},
			GapPercent: pcv.Gap(rota.Custo, pcv.ValoresOtimos[in.Nome]),
			TempoMs:    tempoMs,
		})
	case "memetico":
		memetico := pcv.AlgoritmoMemetico{Par: pcv.ParametrosMemeticoPadrao()}
		inicio := time.Now()
		rota := memetico.Executa(in, semente)
		tempoMs := float64(time.Since(inicio).Nanoseconds()) / 1e6
		escreveJSON(w, http.StatusOK, SolveResponse{
			OK:         true,
			Metodo:     metodo,
			Antes:      nil,
			Rota:       RotaDTO{Ordem: rota.IdsOriginais(in), Custo: rota.Custo},
			GapPercent: pcv.Gap(rota.Custo, pcv.ValoresOtimos[in.Nome]),
			TempoMs:    tempoMs,
		})
	default:
		escreveJSON(w, http.StatusOK, falhaResolveU3(metodo, "metodo desconhecido: "+metodo))
	}
}

// falhaResolveU3 monta um SolveResponse de erro de negocio
func falhaResolveU3(metodo, erro string) SolveResponse {
	return SolveResponse{
		OK:     false,
		Erro:   erro,
		Metodo: metodo,
		Rota:   RotaDTO{Ordem: []int{}},
	}
}
