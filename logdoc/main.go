package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	rootFolder     string // Analüüsitava koodibaasi kaust.
	logDocFileName string // Logilausekirjelduste fail.
	nof            int    // Töödeldud faile.
	noLogStmts     int    // Leitud logilauseid.
)

type logStmnt string // LogStmnt on logilause, AST "pretty-print" kujul.

type logStmntDesc struct { // LogStmntDesc on logilause kirjeldus.
	locations []string
	comments  []string
}

// logStmntDescs on logilausete kirjeldused, mäpina.
var logStmntDescs map[logStmnt]logStmntDesc

func main() {

	readFlags()

	// Loe logilausekirjeldused failist mäppi.
	readLogStmntDescs(logDocFileName)

	// Väljasta kontrolliks sisseloetud logilausekirjeldused.
	writeLogStmntDescs("./Abi.txt")

	// Läbi koodibaas, otsides logilauseid.
	walk()

	// Salvesta logilausete kirjelduste mäpp faili.
	writeLogStmntDescs(logDocFileName)

	// Väljasta statistika
	fmt.Printf("Go-faile (v.a testid): %v\n", nof)
	fmt.Printf("Logilauseid: %v\n", noLogStmts)

}

// readFlags loeb ja kontrollib käsureaohjeväärtused.
func readFlags() {
	rootFolderPtr := flag.String("dir", "", "Koodikaust")
	logDocFileNamePtr := flag.String("logdocfile", "", "Logilausekirjelduste fail")

	flag.Parse()
	if *rootFolderPtr == "" {
		fmt.Println("Anna koodikaust: -dir <kausta nimi>")
		os.Exit(1)
	}
	rootFolder = *rootFolderPtr
	if *logDocFileNamePtr == "" {
		fmt.Println("Anna LogDoc faili asukoht: -logdocfile <failitee>")
		os.Exit(1)
	}
	logDocFileName = *logDocFileNamePtr

	fmt.Printf("LogDoc\nKaust: %s\nLogDoc fail: %s\n",
		rootFolder, logDocFileName)
}

// walk läbib koodibaasi, otsides logilauseid.
func walk() {
	err := filepath.Walk(
		rootFolder,
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
}

// writeLogStmntDescs kirjutab  logilausekirjeldused mäpist faili.
func writeLogStmntDescs(fname string) {
	txtFile, err := os.Create(fname)
	if err != nil {
		fmt.Println("Viga faili avamisel")
	}
	for key, element := range logStmntDescs {
		txtFile.Write([]byte(fmt.Sprintf("%v\n", key)))
		for _, l := range element.locations {
			txtFile.Write([]byte(fmt.Sprintf("%v\n", l)))
		}
		for _, c := range element.comments {
			txtFile.Write([]byte(fmt.Sprintf("%v\n", c)))
		}
		txtFile.Write([]byte("\n"))
	}
	txtFile.Close()
}

// readLogStmntDescs loeb logilausekirjeldused failist mäppi.
func readLogStmntDescs(fname string) {

	const (
		BOF   = iota // Faili algus.
		STMNT        // Logilause.
		LOCS         // Logilause asukohad.
		CMNT         // Selgitus.
		EMPTY        // Tühi rida.
	)

	var cLS logStmnt // Failist parajasti sisseloetav logilause

	// Lähtesta mäpp.
	logStmntDescs = make(map[logStmnt]logStmntDesc)

	dat, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println("Viga faili lugemisel. Faili ei ole?")
		return
	}
	fcontents := string(dat)
	lines := strings.Split(fcontents, "\n")
	pLT := BOF // Eelmise rea tüüp.
	for _, line := range lines {
		switch pLT {
		case BOF: // Faili algus
			cLS = logStmnt(line)
			pLT = STMNT
			break
		case STMNT:
			// Asukohti sisse ei ole; need kirjutame alati üle.
			pLT = LOCS
			break
		case LOCS:
			if line == "" {
				pLT = EMPTY
				break
			}
			if line[0] == '(' {
				// Asukohti sisse ei ole; need kirjutame alati üle.
				break
			}
			abi := logStmntDescs[cLS]
			abi.comments = append(abi.comments, line)
			logStmntDescs[cLS] = abi
			pLT = CMNT
			break
		case CMNT:
			if line == "" {
				pLT = EMPTY
			} else {
				abi := logStmntDescs[cLS]
				abi.comments = append(abi.comments, line)
				logStmntDescs[cLS] = abi
				pLT = CMNT
			}
			break
		case EMPTY:
			cLS = logStmnt(line)
			pLT = STMNT
			break
		}
	}
}

// walkFile otsib failist fname logilaused.
func walkFile(fname string) {

	var buf bytes.Buffer // Logilause AST "pretty-print" kuju koostamiseks.

	// Parsi fail fname AST puuks.
	fset := token.NewFileSet()
	// Tagastab ast.File tipu.
	// https://golang.org/pkg/go/parser/#ParseFile
	node, err := parser.ParseFile(fset, fname, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Viga faili parsimisel: %v", err)
		os.Exit(1)
	}

	// treeStack on AST puu läbimise pinu.
	var treeStack []ast.Node

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
			if id.Name == "Info" || id.Name == "Error" || id.Name == "Debug" {
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
				var funcName string
				if fnp != -1 {
					fN := treeStack[fnp].(*ast.FuncDecl)
					funcName = fN.Name.Name
				} else {
					funcName = ""
				}

				if esp != -1 {
					// Kanna mäppi.
					// Valmista väljastatav logilause
					buf.Reset()
					printer.Fprint(&buf, fset, treeStack[esp])
					k := logStmnt(buf.String())
					// Valmista logilauseasukoht.
					l := fmt.Sprintf("(%v, %v, %v)", fname, funcName, fset.Position(treeStack[esp].Pos()).Line)

					// Kas on juba mäpis?
					_, ok := logStmntDescs[k]
					if ok {
						// Kirjuta asukohad üle. Õigemini, asukohti sisse ei loegi.
					} else {
						logStmntDescs[k] = logStmntDesc{
							locations: []string{l},
							comments:  []string{},
						}
					}
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

// AST printer
// https://pkg.go.dev/go/printer#example-Fprint
