# Trabalho de Grafos - 2026

## Estrutura

```
cmd/main.go                 → ponto de entrada
internal/grafo/grafo.go     → estrutura do grafo
internal/algoritmos/        → cada membro implementa aqui
internal/conversoes/        → cada membro implementa aqui
internal/leitor/leitor.go   → lê os arquivos da pasta inputs
internal/relatorio/relatorio.go → gera a saída padronizada
```

## Como rodar

```
go build -o projeto ./cmd/main.go
./projeto
```

As saídas são geradas em `outputs/`.

## Como cada membro adiciona sua parte

### Passo 1 — Implemente a função

Abra o arquivo em `internal/conversoes/` ou `internal/algoritmos/` correspondente ao seu item e preencha a lógica.

### Passo 2 — Adicione a formatação ao relatório

Se a saída precisa de um texto formatado, adicione uma função `Formata<Algo>(g)` em `internal/relatorio/relatorio.go`.

### Passo 3 — Chame no main

No `cmd/main.go`, adicione a linha de relatório. Se sua função depender de uma conversão prévia, chame-a antes:

```go
conversoes.ListaParaMatriz(g)
r.Adiciona("MATRIZ_DE_ADJACENCIA", relatorio.FormataMatriz(g))
```

### Passo 4 — Rode e verifique

```
go run ./cmd/main.go
cat outputs/GRAFO_0.txt
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

---

## Funções do grafo:

| Função | O que faz |
|---|---|
| `g.AdicionarVertice(id)` | cria vértice se não existir |
| `g.RemoverVertice(id)` | remove vértice e conexões |
| `g.AdicionarAresta(a, b)` | conecta dois vértices |
| `g.RemoverAresta(a, b)` | remove conexão |
| `g.NumVertices()` | total de vértices |
| `g.NumArestas()` | total de arestas |
| `g.GetVizinhos(id)` | retorna vizinhos de um vértice |

