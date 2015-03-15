package cludo

import (
	"fmt"
	"io"
)

type UnitAttrPair struct {
	// Key-Value pair in a unit section
	key, value string
}

type Unit struct {
	keyOrder []string
	// Unit file sections.  Map header -> slice of UnitAttrPair
	sections map[string][]UnitAttrPair
}

// Construct Unit, and initialize empty sections map
func MakeUnit() *Unit {
	unit := new(Unit)
	unit.sections = make(map[string][]UnitAttrPair)
	return unit
}

func (unit *Unit) AddSection(header string) {
	// Preserver order of sections
	unit.keyOrder = append(unit.keyOrder, header)
	// Initialize empty slice
	unit.sections[header] = []UnitAttrPair{}
}

func (unit *Unit) AddItem(header string, key string, value string) {
	unit.sections[header] = append(unit.sections[header], UnitAttrPair{key, value})
}

// Export Unit object to native coreos/systemd unit file format
func (unit *Unit) Export(buf io.Writer) {
	for _, header := range unit.keyOrder {
		buf.Write([]byte(fmt.Sprintf("[%s]\n", header)))
		for attr := range unit.sections[header] {
			buf.Write([]byte(fmt.Sprintf("%s=%s\n", unit.sections[header][attr].key, unit.sections[header][attr].value)))
		}
		buf.Write([]byte(fmt.Sprintln()))
	}
}
