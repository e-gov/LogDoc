/*
LogDoc kogub  koodibaasist kokku logilaused ( .Info(), .Debug() või
.Error() sisaldavad avaldislaused) ja esitab need inimloetava koondina
(väljundina konsoolile). Testifaile ei analüüsita.

Kasutamine:

go run . -dir <kausta nimi> [-level <logitase>]

-dir on kaust, millest ja mille alamkaustadest logilauseid kogutakse.

-logitase väärtuseks anda Info, Debug või Error. Vaikimisi haaratakse
väljundisse kõik logitasemed.

*/
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"time"
)

var (
	// "Info", "Debug", "Error"
	logLevel string
	// Leitud logilauseid
	noLogStmts int
)

func main() {

	rootFolder := flag.String("dir", "", "Koodikaust")
	logLevelPtr := flag.String("level", "", "Logitase (Info, Debug, Error")

	flag.Parse()
	if *rootFolder == "" {
		fmt.Println("Anna koodikaust: -dir <kausta nimi>")
		os.Exit(1)
	}
	logLevel = *logLevelPtr

	fmt.Printf("\nLOGILAUSETE KOOND\n\nKaust: %s\nLogitase: %s\n",
		*rootFolder, *logLevelPtr)
	fmt.Printf("%v\n\n", time.Now().Format("2006.01.02 15:04:05"))

	// Statistika kogujad.
	// Töödeldud faile
	nof := 0

	err := filepath.Walk(
		*rootFolder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Töötle fail
			// fmt.Println(path, info.Size())
			// Ei ole kaust
			if !info.IsDir() &&
				// on go fail
				filepath.Ext(path) == ".go" &&
				// ei ole testifail
				path[len(path)-8:] != "_test.go" {
				walkFile(path)
				nof = nof + 1
			}
			return nil
		},
	)

	if err != nil {
		fmt.Printf("Viga kausta läbijalutamisel: ")
		os.Exit(1)
	}

	// Väljasta statistika
	fmt.Printf("\nGo-faile (v.a testid): %v\n", nof)
	fmt.Printf("\nLogilauseid: %v\n", noLogStmts)

}

// walkFile otsib failist fname logilaused.
func walkFile(fname string) {

	// Parsi fail fname AST puuks.
	fset := token.NewFileSet()
	// Tagastab ast.File tipu.
	// https://golang.org/pkg/go/parser/#ParseFile
	node, err := parser.ParseFile(fset, fname, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Viga faili parsimisel: %v", err)
		os.Exit(1)
	}

	var (
		// Väljasta faili nimi alles esimese huvitava tipu leidmisel;
		// sellega tagad, et ei väljasta ebahuvitavate failide nimesid.
		// fnamePrinted bool
		// treeStack on AST puu läbimise pinu.
		treeStack []ast.Node
	)

	// Jaluta kood läbi
	ast.Inspect(node, func(n ast.Node) bool {

		// Inspekteerimisf-n kutsutakse argumendiga nil välja
		// pärast kõigi alamtippude läbimist.
		if n == nil {
			// Eemalda tipp pinust.
			if len(treeStack) > 0 {
				n := len(treeStack) - 1 // Top element
				treeStack = treeStack[:n]
			}
			return true
		}

		// Lisa tipp pinussse
		treeStack = append(treeStack, n)

		// Töötlus vastavalt tipu tüübile
		// Kas tipp on identifikaator?
		id, isID := n.(*ast.Ident)
		if isID {
			// Kas id on "logLevel"?
			if (logLevel == "" && (id.Name == "Info" || id.Name == "Error" || id.Name == "Debug")) ||
				(logLevel == "" && id.Name == logLevel) {
				// Väljasta faili nimi (kui seda ei ole veel tehtud)
				// if !fnamePrinted {
				// 	fmt.Printf("\n%s\n", fname)
				// 	fnamePrinted = true
				// }

				// Otsi pinust viimane avaldislause
				// si (stackIndex) on pinu järgmisena vaadeldava elemendi
				// indeks.
				si := len(treeStack) - 2
				// esp (expressionStatementPos) on pinu läbimisel leitud
				// avaldislause tüüpi tipu indeks; kui pole leitud, siis -1.
				esp := -1
				for si >= 0 && esp == -1 {
					// Kas pinu elemendis si on avaldislause?
					_, isES := treeStack[si].(*ast.ExprStmt)
					if isES {
						// Avaldislause leitud; ära enam otsi
						esp = si
					}
					si--
				}
				// Otsi ka f-ni nimi
				// fnp (functionPos) on pinu läbimisel leitud funktsioonidekl-i tüüpi
				// tipu indeks; kui pole leitud, siis -1.
				fnp := -1
				for si >= 0 && fnp == -1 {
					// kas pinu elemendis on funktsioonideklaratsioon?
					_, isFD := treeStack[si].(*ast.FuncDecl)
					if isFD {
						// F-dekl-n leitud, ära enam otsi
						fnp = si
					}
					si--
				}

				// Valmista ette f-ni nimi (koos sulgudega) printimiseks.
				var fn string
				if fnp != -1 {
					fN := treeStack[fnp].(*ast.FuncDecl)
					fn = "(" + fN.Name.Name + ")"
				} else {
					fn = ""
				}

				// Väljasta logilause kujul:
				// failinimi reanr (f-ni nimi)
				// logilause
				if esp != -1 {
					fmt.Printf("\n%s ", fname)
					fmt.Printf(
						"%d: %v\n",
						fset.Position(treeStack[esp].Pos()).Line,
						fn,
					)
					// Fprint "pretty-prints" an AST node to output.
					printer.Fprint(os.Stdout, fset, treeStack[esp])
					fmt.Printf("\n")

					noLogStmts = noLogStmts + 1
				}

			}
		}

		return true
	})

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
