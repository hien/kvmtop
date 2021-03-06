package printers

import (
	"fmt"
	"strconv"

	"github.com/cha87de/kvmtop/models"
)

// JSONPrinter describes the json printer
type JSONPrinter struct {
	models.Printer
}

// Open opens the printer
func (printer *JSONPrinter) Open() {
	outputOpen()
}

// Screen prints the measurements on the screen
func (printer *JSONPrinter) Screen(printable models.Printable) {
	output(fmt.Sprintf("{ \"host\": {"))
	hostFields := printable.HostFields
	hostValues := printable.HostValues
	i := 0
	for _, value := range hostValues {
		if i > 0 {
			output(fmt.Sprintf(","))
		}

		// but """ only for strings
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			output(fmt.Sprintf("\"%s\": %d", hostFields[i], intValue))
		} else if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			output(fmt.Sprintf("\"%s\": %f", hostFields[i], floatValue))
		} else {
			output(fmt.Sprintf("\"%s\": \"%s\"", hostFields[i], value))
		}
		i++
	}

	output(fmt.Sprintf("}, \"domains\": ["))

	domainFields := printable.DomainFields
	domainValues := printable.DomainValues
	i = 0
	for domvalue := range domainValues {
		if i > 0 {
			output(fmt.Sprintf(","))
		}
		output(fmt.Sprintf("{"))
		for j, value := range domainValues[domvalue] {
			if j > 0 {
				output(fmt.Sprintf(","))
			}

			// but """ only for strings
			if _, err := strconv.ParseInt(value, 10, 64); err == nil {
				output(fmt.Sprintf("\"%s\": %s", domainFields[j], value))
			} else if _, err := strconv.ParseFloat(value, 64); err == nil {
				output(fmt.Sprintf("\"%s\": %s", domainFields[j], value))
			} else {
				output(fmt.Sprintf("\"%s\": \"%s\"", domainFields[j], value))
			}
		}
		output(fmt.Sprintf("}"))
		i++
	}
	output(fmt.Sprintf("]}\n"))
}

// Close terminates the printer
func (printer *JSONPrinter) Close() {
	outputClose()
}

// CreateJSON creates a new simple text printer
func CreateJSON() JSONPrinter {
	return JSONPrinter{}
}
