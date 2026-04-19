# Trabalho de Grafos - 2026

## Estrutura do projeto

```
cmd/main.go                      → ponto de entrada
internal/grafo/grafo.go          → estrutura e métodos do grafo
internal/algoritmos/             → implementações de algoritmos
internal/conversoes/             → conversões entre representações
internal/leitor/leitor.go        → lê os arquivos de inputs/
internal/relatorio/relatorio.go  → gera saída padronizada em texto
internal/relatorio/png.go        → geração de imagens PNG (Graphviz)
inputs/                          → arquivos de entrada (.txt por grafo)
outputs/                         → relatórios e imagens gerados
```

---

## Build e execução

### Pré-requisitos

- [Go 1.26+](https://go.dev/dl/)
- [Graphviz](https://graphviz.org/download/) (para geração de PNGs)

### Compilar e executar

```bash
# Executar direto (sem compilar)
go run ./cmd/main.go

# Compilar e executar o binário
go build -o projeto_grafos_unidade_1 ./cmd/main.go
./projeto_grafos_unidade_1
```

### Verificar saídas

```bash
cat outputs/GRAFO_1.txt
cat outputs/DIGRAFO_DO_GRUPO.txt
```

Os arquivos gerados em `outputs/` são:
- `<NOME>.txt` — relatório textual completo
- `<NOME>.png` — visualização do grafo
- `<NOME>_BFS.png` — árvore BFS (GRAFO_1, GRAFO_3, DIGRAFO_DO_GRUPO)
- `<NOME>_DFS.png` — árvore DFS (GRAFO_1, GRAFO_3, DIGRAFO_DO_GRUPO)
- `GRAFO_1_ADD_VERTEX.png` — grafo após inclusão de vértice
- `GRAFO_1_REMOVE_VERTEX.png` — grafo após exclusão de vértice

---

## Como cada membro adiciona sua parte

### Passo 1 — Implemente a função

Crie ou edite um arquivo em `internal/algoritmos/` ou `internal/conversoes/`.

### Passo 2 — Adicione a formatação ao relatório

Se a saída precisa de texto formatado, adicione `Formata<Algo>(g)` em `internal/relatorio/relatorio.go` ou `formata_algoritmos.go`.

### Passo 3 — Chame no main

Em `cmd/main.go`, adicione dentro da função correspondente (`adicionarDadosBasicos`, `processarDigrafo`, etc.):

```go
r.Adiciona("TITULO_DA_SECAO", relatorio.FormataAlgo(g))
```

### Passo 4 — Rode e verifique

```bash
go run ./cmd/main.go
cat outputs/GRAFO_1.txt
```
---

## Estrutura do Grafo

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
| `algoritmos.EhConexo(g)` | conectividade  |
| `algoritmos.BFS(g, inicio)` | busca em largura |
| `algoritmos.DFS(g, inicio)` | busca em profundidade (não-direcionado) |
| `algoritmos.DFSDigrafo(g)` | DFS com classificação de arestas (dígrafo) |
| `algoritmos.Biconectividade(g)` | articulações e blocos via lowpt |
| `algoritmos.Bipartido(g)` | verifica bipartição |
| `algoritmos.DeterminaGrafoSubjacente(g)` | grafo subjacente de um dígrafo |
| `algoritmos.EstrelaDireta(matriz, vertices)` | converte para estrela direta |

---

## Divisão de tarefas

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
