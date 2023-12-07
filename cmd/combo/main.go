package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/timer"
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
	timeout = time.Second * 30
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
	incantate incantate
	spells    []string
	point     int
	timer     timer.Model
}

func main() {
	m := model{
		timer:     timer.NewWithInterval(timeout, time.Second),
		orbs:      make([]string, 3, 3),
		invoked:   make([]spell, 2, 2),
		incantate: make(map[spell]struct{}),
		cast:      true,
	}
	m.cast, m.incantate, m.spells = gen(nil)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q", "w", "e":
			m.orbs[2] = m.orbs[1]
			m.orbs[1] = m.orbs[0]
			m.orbs[0] = strings.ToUpper(msg.String())
		case "r":
			i := invoke(m.orbs)
			if i == Undefined {
				break
			}

			if !m.cast {
				if _, ok := m.incantate[i]; ok {
					delete(m.incantate, i)
					m.spells = reOrder(m.spells, spellMap[i])
				}
			}
			if m.invoked[0] == i {
				break
			}
			m.invoked[1] = m.invoked[0]
			m.invoked[0] = i
		case "d":
			if _, ok := m.incantate[m.invoked[0]]; ok {
				delete(m.incantate, m.invoked[0])
				m.spells = reOrder(m.spells, spellMap[m.invoked[0]])
			}
		case "f":
			if _, ok := m.incantate[m.invoked[1]]; ok {
				delete(m.incantate, m.invoked[1])
				m.spells = reOrder(m.spells, spellMap[m.invoked[0]])
			}
		}
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		return m, tea.Quit
	}

	if len(m.incantate) == 0 {
		if !m.cast {
			m.cast, m.incantate, m.spells = gen(m.incantate)
			m.point++
		} else {
			m.cast, m.incantate, m.spells = gen(nil)
			m.point++
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
	s += strings.Join(m.spells, " | ")
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

	if m.timer.Timeout == 0 {
		s += fmt.Sprintf("You Make %d Combos\n\n", m.point)
	}
	return s
}

func invoke(orbs []string) spell {
	if len(orbs) < 3 {
		return Undefined
	}

	var invokerOrb rune = 1

	orb := strings.Join(orbs, "")

	for _, o := range orb {
		invokerOrb *= o
	}

	return spellValue[invokerOrb]
}

func reOrder(spells []string, answer string) []string {
	newSpells := make([]string, 0, len(spells)-1)

	for _, spell := range spells {
		if spell != answer {
			newSpells = append(newSpells, spell)
		}
	}

	return newSpells
}

func gen(prev incantate) (bool, incantate, []string) {
	cast := true
	r := rand.New(rand.NewSource(time.Now().Unix()))

	next := make(map[spell]struct{})

	spells := make([]string, 0, 2)

	for len(next) < 2 {
		n := spell(1 + r.Intn(10))
		if prev != nil {
			if _, ok := prev[n]; ok {
				continue
			}
		}
		if _, ok := next[n]; ok {
			continue
		}
		next[n] = struct{}{}
		spells = append(spells, spellMap[n])
	}

	return cast, next, spells
}
