// throw away
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TODO can set non default button
// useful for legacy keys user
type button string

const (
	q button = "q"
	w        = "w"
	e        = "e"
	r        = "r"
	d        = "d"
	f        = "f"
)

type spell string

// TODO use Invoker responses
const (
	Undefined      spell = "N O O B"
	ColdSnap             = "Cold Snap"
	GhostWalk            = "Ghost Walk"
	IceWall              = "Ice Wall"
	Tornado              = "Tornado"
	EMP                  = "EMP"
	Alacrity             = "Alacrity"
	SunStrike            = "Sun Strike"
	ForgeSpirit          = "Forge Spirit"
	ChaosMeteor          = "Chaos Meteor"
	DeafeningBlast       = "Deafening Blast"
)

func main() {
	// TODO game choice
	orbs := make([]button, 3, 3)
	invoked := make([]spell, 2, 2)
	s := bufio.NewScanner(os.Stdin)
	s.Buffer([]byte{}, 2)
	for {
		var in string
		// TODO use tcell/bubbletea
		// am too stoopid to figure out how to listen
		// to key events myself
		if s.Scan() {
			in = s.Text()
		}
		if err := s.Err(); err != nil {
			panic(err)
		}
		// not in map might be cheaper than branching?
		// loop with goto might be clearer?
		// maybe use byte operations?
		if strings.EqualFold(in, string(q)) ||
			strings.EqualFold(in, string(w)) ||
			strings.EqualFold(in, string(e)) {
			orbs[2] = orbs[1]
			orbs[1] = orbs[0]
			orbs[0] = button(in)
		} else if strings.EqualFold(in, string(r)) {
			i := invoke(orbs)
			if i != Undefined {
				invoked[1] = invoked[0]
				invoked[0] = i
			}
		} else if strings.EqualFold(in, string(d)) {
			fmt.Println(invoked[0])
		} else if strings.EqualFold(in, string(f)) {
			fmt.Println(invoked[1])
		}
		// TODO randomized combo
	}
}

func invoke(orbs []button) spell {
	var count int
	combo := make(map[button]int)
	for _, o := range orbs {
		if o == "" {
			continue
		}
		combo[o]++
		count++
	}
	if count < 3 {
		return Undefined
	}
	if combo[q] == 3 {
		return ColdSnap
	} else if combo[w] == 3 {
		return EMP
	} else if combo[e] == 3 {
		return SunStrike
	} else if combo[q] == 2 {
		if combo[w] == 1 {
			return GhostWalk
		} else if combo[e] == 1 {
			return IceWall
		}
	} else if combo[w] == 2 {
		if combo[q] == 1 {
			return Tornado
		} else if combo[e] == 1 {
			return Alacrity
		}
	} else if combo[e] == 2 {
		if combo[q] == 1 {
			return ForgeSpirit
		} else if combo[w] == 1 {
			return ChaosMeteor
		}
	}
	return DeafeningBlast
}
