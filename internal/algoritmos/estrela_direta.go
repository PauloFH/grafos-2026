package algoritmos

import (
	"fmt"
	"strings"

	"github.com/PauloFH/grafos-2026/internal/grafo"
)

// Cria a matriz a partir do grafo atual
func MatrizIncidencia(g *grafo.Grafo) ([][]int, []string) {
	numVertices := len(g.Vertices)
	numArestas := g.NumArestas()

	// Inicializa a matriz linhas x colunas
	matriz := make([][]int, numVertices)
	for i := range matriz {
		matriz[i] = make([]int, numArestas)
	}

	arestaIdx := 0
	for _, u := range g.Vertices {
		for _, v := range g.ListaAdj[u] {
			// Encontra o índice dos vértices
			uIdx := IndexFinder(g.Vertices, u)
			vIdx := IndexFinder(g.Vertices, v)

			//Calculando os graus de entrada e saída para cada vértice
			matriz[uIdx][arestaIdx] = -1 // Saída
			matriz[vIdx][arestaIdx] = 1  // Entrada
			arestaIdx++
		}
	}
	return matriz, g.Vertices
}

// Recebe a matriz e converte para Estrela Direta
func EstrelaDireta(matriz [][]int, vertices []string) (A []string, IP []int) {
	numV := len(vertices)

	IP = make([]int, numV+1) // Vetor de ponteiros
	if len(matriz) == 0 {
		return []string{}, IP
	}

	numE := len(matriz[0])
	if numE == 0 {
		return []string{}, IP
	}

	A = make([]string, 0, numE) // Vetor de sucessores

	ponteiro := 0
	for i := 0; i < numV; i++ {
		IP[i] = ponteiro
		// Para cada linha (vértice), procuramos as arestas onde ele é origem (-1)
		for j := 0; j < numE; j++ {
			if matriz[i][j] == -1 {
				// Encontramos uma aresta saindo de 'i'.
				// Agora buscamos quem é o destino (onde está o '1' nessa coluna j)
				for k := 0; k < numV; k++ {
					if matriz[k][j] == 1 {
						A = append(A, vertices[k])
						ponteiro++
						break
					}
				}
			}
		}
	}
	IP[numV] = ponteiro
	return A, IP
}

// Helper para encontrar índice de string no slice
func IndexFinder(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

// FormatarEstrelaDireta gera a string para o relatório
func FormatarEstrelaDireta(A []string, IP []int, vertices []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  Vetor A (Sucessores): [%s]\n", strings.Join(A, ", ")))

	strIP := make([]string, len(IP))
	for i, val := range IP {
		strIP[i] = fmt.Sprintf("%d", val)
	}
	sb.WriteString(fmt.Sprintf("  Vetor IP (Ponteiros): [%s]\n", strings.Join(strIP, ", ")))
	return sb.String()
}
