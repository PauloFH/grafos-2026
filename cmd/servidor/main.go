// Command servidor sobe a aplicacao web unica do trabalho, servindo duas
// unidades: U2 (EcoUrbano, rotas eulerianas) em /u2/ e U3 (Caixeiro
// Viajante) em /u3/, com landing em /. Reusa os pacotes internos
// (leitor + euleriano + pcv) e serve os fronts em SVG embutidos.
//
// Uso:
//
//	go run ./cmd/servidor                  # http://localhost:8080
//	go run ./cmd/servidor -addr :9000      # outra porta
//	go run ./cmd/servidor -inputs-u2 dir   # outro diretorio de datasets da U2
//	go run ./cmd/servidor -inputs-u3 dir   # outro diretorio de instancias do PCV
package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"

	"github.com/PauloFH/grafos-2026/internal/leitor"
	"github.com/PauloFH/grafos-2026/internal/pcv"
	"github.com/PauloFH/grafos-2026/internal/web"
)

//go:embed web
var assets embed.FS

func main() {
	addr := flag.String("addr", ":8080", "endereco de escuta (ex.: :8080)")
	inputsU2 := flag.String("inputs-u2", "inputs_u2", "diretorio com os datasets .txt da U2")
	inputsU3 := flag.String("inputs-u3", "inputs_u3", "diretorio com as instancias .txt do PCV (U3)")
	flag.Parse()

	datasets, err := leitor.LerDiretorio(*inputsU2)
	if err != nil {
		log.Fatalf("erro ao carregar datasets de %q: %v", *inputsU2, err)
	}
	if len(datasets) == 0 {
		log.Fatalf("nenhum dataset .txt encontrado em %q", *inputsU2)
	}

	instancias, err := pcv.CarregaDiretorio(*inputsU3)
	if err != nil {
		log.Fatalf("erro ao carregar instancias de %q: %v", *inputsU3, err)
	}
	if len(instancias) == 0 {
		log.Fatalf("nenhuma instancia .txt encontrada em %q", *inputsU3)
	}

	static, err := fs.Sub(assets, "web")
	if err != nil {
		log.Fatalf("erro ao abrir assets embutidos: %v", err)
	}

	srv := web.NovoServer(datasets, instancias, static)

	log.Printf("U2 EcoUrbano — %d dataset(s) carregado(s) de %q", len(datasets), *inputsU2)
	for nome := range datasets {
		log.Printf("  - %s", nome)
	}
	log.Printf("U3 Caixeiro Viajante — %d instancia(s) carregada(s) de %q", len(instancias), *inputsU3)
	for _, in := range instancias {
		log.Printf("  - %s (%s, N=%d)", in.Nome, in.Medida, in.N)
	}
	log.Printf("Servindo em http://localhost%s  (Ctrl+C para parar)", *addr)

	if err := http.ListenAndServe(*addr, srv.Rotas()); err != nil {
		log.Fatal(err)
	}
}
