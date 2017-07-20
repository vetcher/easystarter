package printer

import (
	"fmt"
	"log"
	"os"
)

type Printer struct {
	logger *log.Logger
}

var defaultPrinter = Printer{
	logger: log.New(os.Stdout, "", 0),
}

func (p *Printer) Print(marker string, args ...interface{}) {
	p.logger.Printf("[%v] %v\n", marker, fmt.Sprint(args...))
}

func (p *Printer) Printf(marker string, format string, args ...interface{}) {
	p.logger.Printf("[%v] %v\n", marker, fmt.Sprintf(format, args...))
}

func (p *Printer) PrintRaw(args ...interface{}) {
	p.logger.Print(args...)
}

func Print(marker string, args ...interface{}) {
	defaultPrinter.Print(marker, args...)
}

func Printf(marker string, format string, args ...interface{}) {
	defaultPrinter.Printf(marker, format, args...)
}

func PrintRaw(args ...interface{}) {
	defaultPrinter.PrintRaw(args...)
}
