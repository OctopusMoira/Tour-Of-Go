package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	slice := make([][]uint8, dy)
	for y := range slice {
		xslice := make([]uint8, dx)
		for x := range xslice{
			xslice[x] = uint8(x*y)
		}
		slice[y] = xslice
	}
	return slice
}

func main() {
	pic.Show(Pic)
}

