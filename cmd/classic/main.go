package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type spell int

const (
	Undefined spell = iota
	ColdSnap
	GhostWalk
	IceWall
	Tornado
	EMP
	Alacrity
	SunStrike
	ForgeSpirit
	ChaosMeteor
	DeafeningBlast
)

var spellMap = map[spell]string{
	Undefined:      "",
	ColdSnap:       "Cold Snap",
	GhostWalk:      "Ghost Walk",
	IceWall:        "Ice Wall",
	Tornado:        "Tornado",
	EMP:            "EMP",
	Alacrity:       "Alacrity",
	SunStrike:      "Sun Strike",
	ForgeSpirit:    "Forge Spirit",
	ChaosMeteor:    "Chaos Meteor",
	DeafeningBlast: "Deafening Blast",
}

// spellValue
// map invoker orb spells to invoker skills/spells.

// the reason we use multiplication (*) here instead of addition (+) is that the
// result wouldn't be unique with addition.
// So, in this case, we opt for multiplication to ensure a unique value.
var spellValue = map[rune]spell{
	'Q' * 'Q' * 'Q': ColdSnap,
	'Q' * 'Q' * 'W': GhostWalk,
	'Q' * 'Q' * 'E': IceWall,
	'W' * 'W' * 'W': EMP,
	'W' * 'W' * 'Q': Tornado,
	'W' * 'W' * 'E': Alacrity,
	'E' * 'E' * 'E': SunStrike,
	'E' * 'E' * 'Q': ForgeSpirit,
	'E' * 'E' * 'W': ChaosMeteor,
	'Q' * 'W' * 'E': DeafeningBlast,
}

type incantate map[spell]struct{}

type model struct {
	orbs      []string
	invoked   []spell
	cast      bool
	spell     spell
	point     int
	timeSince time.Time
	record    float64
}

func main() {
	m := model{
		orbs:      make([]string, 3, 3),
		invoked:   make([]spell, 2, 2),
		timeSince: time.Now(),
	}
	m.spell = generate(0)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// NOTe
	// to guard if the point is already 10, immediately quit
	if m.point == 10 {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q", "w", "e": // handle upper case?
			m.orbs[2] = m.orbs[1]
			m.orbs[1] = m.orbs[0]
			m.orbs[0] = strings.ToUpper(msg.String())
		case "r":
			i := invoke(m.orbs)
			if i == Undefined {
				break
			}

			m.invoked[1] = m.invoked[0]
			m.invoked[0] = i

			if i == m.spell {
				m.spell = generate(i)
				m.point++
			}

			if m.point == 10 {
				m.record = time.Since(m.timeSince).Seconds()
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string

	if !m.cast {
		s += "INVOKE\n\n"
	} else {
		s += "INVOKE AND CAST\n\n"
	}

	s += "| "
	s += spellMap[m.spell]
	s += " |"
	s += "\n--------------------------------------------------\n"

	s += "| "
	s += strings.Join(m.orbs, " | ")
	s += " |"
	s += "\n--------------------------------------------------\n"

	s += "| Q | W | E |"
	for _, i := range m.invoked {
		s += " "
		s += spellMap[i]
		s += " |"
	}
	s += "\n--------------------------------------------------\n"

	s += fmt.Sprintf("%d POINTS\n\n", m.point)

	if m.point == 10 {
		s += "Your injoker classic Record is:\n"
		s += fmt.Sprintf("%v Seconds\n\n", m.record)
	}

	return s
}

func invoke(orbs []string) spell {
	if len(orbs) < 3 {
		return Undefined
	}

	// NOTE
	// need to assign initial value to 1 (default 0)
	// because in below we need do operation using multiplication
	// to determine what type of spell
	var invokerOrb rune = 1

	orb := strings.Join(orbs, "")

	for _, o := range orb {
		invokerOrb *= o
	}

	return spellValue[invokerOrb]
}

func generate(prevSpell spell) spell {
	randSource := rand.NewSource(time.Now().Unix())
	r := rand.New(randSource)

	randSpellNumber := r.Intn(10-1) + 1

	if randSpellNumber == int(prevSpell) {
		return generate(prevSpell)
	}

	return spell(randSpellNumber)
}
