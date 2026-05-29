# Trabalho de Grafos - 2026

Projeto da disciplina **DIM0549 Grafos** (UFRN). O repositório reúne as duas unidades,
compartilhando o mesmo núcleo em `internal/`:

- **Unidade 1** — representações, conversões e algoritmos clássicos sobre grafos/dígrafos.
- **Unidade 2** — **Projeto 3: Grafos Eulerianos** (validação de eulerianidade + Hierholzer).

## Estrutura do projeto

```
cmd/
  unidade1/main.go                 → ponto de entrada da Unidade 1
  unidade2/main.go                 → ponto de entrada da Unidade 2 (eulerianos)

internal/                          → NÚCLEO COMPARTILHADO pelas duas unidades
  grafo/grafo.go                   → estrutura e métodos do grafo
  leitor/leitor.go                 → lê os arquivos de inputs/
  algoritmos/                      → BFS, DFS, conexo, contagem, pilha, etc.
  conversoes/                      → conversões entre representações
  relatorio/relatorio.go           → relatório textual padronizado
  relatorio/png*.go                → geração de imagens PNG (Graphviz)
  euleriano/                       → Unidade 2: eulerianidade + Hierholzer

inputs/      outputs/              → entradas/saídas da Unidade 1
inputs_u2/   outputs_u2/           → entradas/saídas da Unidade 2
```

---

## Build e execução

### Pré-requisitos

- [Go 1.26+](https://go.dev/dl/)
- [Graphviz](https://graphviz.org/download/) (apenas para gerar os PNGs)
  - Arch/CachyOS: `sudo pacman -S graphviz`
  - Sem o Graphviz o programa roda normal e só pula a geração de imagens.

### Com o Makefile (recomendado)

```bash
make           # lista os alvos disponíveis
make u1        # roda a Unidade 1            -> outputs/
make u2        # roda a Unidade 2 (eulerianos) -> outputs_u2/
make run       # roda as duas em sequência
make build     # compila as duas em bin/
make check     # go build ./... + go vet ./... (antes de commitar)
make fmt       # gofmt em todo o código
make clean     # remove binários e .dot gerados
```

### Sem o Makefile (go puro)

> Rode **sempre a partir da raiz** do repositório — os caminhos `inputs*/`/`outputs*/`
> são relativos ao diretório atual.

```bash
# Unidade 1
go run ./cmd/unidade1
go build -o bin/unidade1 ./cmd/unidade1 && ./bin/unidade1

# Unidade 2 (eulerianos)
go run ./cmd/unidade2
go build -o bin/unidade2 ./cmd/unidade2 && ./bin/unidade2
```

### Verificar saídas

```bash
cat outputs/GRAFO_1.txt
cat outputs_u2/GRAFO_EULER.txt
```

---

## Formato dos arquivos de entrada

Cada grafo é um `.txt` em `inputs/` (Unidade 1) ou `inputs_u2/` (Unidade 2):

```
20            # 1ª linha: número de vértices (informativo)
1,2           # demais linhas: origem,destino
1,5
2,3
...
```

O tipo é decidido pelo **nome do arquivo**: se contém `DIGRAFO`, é tratado como dígrafo
(arestas direcionadas); caso contrário, como grafo não-direcionado.

---

# Unidade 2 — Projeto 3: Grafos Eulerianos

Cenário: roteamento de coleta de resíduos que precisa percorrer **todas as ruas exatamente
uma vez** (trilha/circuito euleriano). O software deve:

1. **Validar a eulerianidade** do grafo e do dígrafo.
2. **Algoritmo de Hierholzer (grafos)** — extrair a cadeia euleriana (arestas bidirecionais).
3. **Algoritmo de Hierholzer (dígrafos)** — mesma extração com arestas direcionadas.

### Como rodar

```bash
make u2          # ou: go run ./cmd/unidade2
cat outputs_u2/GRAFO_EULER.txt
```

Saídas em `outputs_u2/`:
- `<NOME>.txt` — relatório com as seções `EULERIANIDADE` e `TRILHA_EULERIANA`
- `<NOME>.png` — visualização do grafo
- `<NOME>_TRILHA.png` — trilha euleriana com as arestas numeradas na ordem de percurso

### Arquitetura (pacote `internal/euleriano`)

Contrato compartilhado (em `eulerianidade.go`) que liga validação ↔ Hierholzer ↔ relatório:

| Tipo / Função | Papel |
|---|---|
| `ResultadoEuler` | classe (não/semi/euleriano), conectividade, graus ímpares, vértice inicial |
| `TrilhaEuler` | sequência ordenada de vértices + se é circuito |
| `Classifica(g)` | aplica os critérios de eulerianidade (grafo ou dígrafo) |
| `HierholzerGrafo(g, inicio)` | cadeia euleriana de grafo não-direcionado |
| `HierholzerDigrafo(g, inicio)` | cadeia euleriana de dígrafo |

Reuso do núcleo: `grafo.Clone/RemoverAresta`, `leitor`, `algoritmos.EhConexo/BFS/Pilha`,
`relatorio.Relatorio/GeradorPNG`.

### Divisão de tarefas (Unidade 2)

| Responsável | Tarefa | Arquivo |
|---|---|---|

---

# Unidade 1

### Estrutura do Grafo

```go
type Grafo struct {
    NomeArquivo string
    Direcionado bool
    Vertices    []string            // vértices na ordem de leitura
    ListaAdj    map[string][]string // vértice -> vizinhos
}
```

### Métodos disponíveis

| Método | O que faz |
|---|---|
| `g.AdicionarVertice(id)` | cria vértice se não existir |
| `g.RemoverVertice(id)` | remove vértice e todas as suas conexões |
| `g.AdicionarAresta(a, b)` | conecta dois vértices (bidirecional se não-direcionado) |
| `g.RemoverAresta(a, b)` | remove conexão |
| `g.GetVizinhos(id)` | retorna slice de vizinhos |
| `g.GrauVertice(id)` | grau de saída do vértice |
| `g.GrausVertices()` | map com grau de saída de todos os vértices |
| `g.Clone()` | retorna cópia independente do grafo |

### Funções utilitárias (`internal/algoritmos/`)

| Função | O que faz |
|---|---|
| `algoritmos.TotalVertices(g)` | número de vértices |
| `algoritmos.TotalArestas(g)` | número de arestas |
| `algoritmos.SaoAdjacentes(g, a, b)` | verifica adjacência entre dois vértices |
| `algoritmos.ParesAdjacentes(g)` | lista todos os pares adjacentes |
| `algoritmos.EhConexo(g)` | conectividade |
| `algoritmos.BFS(g, inicio)` | busca em largura |
| `algoritmos.DFS(g, inicio)` | busca em profundidade (não-direcionado) |
| `algoritmos.DFSDigrafo(g)` | DFS com classificação de arestas (dígrafo) |
| `algoritmos.Biconectividade(g)` | articulações e blocos via lowpt |
| `algoritmos.Bipartido(g)` | verifica bipartição |
| `algoritmos.DeterminaGrafoSubjacente(g)` | grafo subjacente de um dígrafo |
| `algoritmos.EstrelaDireta(matriz, vertices)` | converte para estrela direta |

### Divisão de tarefas (Unidade 1)

| # | Descrição | Grafos | Responsável |
|---|---|---|---|
| 1 | Representação por Lista de Adjacências | GRAFO1, GRAFO2 | Paulo Roberto |
| 2 | Representação por Matriz de Adjacências | GRAFO1, GRAFO2 | Paulo Roberto |
| 3 | Representação por Matriz de Incidência | GRAFO1, GRAFO2 | Vinicius |
| 4 | Conversão Matriz de Adj. ↔ Lista de Adj. | GRAFO1, GRAFO2 | Paulo Roberto |
| 5 | Calcular o grau de cada vértice | GRAFO1, GRAFO2 | João Marcelo |
| 6 | Determinar se dois vértices são adjacentes | GRAFO1, GRAFO2 | Paulo Roberto |
| 7 | Determinar número total de vértices | GRAFO1, GRAFO2 | Vinicius |
| 8 | Determinar número total de arestas | GRAFO1, GRAFO2 | Vinicius |
| 9 | Inclusão de um novo vértice | GRAFO1 | João Marcelo |
| 10 | Exclusão de um vértice existente | GRAFO1 | João Marcelo |
| 11 | Determinar se o grafo é conexo | GRAFO1, GRAFO2 | Vinicius |
| 12 | OPCIONAL: Determinar se é bipartido (1,0 pt) | GRAFO1, GRAFO2 | Diego |
| 13 | Busca em Largura (BFS) | GRAFO1, GRAFO3 | Diego |
| 14 | Busca em Profundidade (DFS) | GRAFO1, GRAFO3 | Diego |
| 15 | Articulações e Blocos (Biconectividade via lowpt) | GRAFO3 | Diego |
| 16 | Representação por Matriz de Adjacências | DIGRAFO1, DIGRAFO2 | Paulo Roberto |
| 17 | Representação por Matriz de Incidência | DIGRAFO1, DIGRAFO2 | Vinicius |
| 18 | OPCIONAL: Determinação do Grafo subjacente (0,5 pt) | DIGRAFO1 | João Victor |
| 19 | OPCIONAL: Matriz Incidência ↔ Estrela Direta (0,5 pt) | DIGRAFO1 | João Victor |
| 20 | DFS (Profundidade entrada/saída e tipos de arestas) | DIGRAFO2, DIGRAFO3 | João Victor |
| 21 | OPCIONAL: Aplicação Real de DFS (1,0 pt) | Exemplo do grupo (≥ 10 vértices) | João Victor |
