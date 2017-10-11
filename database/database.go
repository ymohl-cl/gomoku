package database

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/ymohl-cl/gomoku/conf"
)

const (
	Bot = "AI"
)

func Get() (*Data, error) {
	D := new(Data)

	f, err := os.OpenFile(conf.ProtoBufFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadFile(conf.ProtoBufFile)
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(buf, D)
	if err != nil {
		fmt.Println("Unmarshal")
		return nil, err
	}

	err = D.initSave()

	return D, err
}

func (D *Data) initSave() error {
	fmt.Println("Init save protobuf")
	if len(D.Players) == 0 {
		AI := CreatePlayer(Bot)
		D.Players = append(D.Players, AI)
	}

	D.Current = new(Session)
	return nil
}

func (D *Data) UpdateCurrent(p *Player, position int) error {
	if position == 0 {
		D.Current.P1 = p
		if D.Current.P2 != nil && D.Current.P2.Name == p.Name {
			D.Current.P2 = nil
		}
	} else if position == 1 {
		D.Current.P2 = p
		if D.Current.P1 != nil && D.Current.P1.Name == p.Name {
			D.Current.P2 = nil
		}
	} else {
		return errors.New("Position to player is not allowed")
	}

	return nil
}

func (d *Data) AddPlayer(p *Player) {
	d.Players = append(d.Players, p)
	d.Write()
}

func (d *Data) DeletePlayer(p *Player) (int, error) {
	fmt.Println("len player: ", len(d.Players))
	for id, pt := range d.Players {
		if pt.Name == p.Name {
			d.Players = append(d.Players[:id], d.Players[id+1:]...)
			//d.Current = nil
			if d.Current.P1 == p {
				d.Current.P1 = nil
			} else if d.Current.P2 == p {
				d.Current.P2 = nil
			}

			for _, i := range p.Saves {
				d.Sessions = append(d.Sessions[:i], d.Sessions[i+1:]...)
			}

			d.Write()
			//			d.Current = new(Session)
			fmt.Println("len player: ", len(d.Players))
			return id, nil
		}
	}
	return 0, errors.New("Player name not found")
}

func (d *Data) GetPlayerByName(name string) (*Player, error) {
	for _, p := range d.Players {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, errors.New("Player not found")
}

// SaveSession, save the current session
func (d *Data) SaveSession() {
	d.Sessions = append(d.Sessions, d.Current)

	for index, s := range d.Sessions {
		if s == d.Current {
			d.Current.P1.Saves = append(d.Current.P1.Saves, int32(index))
		}
	}

	d.Current = nil
	d.Write()
	d.Current = new(Session)
}

func (d *Data) Write() {
	// get the buffer generated by protobuf
	buf, err := proto.Marshal(d)
	if err != nil {
		panic(err)
	}

	// save on the file
	if err := ioutil.WriteFile(conf.ProtoBufFile, buf, 0755); err != nil {
		panic(err)
	}
}
