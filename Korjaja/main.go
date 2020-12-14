// `Korjaja` on rakendus, korjab koodibaasist kokku logilaused ja väljastab väljundvoogu.
// Vt lähemalt ülakausta README.md-failist.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	rootFolder string // Analüüsitav koodibaas (kaust).
	noModules  int    // Töödeldud mooduleid.
	noPackages int    // Töödeldud pakke.
	noLogStmts int    // Leitud logilauseid.

	moduleLine     string // Parajasti läbitava mooduli nimi (1. rida go.mod failist).
	pkgName        string // Parajasti läbitava paki nimi.
	pkgNamePrinted string // Viimati prinditud paki nimi. Kasutusel mitmekordse printimise vältimiseks.
)

type mpf string      // Logilause asukohta näitav sõne, kujul <moodul>/<pakk> <f-n>.
type logStmnt string // LogStmnt on logilause, AST "pretty-print" kujul.

func main() {
	// Loe ja töötle ohjelipud.
	codeBasePath := flag.String("dir", "", "Koodikaust")
	flag.Parse()
	if *codeBasePath == "" {
		fmt.Println("Anna koodikaust: -dir <kausta nimi>")
		os.Exit(1)
	}
	rootFolder = *codeBasePath
	fmt.Printf("*** LogDoc ****\n\nLogilausete korje koodibaasist\n\n")
	t := time.Now()
	fmt.Printf(t.Format("Korje tehtud:\t\t\t02.01.2006 15:04\n\n"))
	fmt.Printf("Korjan logilauseid kaustast:\t%s\n", rootFolder)

	// Läbi koodibaas, otsides logilauseid.
	err := filepath.Walk(rootFolder, walkFn)
	if err != nil {
		fmt.Printf("Viga kausta läbijalutamisel: %v\n", err)
		os.Exit(1)
	}

	// Väljasta statistika
	fmt.Printf("\n\nStatistika\n\n")
	fmt.Printf("Mooduleid (Go module):\t\t%v\n", noModules)
	fmt.Printf("Pakke (Go package):\t\t%v\n", noPackages)
	fmt.Printf("Logilauseid (v.a testides):\t%v\n", noLogStmts)
}

// Märkmed
// https://yourbasic.org/golang/list-files-in-directory/
// Kausta läbijalutamine
// "github.com/davecgh/go-spew/spew"
// AST ilus, sügav väljastus
// https://golang.org/pkg/go/ast/#Inspect
// ast.Inspect, koos näitega
// https://yourbasic.org/golang/implement-stack/
// Pinu teostus
// AST printer
// https://pkg.go.dev/go/printer#example-Fprint
// https://gobyexample.com/command-line-subcommands
// Go alamkäsud
