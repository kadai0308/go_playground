package main

import (
	"fmt"
	"math/rand"
	"time"
)

var cardTypeEnum = map[int]string{
	1: "Club",
	2: "Diamond",
	3: "Heart",
	4: "Spade",
}

type Desk struct {
	order []int
}

type Player struct {
	handCard   []int
	playedCard []int
}

func (d *Desk) newDeck() {
	newDeck := make([]int, 54)
	for i := range newDeck {
		newDeck[i] = i + 1
	}
	d.order = newDeck
	d.deskShuffle()
}

func (d *Desk) deskShuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(d.order), func(i, j int) {
		d.order[i], d.order[j] = d.order[j], d.order[i]
	})
}

func (d *Desk) convertCardInfo(cardNo int) (string, int) {
	cardType := cardNo / 12
	cardNum := cardNo - cardType*13
	return cardTypeEnum[cardType+1], cardNum
}

func (d *Desk) dealDesk(cardAmount int) []int {
	returnValue := d.order[len(d.order)-cardAmount:]
	d.order = d.order[:len(d.order)-cardAmount]
	return returnValue
}

func main() {
	d := Desk{}
	d.newDeck()
	fmt.Println(d.order)
	fmt.Println(d.convertCardInfo(d.order[0]))
	fmt.Println(d.dealDesk(5))
}
