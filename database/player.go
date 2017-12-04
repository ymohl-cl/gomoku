package database

import "fmt"

func (P Player) Play() {
	fmt.Println("Player play !")
}

func CreatePlayer(name string) *Player {
	p := new(Player)

	p.Name = name
	p.Saves = -1
	return p
}

func (P *Player) DeleteSave() string {
	if P.Saves == -1 {
		return "No session saved"
	}
	P.Saves = -1
	return ""
}
