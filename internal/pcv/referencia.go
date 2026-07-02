package pcv

// ValoresOtimos mapeia Instancia.Nome para o melhor valor conhecido do
// artigo GRASP_SBPO2025 (Tabelas 2 e 4); usado no calculo do gap.
// Consultar com "v, ok :=" - instancia ausente e responsabilidade do chamador.
var ValoresOtimos = map[string]float64{
	"PROBLEMA_01": 1942.30,
	"PROBLEMA_02": 1973.00,
	"PROBLEMA_03": 1705.60,
	"PROBLEMA_04": 1664.00,
	"PROBLEMA_05": 1321.00,
	"PROBLEMA_06": 1223.00,
	"PROBLEMA_07": 672.70,
	"PROBLEMA_08": 606.00,
	"PROBLEMA_09": 438.30,
	"PROBLEMA_10": 364.00,
	"PROBLEMA_11": 344.90,
	"PROBLEMA_12": 305.00,
}

// Gap devolve o desvio percentual (custo-otimo)/otimo*100, ou 0 se
// otimo <= 0. Pode ser negativo (custo < otimo indica bug e queremos ve-lo).
func Gap(custo, otimo float64) float64 {
	if otimo <= 0 {
		return 0
	}
	return (custo - otimo) / otimo * 100
}
