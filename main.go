package main

import (
	"fmt"
	Classes "./classes"
	"math/rand"
	"time"
)


func main(){
	//Each player has a board
	// var board classes.Board
	// board = *classes.NewBoard(10,10)
	// board.Size_x = 10
	// board.Size_y = 10
	
	//Init player
	var player1 Classes.Player
	var computer Classes.Player
	player1, err := player1.Initialize(false)
	if err!= nil {
		fmt.Println(err)
	}
	fmt.Println("This is your ocean grid map:")
	player1.Showgrid()
	
	// Initialize a computer controled player
	computer, err = computer.Initialize(true)
	if err!= nil {
		fmt.Println(err)
	}

	fmt.Println("Computer adversary created")
	fmt.Println("Starting game!")

	fmt.Println("Player attack first!")

	victory := false
	for victory == false{
		fmt.Println("Your turn to strike.")
		player1.Showgrid()
		x,y := shoot(computer.Board.Size_x, computer.Board.Size_y)
		fmt.Printf("You strike at: (%d - %d) \n", x,y)
		player1 = player1.MarkPoints(x,y,verify(computer, x, y))
		if player1.Wins(){
			victory=true
			break
		}

		x,y = machineShoot(player1.Board.Size_x, player1.Board.Size_y)
		fmt.Printf("Computer strikes at: (%d - %d) \n", x,y)
		computer = computer.MarkPoints(x,y,verify(player1, x, y))
		if computer.Wins(){
			victory=true
			break
		}
	}
}

func getRandomInt(ran int) int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(ran)
}

func machineShoot(maxX, maxY int) (int, int){
	var x int
	x = getRandomInt(maxX)
	var y int
	y = getRandomInt(maxY)
	return x,y
}

func shoot(maxX, maxY int) (int, int){
	fmt.Println("Enter X coordinate: ")
	valid := false
	var x int
	for valid == false{
		fmt.Scanf("%d", &x)
		if x >= 0 || x < maxX{
			valid = true
		}else{
			fmt.Printf("Invalid X coordinate. Please shoot between %v and %v\n", 0, maxX)
		}
	}

	fmt.Println("Enter Y coordinate: ")
	var y int
	valid = false
	for valid == false{
		fmt.Scanf("%d", &y)
		if y >= 0 || y < maxY{
			valid = true
		}else{
			fmt.Printf("Invalid Y coordinate. Please shoot between %v and %v\n", 0, maxY)
		}
	}
	return x,y
}

func verify(p Classes.Player, x int, y int) bool{
	r, _:= p.CheckPoints(x, y)
	return r
}