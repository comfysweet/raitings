package main

import (
	"github.com/comfysweet/ratings/domain/model"
	"github.com/emirpasic/gods/maps/treemap"
)

func main() {
	storage := model.PlayerStorage{
		PlayerById:      map[int]*model.Player{},
		PlayersByPoints: treemap.NewWith(IntComparator),
		Ratings:         model.Ratings{Places: []int{}}}

	storage.ChangeRating(1, 1)
	storage.PrintStorage()
	storage.ChangeRating(2, 0)
	storage.PrintStorage()
	storage.ChangeRating(2, 3)
	storage.PrintStorage()
	storage.ChangeRating(3, 1)
	storage.PrintStorage()
	storage.ChangeRating(3, 0)
	storage.PrintStorage()
	storage.ChangeRating(3, 1)
	storage.PrintStorage()
}

func IntComparator(a, b interface{}) int {
	aAsserted := a.(int)
	bAsserted := b.(int)
	switch {
	case aAsserted < bAsserted:
		return 1
	case aAsserted > bAsserted:
		return -1
	default:
		return 0
	}
}
