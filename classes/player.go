package classes

import (
	"fmt"
	"io/ioutil"
	"log"
	"errors"
	"strconv"
	"strings"
	"gopkg.in/yaml.v2"
	"math/rand"
	"time"
)

// Player structure
type Player struct {
	Board      Board
	// Carrier    Ship
	// Battleship Ship
	// Cruiser    Ship
	// Submarine  Ship
	// Destroyer  Ship
	Ships 	   []Ship
	Machine    bool
	Enemy      Board
	Hits	   int
}

// Game structure
type Game struct{
	Board map[string]int `yaml:"board"`
	Carrier map[string]string `yaml:"Carrier"`
	Battleship map[string]string `yaml:"Battleship"`
	Cruiser map[string]string `yaml:"Cruiser"`
	Submarine map[string]string `yaml:"Submarine"`
	Destroyer map[string]string `yaml:"Destroyer"`
}

func (p Player) validatePoints(ix int, iy int, fx int, fy int, size int) bool {
	if ix >= p.Board.Size_x || iy >= p.Board.Size_y || fx >= p.Board.Size_x || fy >= p.Board.Size_y {
		return false
	}
	if ix == fx {
		return Abs(iy-fy) == size-1
	}
	return Abs(ix-fx) == size-1 && ix == iy
}

func (p Player) generateMissingPoints(i1, i2, f1, f2 int) []Point {
	InitPoint := Point{x: i1, y: i2}
	FinalPoint := Point{x: f1, y: f2}
	listPoints := []Point{InitPoint, FinalPoint}
	if i1 == f1 {
		root := Min(i2, f2)
		for i := 1; i < Abs(i2-f2); i++ {
			listPoints = append(listPoints, Point{x: i1, y: root + i})
		}
	} else {
		root := Min(i1, f1)
		for i := 1; i < Abs(i1-f1); i++ {
			listPoints = append(listPoints, Point{x: root + i, y: i2})
		}
	}
	return listPoints
}

func (p Player) extractPoints(size int, p1 string, p2 string) ([]Point, error){
	var splits = strings.Split(p1, ",")
	i1, _ := strconv.Atoi(strings.Trim(splits[0], " "))
	i2, _ := strconv.Atoi(strings.Trim(splits[1], " "))

	splits = strings.Split(p2, ",")
	f1, _ := strconv.Atoi(strings.Trim(splits[0], " "))
	f2, _ := strconv.Atoi(strings.Trim(splits[1], " "))

	if p.validatePoints(i1, i2, f1, f2, size) {
		fmt.Printf("Correct positions: (%d, %d) - (%d, %d) \n", i1, i2, f1, f2) 
	}else{
		return []Point{}, errors.New(fmt.Sprintf("Not valid points. Check the size (%d, %d) - (%d, %d) \n", i1, i2, f1, f2))
	}
	return p.generateMissingPoints(i1, i2, f1, f2), nil
}

func (p Player) validateShips(s Ship, size int, dimensions map[string]string) (Ship, error){
	s.size = size
	valid := false
	for valid == false {
		var err error
		s.points, err = p.extractPoints(size, dimensions["init"], dimensions["final"])
		if err != nil {
			panic(err)
		}
		r, err := p.ValidShip(s)
		if err != nil{
			panic(err)
		}
		if r{
			return s, nil
		}
	}
	return s, errors.New("Error during validation")
}

func (p Player) generateRandomStringPoints(s Ship, size int) (Ship, error){
	rand.Seed(time.Now().UnixNano())
	valid := false
	for valid ==false{
		initx := rand.Intn(p.Board.Size_x)
		inity := rand.Intn(p.Board.Size_y)
		// Vertical :0
		// horizontal :1
		verhoriz := rand.Intn(2)
		var finalx int
		var finaly int
		if verhoriz > 0{
			finalx = initx
			finaly = inity + size -1
		}else{
			finalx = initx + size -1
			finaly = inity
		}
		if p.validatePoints(initx, inity, finalx, finaly, size){
			fmt.Printf("Correct positions: (%d, %d) - (%d, %d) \n", initx, inity, finalx, finaly) 
			s.points = p.generateMissingPoints(initx, inity, finalx, finaly)
			r, _ := p.ValidShip(s)
			if r {
				valid = true
			}
		}
	}
	return s, nil
}

// Addship adds a ship to the board grid
func (p Player) Addship(s Ship) (bool, error) {
	// fmt.Println("Current points: ", p.Board.Points)
	// p.Board.Points = append(p.Board.Points, s.points...)
	// fmt.Println(p.Board.Points)
	// fmt.Printf("Added ship: %v\n", s.ship)
	
	a := []rune(s.ship)
	for _, point := range s.points {
		p.Board.Grid[point.x][point.y] = string(a[0:2])
	}
	// for i := 0; i < p.Board.Size_x; i++ {
	// 	fmt.Println(p.Board.Grid[i])
	// }
	return true, nil
}

// Initialize function:
// If is a human player, machine is false
// if is a computer controled, machine is true
func (p Player) Initialize(machine bool) (Player, error){
	p.Machine = machine
	data, err := ioutil.ReadFile("defs/game.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Game
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
		return p, err
	}
	fmt.Printf("%+v\n", config)

	p.Board = Board{Size_x: config.Board["x"], Size_y: config.Board["y"]}
	p.Board.Grid = p.Initgrid(config.Board["x"], config.Board["y"], "OO")
	p.Enemy.Grid = p.Initgrid(config.Board["x"], config.Board["y"], "__")
	p.Board.Points = []Point{}
	
	listShips := []string{
		"Carrier", "Battleship", "Cruiser", "Submarine", "Destroyer",
	}
	for _, ship := range listShips {
		// p.Showgrid()
		var s Ship
		s.ship = ship
		switch ship {
		case "Carrier":
			p = p.processShip(s, 5, config.Carrier)
		case "Battleship":
			p = p.processShip(s, 4, config.Battleship)
		case "Cruiser":
			p = p.processShip(s, 3, config.Cruiser)
		case "Submarine":
			p = p.processShip(s, 3, config.Submarine)
		case "Destroyer":
			p = p.processShip(s, 2, config.Destroyer)
		}
	}
	fmt.Println("Ships in map completed.")
	return p, nil
}

func (p Player) processShip(s Ship, size int, config map[string]string) Player{
	fmt.Printf("Processing player machine: %v \n", p.Machine)
	if !p.Machine{
		s, _ = p.validateShips(s, size, config)
	}else{
		s, _ = p.generateRandomStringPoints(s, size)
	}	
	p.Ships = append(p.Ships, s)
	p.Board.Points = append(p.Board.Points, s.points...)
	p.Addship(s)
	return p
}

// Showgrid shows the grid and the ships on it
func (p Player) Showgrid() {
	for i := 0; i < p.Board.Size_x; i++ {
		fmt.Printf("%v \t %v\n", p.Board.Grid[i], p.Enemy.Grid[i])
	}
}

// ValidShip verifies if a ship does not position over other ships
func (p Player) ValidShip(s Ship) (bool, error){
	fmt.Printf("Verifying position of %v\n", s.ship)
	for _, bpoint := range p.Board.Points {
		for _, shipPos := range s.points {
			// fmt.Printf("Points: (%v,%v),(%v,%v)", bpoint.x, shipPos.x, bpoint.y, shipPos.y)
			if bpoint.x == shipPos.x && bpoint.y == shipPos.y {
				fmt.Println("Invalid ship position")
				return false, errors.New("Invalid ship position. Rewrite positions.")
			}
		}
	}
	fmt.Println("Valid ship position")
	return true, nil
}

// CheckPoints checks for a hit
func (p Player) CheckPoints(x int, y int) (bool, error){
	for _, bpoint := range p.Board.Points {
		if bpoint.x == x && bpoint.y == y {
			fmt.Println("Your attack hit!")
			return true, nil
		}
	}
	fmt.Println("No hit")
	return false, nil
}

// MarkPoints marks a hit or a miss
func (p Player) MarkPoints(x int, y int, hit bool) Player{
	if hit{
		p.Enemy.Grid[x][y] = "XX"
		p.Hits++
	}else{
		p.Enemy.Grid[x][y] = "OO"
	}
	return p
}

// Wins check if all ships have been defeated
// 17 points is all is needed
func (p Player) Wins() bool{
	if p.Hits >= 17{
		return true
	}
	return false
}

// Initgrid initializes a grid
func (p Player) Initgrid(x, y int, symb string) [][]string {
	grid := make([][]string, x)
	for i := 0; i < x; i++ {
		row := make([]string, y)
		for j := 0; j < y; j++ {
			row[j] = symb
		}
		grid[i] = row
	}
	return grid
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Min returns the minimum value of two.
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
