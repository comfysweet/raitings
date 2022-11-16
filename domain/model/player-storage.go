package model

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"sync"
)

type PlayerStorage struct {
	PlayerById      map[int]*Player
	PlayersByPoints *treemap.Map
	Ratings         Ratings
	Lock            sync.RWMutex
}

func (storage *PlayerStorage) PrintStorage() {
	fmt.Printf("GetPlaces: %v \n", storage.GetPlaces(0, len(storage.PlayerById)))
}

// O(1)
func (storage *PlayerStorage) GetPlace(id int) int {
	storage.Lock.RLock()
	defer storage.Lock.RUnlock()

	player, ok := storage.PlayerById[id]
	if !ok {
		return 0
	}
	return (*player).Place
}

// O(1)
func (storage *PlayerStorage) GetPlaces(firstPlace int, lastPlace int) []Player {
	storage.Lock.RLock()
	defer storage.Lock.RUnlock()

	players := make([]Player, 0, lastPlace-firstPlace)
	ratings := storage.Ratings.Places[firstPlace:lastPlace]
	for _, id := range ratings {
		players = append(players, *storage.PlayerById[id])
	}
	return players
}

// O(N)
func (storage *PlayerStorage) ChangeRating(id int, points int) {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	_, ok := storage.PlayerById[id]
	if !ok {
		storage.addPlayer(id, points)
	} else {
		storage.updatePlayer(id, points)
	}
	storage.updatePlaces()
}

// O(logN)
func (storage *PlayerStorage) addPlayer(id int, points int) {
	player := &Player{Id: id, Points: points}
	storage.PlayerById[id] = player

	existPlayers, found := storage.PlayersByPoints.Get(points)
	if found {
		players := existPlayers.([]*Player)
		storage.putToStorage(players, player, points)
	} else {
		players := make([]*Player, 0)
		storage.putToStorage(players, player, points)
	}
}

// O(logN)
func (storage *PlayerStorage) updatePlayer(id int, newPoints int) {
	player := storage.PlayerById[id]
	existPlayers, _ := storage.PlayersByPoints.Get(player.Points)
	players := existPlayers.([]*Player)
	for i, pl := range players {
		if pl.Id == id {
			players = remove(players, i)
			break
		}
	}
	storage.PlayersByPoints.Put(player.Points, players)

	player.Points = newPoints
	existPlayersWithNewPoints, found := storage.PlayersByPoints.Get(newPoints)
	if !found {
		storage.PlayersByPoints.Put(newPoints, []*Player{player})
	} else {
		storage.putToStorage(existPlayersWithNewPoints.([]*Player), player, newPoints)
	}
}

/** TODO:
можно обновлять рейтинг не всех игроков, а только тех, у кого он поменялся

Первоначальный рейтинг
[ a:200, b:30, c:10, f:5, e:1 ]

Обновленное значение для f
f:20
[ a:200, b:30, f:20, c:10, e:1]

oldPoints = 5, newPoints:20

O(k)   k << n - количество обнволенных игроков
Если между собой игроки играют близкие по уровню, то им можно пренебречь. Применимо, когда N гораздо больше чем k

Получаем нужный индекс (3 у f:5), далее нужно идти в сторону изменения рейтинга
если увеличился - в начало списка, если уменьшился - в конец
пока не дойдем до места, рейтинг которого не будет меняться (b:30)
*/
// O(N)
func (storage *PlayerStorage) updatePlaces() {
	ratings := Ratings{}
	for i, value := range storage.PlayersByPoints.Values() {
		players := value.([]*Player)
		for _, player := range players {
			player.Place = i
			ratings.Places = append(ratings.Places, player.Id)
		}
	}
	storage.Ratings = ratings
}

func (storage *PlayerStorage) putToStorage(players []*Player, player *Player, points int) {
	players = append(players, player)
	storage.PlayersByPoints.Put(points, players)
}

func remove(slice []*Player, s int) []*Player {
	return append(slice[:s], slice[s+1:]...)
}
