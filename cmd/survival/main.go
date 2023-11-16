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

// TODO use Invoker responses
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
	timeout = time.Second * 10
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

// model??? really??? in 2023???
// think of a better name
// admittedly this is a copy paste from tutorial
type model struct {
	orbs    []string
	invoked []spell
	cast    bool
	// The reason we still keep the `incantate` map is to prevent unnecessary
	// reordering of the `spells`(slice of string) each time
	// a casting is triggered (when `d` and `f` are invoked).
	incantate  incantate
	spells     []string
	timer      timer.Model
	timeTaken  time.Time
	timeRecord float64
}

func main() {
	m := model{
		timer:     timer.NewWithInterval(timeout, time.Second),
		orbs:      make([]string, 3, 3),
		invoked:   make([]spell, 2, 2),
		incantate: make(map[spell]struct{}),
		timeTaken: time.Now(),
	}
	m.cast, m.incantate, m.spells = gen(nil)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	// TODO init spells to incantate
	// can be 2 or 3 spells
	// can require casting or not
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// if stage only requires invoking, we can still count?
			if !m.cast {
				if _, ok := m.incantate[i]; ok {
					delete(m.incantate, i)
					m.timer.Timeout = timeout
					m.spells = reOrder(m.spells, spellMap[i])
				}
			}
			if m.invoked[0] == i {
				break
			}
			m.invoked[1] = m.invoked[0]
			m.invoked[0] = i
		case "d":
			// Invoke here to prevent unnecessary reordering of spells every time a user triggers
			// a spell casting.
			if _, ok := m.incantate[m.invoked[0]]; ok {
				delete(m.incantate, m.invoked[0])
				m.timer.Timeout = timeout
				m.spells = reOrder(m.spells, spellMap[m.invoked[0]])
			}
		case "f":
			if _, ok := m.incantate[m.invoked[1]]; ok {
				delete(m.incantate, m.invoked[1])
				m.timer.Timeout = timeout
				m.spells = reOrder(m.spells, spellMap[m.invoked[0]])
			}
		}
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		m.timeRecord = time.Since(m.timeTaken).Seconds()
		return m, tea.Quit
	}

	if len(m.incantate) == 0 {
		if !m.cast {
			m.cast, m.incantate, m.spells = gen(m.incantate)
		} else {
			m.cast, m.incantate, m.spells = gen(nil)
		}
	}
	return m, nil
}

// TODO allocate slice
// use strings builder
// nice box
// !!! center div Kappa
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

	s += m.timer.View()

	if m.timer.Timeout == 0 {
		s += fmt.Sprintf("\nYou Take Time %.0f Seconds\n\n", m.timeRecord)
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

// reOrder
// Reordering the spells that will be displayed for the user to cast or invoke.
// This reordering will be triggered if the user invokes or casts the correct answer.
// The time complexity for this operation will be O(N).
func reOrder(spells []string, answer string) []string {
	newSpells := make([]string, 0, len(spells)-1)

	for _, spell := range spells {
		if spell != answer {
			newSpells = append(newSpells, spell)
		}
	}

	return newSpells
}

// TODO performance?
// create random array of spells to incantate
// it can either be 2 or 3 spells
// if previous incantate is nil, create a new one with no overlap
// else overlap is fine
func gen(prev incantate) (bool, incantate, []string) {
	var cast bool
	r := rand.New(rand.NewSource(time.Now().Unix()))
	if c := r.Intn(2); c != 0 {
		cast = true
	}
	next := make(map[spell]struct{})
	length := 2

	spells := make([]string, 0, length)

	// TODO check how many passes do we need to generate
	for len(next) < length {
		n := spell(1 + r.Intn(10)) // this will rarely gets to deafening?
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
