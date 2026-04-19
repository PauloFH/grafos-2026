package algoritmos

import "fmt"

type nodeFila struct {
	data string
	next *nodeFila
}

type Fila struct {
	head *nodeFila
	tail *nodeFila
	size int
}

func (q *Fila) Tamanho() int {
	return q.size
}

func (q *Fila) Enfileira(value string) {
	newNode := &nodeFila{data: value}
	if q.tail != nil {
		q.tail.next = newNode
	}
	q.tail = newNode
	if q.head == nil {
		q.head = newNode
	}
	q.size++
}

func (q *Fila) Desenfileira() (string, error) {
	if q.head == nil {
		return "", fmt.Errorf("fila vazia")
	}
	value := q.head.data
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	q.size--
	return value, nil
}

func (q *Fila) Frente() (string, error) {
	if q.head == nil {
		return "", fmt.Errorf("fila vazia")
	}
	return q.head.data, nil
}
