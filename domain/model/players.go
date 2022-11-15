package model

import "fmt"

type Player struct {
	Id     int
	Points int // количество очков
	Place  int
}

type Ratings struct {
	Places []int // место в турнирной таблице = индекс массива и его, id игрока = значение
}

func (p Player) String() string {
	return fmt.Sprintf("{Id: %d Points: %d Place: %d}", p.Id, p.Points, p.Place)
}
