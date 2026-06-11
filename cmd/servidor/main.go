// Command servidor sobe a aplicacao web do Projeto 3 (EcoUrbano):
// roteirizacao de coleta de residuos via grafos eulerianos. Reusa os pacotes
// internos (leitor + euleriano) e serve um front em SVG embutido.
//
// Uso:
//
//	go run ./cmd/servidor               # http://localhost:8080, datasets de inputs_u2
//	go run ./cmd/servidor -addr :9000   # outra porta
//	go run ./cmd/servidor -inputs dir   # outro diretorio de datasets
package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"

	"github.com/PauloFH/grafos-2026/internal/leitor"
	"github.com/PauloFH/grafos-2026/internal/web"
)

//go:embed web
var assets embed.FS

func main() {
	addr := flag.String("addr", ":8080", "endereco de escuta (ex.: :8080)")
	inputs := flag.String("inputs", "inputs_u2", "diretorio com os datasets .txt")
	flag.Parse()

	datasets, err := leitor.LerDiretorio(*inputs)
	if err != nil {
		log.Fatalf("erro ao carregar datasets de %q: %v", *inputs, err)
	}
	if len(datasets) == 0 {
		log.Fatalf("nenhum dataset .txt encontrado em %q", *inputs)
	}

	static, err := fs.Sub(assets, "web")
	if err != nil {
		log.Fatalf("erro ao abrir assets embutidos: %v", err)
	}

	srv := web.NovoServer(datasets, static)

	log.Printf("EcoUrbano (Projeto 3) — %d dataset(s) carregado(s) de %q", len(datasets), *inputs)
	for nome := range datasets {
		log.Printf("  - %s", nome)
	}
	log.Printf("Servindo em http://localhost%s  (Ctrl+C para parar)", *addr)

	if err := http.ListenAndServe(*addr, srv.Rotas()); err != nil {
		log.Fatal(err)
	}
}
