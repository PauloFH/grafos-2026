package web

import (
	"math"

	"github.com/PauloFH/grafos-2026/internal/pcv"
)

// PosicoesMDS projeta as cidades de uma instancia do PCV no plano via MDS classico sobre a matriz de custos.
func PosicoesMDS(in *pcv.Instancia) []PontoDTO {
	n := in.N
	xs, ys, ok := coordenadasMDS(in)
	if !ok || degenerado(xs, ys) {
		xs, ys = coordenadasCirculo(n)
	}
	normalizaIntervalo(xs)
	normalizaIntervalo(ys)

	pontos := make([]PontoDTO, n)
	for k := 0; k < n; k++ {
		pontos[k] = PontoDTO{ID: in.Cidades[k], Nome: in.Nomes[k], X: xs[k], Y: ys[k]}
	}
	return pontos
}

func coordenadasMDS(in *pcv.Instancia) (xs, ys []float64, ok bool) {
	n := in.N
	if n < 2 {
		return nil, nil, false
	}

	d2 := make([][]float64, n)
	mediaLinha := make([]float64, n)
	mediaTotal := 0.0
	for i := 0; i < n; i++ {
		d2[i] = make([]float64, n)
		soma := 0.0
		for j := 0; j < n; j++ {
			v := in.Custo(i, j)
			d2[i][j] = v * v
			soma += d2[i][j]
		}
		mediaLinha[i] = soma / float64(n)
		mediaTotal += soma
	}
	mediaTotal /= float64(n * n)

	b := make([][]float64, n)
	for i := 0; i < n; i++ {
		b[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			b[i][j] = -0.5 * (d2[i][j] - mediaLinha[i] - mediaLinha[j] + mediaTotal)
		}
	}

	const epsPositivo = 1e-9
	var lambdas [2]float64
	var vetores [2][]float64
	achados := 0
	for tentativa := 0; tentativa < n && achados < 2; tentativa++ {
		lambda, v := autoparDominante(b)
		if math.Abs(lambda) < epsPositivo {
			break
		}
		if lambda > 0 {
			lambdas[achados] = lambda
			vetores[achados] = append([]float64(nil), v...)
			achados++
		}
		deflaciona(b, lambda, v)
	}
	if achados < 2 {
		return nil, nil, false
	}

	xs = make([]float64, n)
	ys = make([]float64, n)
	e1 := math.Sqrt(lambdas[0])
	e2 := math.Sqrt(lambdas[1])
	for i := 0; i < n; i++ {
		xs[i] = vetores[0][i] * e1
		ys[i] = vetores[1][i] * e2
	}
	return xs, ys, true
}

func autoparDominante(b [][]float64) (float64, []float64) {
	n := len(b)
	v := make([]float64, n)
	for i := range v {
		v[i] = 1 + 0.001*float64(i)
	}
	normalizaVetor(v)

	lambda := 0.0
	w := make([]float64, n)
	for it := 0; it < 100; it++ {
		// w := B * v
		for i := 0; i < n; i++ {
			soma := 0.0
			for j := 0; j < n; j++ {
				soma += b[i][j] * v[j]
			}
			w[i] = soma
		}
		norma := normaVetor(w)
		if norma < 1e-15 {
			return 0, v
		}
		for i := range w {
			w[i] /= norma
		}
		novo := 0.0
		for i := 0; i < n; i++ {
			soma := 0.0
			for j := 0; j < n; j++ {
				soma += b[i][j] * w[j]
			}
			novo += w[i] * soma
		}
		copy(v, w)
		if math.Abs(novo-lambda) < 1e-9 {
			lambda = novo
			break
		}
		lambda = novo
	}
	return lambda, v
}

func deflaciona(b [][]float64, lambda float64, v []float64) {
	n := len(b)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			b[i][j] -= lambda * v[i] * v[j]
		}
	}
}

func coordenadasCirculo(n int) (xs, ys []float64) {
	xs = make([]float64, n)
	ys = make([]float64, n)
	for i := 0; i < n; i++ {
		ang := 2 * math.Pi * float64(i) / float64(n)
		xs[i] = math.Cos(ang)
		ys[i] = math.Sin(ang)
	}
	return xs, ys
}

func degenerado(xs, ys []float64) bool {
	for i := range xs {
		if math.IsNaN(xs[i]) || math.IsInf(xs[i], 0) || math.IsNaN(ys[i]) || math.IsInf(ys[i], 0) {
			return true
		}
	}
	const eps = 1e-12
	return amplitude(xs) < eps && amplitude(ys) < eps
}

func amplitude(v []float64) float64 {
	menor, maior := v[0], v[0]
	for _, x := range v {
		menor = math.Min(menor, x)
		maior = math.Max(maior, x)
	}
	return maior - menor
}

func normalizaIntervalo(v []float64) {
	menor, maior := v[0], v[0]
	for _, x := range v {
		menor = math.Min(menor, x)
		maior = math.Max(maior, x)
	}
	largura := maior - menor
	if largura < 1e-12 {
		for i := range v {
			v[i] = 0.5
		}
		return
	}
	for i := range v {
		v[i] = 0.05 + 0.9*(v[i]-menor)/largura
	}
}

func normalizaVetor(v []float64) {
	n := normaVetor(v)
	if n < 1e-15 {
		return
	}
	for i := range v {
		v[i] /= n
	}
}

func normaVetor(v []float64) float64 {
	soma := 0.0
	for _, x := range v {
		soma += x * x
	}
	return math.Sqrt(soma)
}
