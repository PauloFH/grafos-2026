package algoritmos

import "errors"

type nodePilha struct {
	data string
	next *nodePilha
}

type Pilha struct {
	top      *nodePilha
	inserted int
}

func (s *Pilha) Tamanho() int {
	return s.inserted
}

func (s *Pilha) Empilha(e string) {
	newNode := &nodePilha{data: e, next: nil}
	if s.top != nil {
		newNode.next = s.top
	}
	s.top = newNode
	s.inserted++
}

func (s *Pilha) Desempilha() (string, error) {
	if s.top == nil {
		return "", errors.New("pilha vazia")
	}
	dataPop := s.top.data
	s.top = s.top.next
	s.inserted--
	return dataPop, nil
}

func (s *Pilha) Topo() (string, error) {
	if s.top == nil {
		return "", errors.New("pilha vazia")
	}
	return s.top.data, nil
}
