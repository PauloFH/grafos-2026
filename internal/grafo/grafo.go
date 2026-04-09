package grafo

// Vertice representa um nó do grafo
type Vertice struct {
	ID string
}

// Grafo representa um grafo ou digrafo
// Armazena os dados em lista de adjacência
type Grafo struct {
	NomeArquivo string
	Direcionado bool
	Vertices    []string            // lista de vértices na ordem
	ListaAdj    map[string][]string // vértice -> vizinhos
}

// NovoGrafo cria um grafo vazio
func NovoGrafo(direcionado bool, nome string) *Grafo {
	return &Grafo{
		NomeArquivo: nome,
		Direcionado: direcionado,
		Vertices:    make([]string, 0),
		ListaAdj:    make(map[string][]string),
	}
}

// AdicionarVertice inclui um vértice se ele ainda não existir
// Item 9 - Responsável: João Marcelo
func (g *Grafo) AdicionarVertice(id string) {
	// Verifica se já existe
	for _, v := range g.Vertices {
		if v == id {
			return
		}
	}

	// Adiciona
	g.Vertices = append(g.Vertices, id)
	g.ListaAdj[id] = make([]string, 0)
}

// RemoverVertice exclui um vértice e suas conexões
// Item 10 - Responsável: João Marcelo
func (g *Grafo) RemoverVertice(id string) {
	// Remove das listas dos outros vértices
	for v, vizinhos := range g.ListaAdj {
		if v == id {
			continue
		}
		novos := make([]string, 0)
		for _, viz := range vizinhos {
			if viz != id {
				novos = append(novos, viz)
			}
		}
		g.ListaAdj[v] = novos
	}

	// Remove a lista e o vértice
	delete(g.ListaAdj, id)

	novos := make([]string, 0)
	for _, v := range g.Vertices {
		if v != id {
			novos = append(novos, v)
		}
	}
	g.Vertices = novos
}

// AdicionarAresta conecta dois vértices
// Cria os vértices se não existirem
func (g *Grafo) AdicionarAresta(origem, destino string) {
	g.AdicionarVertice(origem)
	g.AdicionarVertice(destino)

	g.ListaAdj[origem] = append(g.ListaAdj[origem], destino)

	if !g.Direcionado {
		g.ListaAdj[destino] = append(g.ListaAdj[destino], origem)
	}
}

// RemoverAresta remove a conexão entre dois vértices
func (g *Grafo) RemoverAresta(origem, destino string) {
	// Remove origem -> destino
	if vizinhos, ok := g.ListaAdj[origem]; ok {
		novos := make([]string, 0)
		for _, v := range vizinhos {
			if v != destino {
				novos = append(novos, v)
			}
		}
		g.ListaAdj[origem] = novos
	}

	// Se não direcionado, remove destino -> origem também
	if !g.Direcionado {
		if vizinhos, ok := g.ListaAdj[destino]; ok {
			novos := make([]string, 0)
			for _, v := range vizinhos {
				if v != origem {
					novos = append(novos, v)
				}
			}
			g.ListaAdj[destino] = novos
		}
	}
}

// NumVertices retorna o total de vértices
func (g *Grafo) NumVertices() int {
	return len(g.Vertices)
}

// NumArestas retorna o total de arestas
func (g *Grafo) NumArestas() int {
	total := 0
	for _, vizinhos := range g.ListaAdj {
		total += len(vizinhos)
	}
	if !g.Direcionado {
		total /= 2
	}
	return total
}

// SaoAdjacentes verifica se dois vértices são vizinhos
func (g *Grafo) SaoAdjacentes(a, b string) bool {
	vizinhos, ok := g.ListaAdj[a]
	if !ok {
		return false
	}
	for _, v := range vizinhos {
		if v == b {
			return true
		}
	}
	return false
}

// GetVizinhos retorna os vizinhos de um vértice
func (g *Grafo) GetVizinhos(id string) []string {
	return g.ListaAdj[id]
}
