package classes

type Point struct{
	x int
	y int
}

type Ship struct{
	ship string
	size int
	points []Point
}

type Board struct{
	Size_x int
	Size_y int
	Points []Point
	Grid [][]string
}



