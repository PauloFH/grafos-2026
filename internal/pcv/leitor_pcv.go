package pcv

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// CarregaInstancia le e valida um arquivo de instancia no formato canonico:
//
//	Linha 1: ID da instancia (ex.: PROBLEMA_01)
//	Linha 2: MEDIDA: KM|MIN
//	Linha 3: N: <n>
//	Linha 4: CIDADES: <n ids> (primeiro = deposito)
//	Linha 5: NOMES: <n nomes separados por ;>
//	Linha 6: MATRIZ:
//	Linhas 7..: n linhas com n floats
//
// O Nome vem da linha 1, NAO do nome do arquivo. Valida contagens, valores
// nao negativos, diagonal zero e simetria (tolerancia 1e-9).
func CarregaInstancia(caminho string) (*Instancia, error) {
	f, err := os.Open(caminho)
	if err != nil {
		return nil, fmt.Errorf("abrindo %s: %w", caminho, err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	// leLinha devolve a proxima linha (com TrimSpace) e se houve linha.
	numLinha := 0
	leLinha := func() (string, bool) {
		if !sc.Scan() {
			return "", false
		}
		numLinha++
		return strings.TrimSpace(sc.Text()), true
	}

	// Linha 1: id da instancia.
	nome, ok := leLinha()
	if !ok || nome == "" {
		return nil, fmt.Errorf("%s: linha 1: id da instancia vazio", caminho)
	}

	// Linha 2: MEDIDA: KM ou MEDIDA: MIN.
	linha, ok := leLinha()
	medida := strings.TrimSpace(strings.TrimPrefix(linha, "MEDIDA:"))
	if !ok || !strings.HasPrefix(linha, "MEDIDA:") || (medida != "KM" && medida != "MIN") {
		return nil, fmt.Errorf("%s: linha 2: esperado 'MEDIDA: KM' ou 'MEDIDA: MIN', obtido %q", caminho, linha)
	}

	// Linha 3: N: <n>.
	linha, ok = leLinha()
	if !ok || !strings.HasPrefix(linha, "N:") {
		return nil, fmt.Errorf("%s: linha 3: N invalido: %q", caminho, linha)
	}
	n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(linha, "N:")))
	if err != nil || n <= 0 {
		return nil, fmt.Errorf("%s: linha 3: N invalido: %q", caminho, linha)
	}

	// Linha 4: CIDADES: <n ids>.
	linha, ok = leLinha()
	if !ok || !strings.HasPrefix(linha, "CIDADES:") {
		return nil, fmt.Errorf("%s: linha 4: esperado 'CIDADES: <ids>', obtido %q", caminho, linha)
	}
	campos := strings.Fields(strings.TrimPrefix(linha, "CIDADES:"))
	if len(campos) != n {
		return nil, fmt.Errorf("%s: linha 4: esperados %d ids em CIDADES, obtidos %d", caminho, n, len(campos))
	}
	cidades := make([]int, n)
	for i, campo := range campos {
		id, err := strconv.Atoi(campo)
		if err != nil {
			return nil, fmt.Errorf("%s: linha 4: id de cidade invalido: %q", caminho, campo)
		}
		cidades[i] = id
	}

	// Linha 5: NOMES: <n nomes separados por ponto-e-virgula>.
	linha, ok = leLinha()
	if !ok || !strings.HasPrefix(linha, "NOMES:") {
		return nil, fmt.Errorf("%s: linha 5: esperado 'NOMES: <nomes>', obtido %q", caminho, linha)
	}
	nomes := strings.Split(strings.TrimSpace(strings.TrimPrefix(linha, "NOMES:")), ";")
	for i := range nomes {
		nomes[i] = strings.TrimSpace(nomes[i])
	}
	if len(nomes) != n {
		return nil, fmt.Errorf("%s: linha 5: esperados %d nomes em NOMES, obtidos %d", caminho, n, len(nomes))
	}

	// Linha 6: MATRIZ:.
	linha, ok = leLinha()
	if !ok || linha != "MATRIZ:" {
		return nil, fmt.Errorf("%s: linha 6: esperado 'MATRIZ:', obtido %q", caminho, linha)
	}

	// Linhas 7..7+n-1: n linhas com n floats cada.
	matriz := make([][]float64, 0, n)
	for i := 0; i < n; i++ {
		linha, ok = leLinha()
		if !ok {
			return nil, fmt.Errorf("%s: matriz incompleta: esperadas %d linhas, obtidas %d", caminho, n, i)
		}
		campos := strings.Fields(linha)
		if len(campos) != n {
			return nil, fmt.Errorf("%s: linha %d: esperadas %d colunas na matriz, obtidas %d", caminho, numLinha, n, len(campos))
		}
		valores := make([]float64, n)
		for j, campo := range campos {
			v, err := strconv.ParseFloat(campo, 64)
			if err != nil {
				return nil, fmt.Errorf("%s: linha %d, coluna %d: valor invalido: %q", caminho, numLinha, j+1, campo)
			}
			if v < 0 {
				return nil, fmt.Errorf("%s: linha %d, coluna %d: custo negativo: %.2f", caminho, numLinha, j+1, v)
			}
			valores[j] = v
		}
		matriz = append(matriz, valores)
	}

	// Validacao global: diagonal zero e simetria (tolerancia 1e-9).
	for i := 0; i < n; i++ {
		if matriz[i][i] != 0 {
			return nil, fmt.Errorf("%s: matriz: diagonal [%d][%d] deveria ser 0.00, vale %.2f", caminho, i, i, matriz[i][i])
		}
		for j := i + 1; j < n; j++ {
			if math.Abs(matriz[i][j]-matriz[j][i]) > 1e-9 {
				return nil, fmt.Errorf("%s: matriz assimetrica em [%d][%d]: %.2f != %.2f", caminho, i, j, matriz[i][j], matriz[j][i])
			}
		}
	}

	return &Instancia{
		Nome:    nome,
		Medida:  medida,
		N:       n,
		Cidades: cidades,
		Nomes:   nomes,
		Matriz:  matriz,
	}, nil
}

// CarregaDiretorio carrega todas as instancias .txt do diretorio, ordenadas
// pelo NOME DO ARQUIVO (garante PROBLEMA_01..12). O primeiro erro aborta
// tudo; diretorio sem nenhum .txt e erro.
func CarregaDiretorio(dir string) ([]*Instancia, error) {
	entradas, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("lendo diretorio %s: %w", dir, err)
	}

	var arquivos []string
	for _, e := range entradas {
		if e.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(e.Name()), ".txt") {
			arquivos = append(arquivos, e.Name())
		}
	}
	if len(arquivos) == 0 {
		return nil, fmt.Errorf("nenhum .txt encontrado em %s", dir)
	}
	sort.Strings(arquivos)

	instancias := make([]*Instancia, 0, len(arquivos))
	for _, nome := range arquivos {
		in, err := CarregaInstancia(filepath.Join(dir, nome))
		if err != nil {
			return nil, err
		}
		instancias = append(instancias, in)
	}
	return instancias, nil
}
