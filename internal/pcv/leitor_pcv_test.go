package pcv

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

// instanciaExemplo e um fixture 3x3 valido no formato canonico (subconjunto
// real da aba Km: ANGICOS, SAO RAFAEL, ITAJA).
const instanciaExemplo = `PROBLEMA_EXEMPLO
MEDIDA: KM
N: 3
CIDADES: 1 2 3
NOMES: ANGICOS;SAO RAFAEL;ITAJA
MATRIZ:
0.00 38.80 29.70
38.80 0.00 31.50
29.70 31.50 0.00
`

// TestCarregaInstancia grava um fixture valido em t.TempDir(), carrega e
// confere os campos parseados e o custo do ciclo 0->1->2->0.
func TestCarregaInstancia(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "PROBLEMA_EXEMPLO_KM_03.txt")
	if err := os.WriteFile(caminho, []byte(instanciaExemplo), 0o644); err != nil {
		t.Fatal(err)
	}

	in, err := CarregaInstancia(caminho)
	if err != nil {
		t.Fatal(err)
	}

	if in.Nome != "PROBLEMA_EXEMPLO" {
		t.Errorf("Nome = %q, esperado PROBLEMA_EXEMPLO (linha 1, nao o nome do arquivo)", in.Nome)
	}
	if in.Medida != "KM" || in.N != 3 {
		t.Errorf("Medida/N = %q/%d, esperado KM/3", in.Medida, in.N)
	}
	if len(in.Cidades) != 3 || in.Cidades[0] != 1 || in.Cidades[2] != 3 {
		t.Errorf("Cidades = %v, esperado [1 2 3]", in.Cidades)
	}
	if len(in.Nomes) != 3 || in.Nomes[1] != "SAO RAFAEL" {
		t.Errorf("Nomes = %v, esperado [ANGICOS SAO RAFAEL ITAJA]", in.Nomes)
	}
	if in.Matriz[0][1] != 38.80 || in.Matriz[2][1] != 31.50 {
		t.Errorf("Matriz[0][1]/[2][1] = %.2f/%.2f, esperado 38.80/31.50", in.Matriz[0][1], in.Matriz[2][1])
	}

	// Ciclo 0->1->2->0: 38.80 + 31.50 + 29.70 = 100.00.
	ordem := []int{0, 1, 2}
	if err := ValidaRota(ordem, in); err != nil {
		t.Fatal(err)
	}
	if custo := CustoRota(ordem, in); math.Abs(custo-100.00) > 1e-9 {
		t.Errorf("CustoRota = %.2f, esperado 100.00", custo)
	}
}

// TestCarregaInstancia_Malformada confere que um arquivo fora do formato
// canonico (medida invalida) retorna erro em vez de instancia.
func TestCarregaInstancia_Malformada(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "ruim.txt")
	conteudo := "PROBLEMA_RUIM\nMEDIDA: HORAS\nN: 3\n"
	if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := CarregaInstancia(caminho); err == nil {
		t.Error("esperado erro para medida invalida, obtido nil")
	}
}
