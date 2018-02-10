package configuration

import "fmt"

type IntMatrix interface {
	At(i, j int) int
	Set(val, i, j int) error
	Dims() (r, c int)
}

func NewIntMatrix(r, c int) IntMatrix {
	return &intMatrix{
		data: make([]int, r*c),
		r:    r,
		c:    c,
	}
}

type intMatrix struct {
	data []int
	r    int
	c    int
}

func (im *intMatrix) Set(val, i, j int) error {
	if i >= im.r {
		return fmt.Errorf("row index out of range")
	}
	if j >= im.c {
		return fmt.Errorf("col index out of range")
	}
	im.data[im.r * i + j] = val
	return nil
}

func (im *intMatrix) At(i, j int) int {
	return im.data[i*im.r+j]
}

func (im *intMatrix) Dims() (r, c int) {
	return im.r, im.c
}
