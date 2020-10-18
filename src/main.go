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
	switch os.Args[1] {
	case "create":
		// Läbi koodibaas, otsides logilauseid.
		// Lähtesta mäpp.
		logStmntDescs = make(map[logStmnt]logStmntDesc)
		walk()
		// Salvesta logilausete kirjelduste mäpp faili.
		writeLogStmntDescs(logDocFileName)
		// Väljasta statistika
		fmt.Printf("Go-faile (v.a testid): %v\n", nof)
		fmt.Printf("Logilauseid: %v\n", noLogStmts)
	case "update":
		// Loe logilausekirjeldused (ilma viideteta koodibaasile) failist mäppi.
		readLogStmntDescs(logDocFileName, false)
		// Salvesta logilausete kirjelduste mäpp kontrolliks faili.
		// writeLogStmntDescs("./Abi.txt")
		// Läbi koodibaas, otsides logilauseid.
		walk()
		// Salvesta logilausete kirjelduste mäpp faili.
		writeLogStmntDescs(logDocFileName)
		// Väljasta statistika
		fmt.Printf("Töödeldud Go-faile (v.a testid): %v\n", nof)
		fmt.Printf("Koodibaasis leitud logilauseid: %v\n", noLogStmts)
		fmt.Printf("LogDoc-failis logilauseid: %v\n", len(logStmntDescs))
	case "clear":
		// Loe logilausekirjeldused, sh viited koodibaasile, failist mäppi.
		readLogStmntDescs(logDocFileName, true)
		clearLogDocFile()
	case "stat":
		// Loe logilausekirjeldused, sh viited koodibaasile, failist mäppi.
		readLogStmntDescs(logDocFileName, true)
		// Salvesta logilausete kirjelduste mäpp kontrolliks faili.
		// writeLogStmntDescs("./Abi.txt")
		statistics()
	}
}

// readFlags loeb ja kontrollib käsurea ohjeväärtused.
func readFlags() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createCodebase := createCmd.String("dir", "", "Koodikaust")
	createLogDocFile := createCmd.String("logdocfile", "", "Logilausekirjelduste fail")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateCodebase := updateCmd.String("dir", "", "Koodikaust")
	updateLogDocFile := updateCmd.String("logdocfile", "", "Logilausekirjelduste fail")

	clearCmd := flag.NewFlagSet("clear", flag.ExitOnError)
	clearLogDocFile := clearCmd.String("logdocfile", "", "Logilausekirjelduste fail")

	statCmd := flag.NewFlagSet("stat", flag.ExitOnError)
	statLogDocFile := statCmd.String("logdocfile", "", "Logilausekirjelduste fail")

	if len(os.Args) < 2 {
		fmt.Println("Anna alamkäsk 'create', 'update', 'stat' või 'clear'")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if *createCodebase == "" {
			fmt.Println("Anna koodikaust: -dir <kausta nimi>")
			os.Exit(1)
		}
		if *createLogDocFile == "" {
			fmt.Println("Anna LogDoc faili asukoht: -logdocfile <failitee>")
			os.Exit(1)
		}
		rootFolder = *createCodebase
		logDocFileName = *createLogDocFile
		fmt.Printf("Kaust: %s\nLogDoc fail: %s\n", rootFolder, logDocFileName)
	case "update":
		updateCmd.Parse(os.Args[2:])
		if *updateCodebase == "" {
			fmt.Println("Anna koodikaust: -dir <kausta nimi>")
			os.Exit(1)
		}
		if *updateLogDocFile == "" {
			fmt.Println("Anna LogDoc faili asukoht: -logdocfile <failitee>")
			os.Exit(1)
		}
		rootFolder = *updateCodebase
		logDocFileName = *updateLogDocFile
		fmt.Printf("Kaust: %s\nLogDoc fail: %s\n", rootFolder, logDocFileName)
	case "clear":
		clearCmd.Parse(os.Args[2:])
		if *clearLogDocFile == "" {
			fmt.Println("Anna LogDoc faili asukoht: -logdocfile <failitee>")
			os.Exit(1)
		}
		logDocFileName = *clearLogDocFile
		fmt.Printf("LogDoc fail: %s\n", logDocFileName)
	case "stat":
		statCmd.Parse(os.Args[2:])
		if *statLogDocFile == "" {
			fmt.Println("Anna LogDoc faili asukoht: -logdocfile <failitee>")
			os.Exit(1)
		}
		logDocFileName = *statLogDocFile
		fmt.Printf("LogDoc fail: %s\n", logDocFileName)
	}
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
			if !info.IsDir() && // Ei ole kaust
				filepath.Ext(path) == ".go" && // on go fail
				path[len(path)-8:] != "_test.go" { // ei ole testifail
				fmt.Printf("%s\n", path)
				walkFile(path)
				nof = nof + 1
			}
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Viga kausta läbijalutamisel: %v\n", err)
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
		txtFile.Write([]byte(fmt.Sprintf("----\n")))    // Väljasta eraldaja.
		txtFile.Write([]byte(fmt.Sprintf("%v\n", key))) // Väljasta logilause.
		txtFile.Write([]byte(fmt.Sprintf("----\n")))    // Väljasta eraldaja.
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
func readLogStmntDescs(fname string, readLocs bool) {
	const (
		BOF  = iota // Faili algus.
		LOGB        // Logilause aluseraldaja ('----')
		LOGS        // Logilause rida.
		LOGE        // Logilause lõpueraldaja ('----')
		LOCS        // Logilause asukohad.
		CMNT        // Selgitus.
	)
	var LSBuf string // Abimuutuja logilause kogumiseks.
	var cLS logStmnt // Failist parajasti sisseloetav logilause.

	// Lähtesta mäpp.
	logStmntDescs = make(map[logStmnt]logStmntDesc)

	dat, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println("Viga LogDoc-faili lugemisel. Faili ei ole?")
		return
	}
	fcontents := string(dat)
	lines := strings.Split(fcontents, "\n")
	prev := BOF // Eelmise rea tüüp.
	for lno, line := range lines {
		// fmt.Println("Read: " + line)
		switch prev {
		case BOF: // Faili algus
			if line == "----" {
				prev = LOGB
			}
		case LOGB: // Logilause alguseraldaja
			if line == "----" {
				fmt.Printf("Ebakorrektne LogDoc-file, rida: %v", lno+1)
				os.Exit(1)
			}
			LSBuf = line
			prev = LOGS
		case LOGS: // Logilause rida
			if line == "----" {
				cLS = logStmnt(LSBuf)
				logStmntDescs[cLS] = logStmntDesc{[]string{}, []string{}}
				// fmt.Printf("Lugesin logilause rea: %v\n", cLS)
				prev = LOGE
			} else {
				if LSBuf == "" {
					LSBuf = line
				} else {
					LSBuf = LSBuf + "\n" + line
				}
			}
		case LOGE: // Logilause lõpueraldaja
			if line == "" { // Tühiridu ei loe kommentaaridena.
				break
			}
			if line == "----" { // Puuduvad nii viited kui kommentaarid.
				prev = LOGB
				break
			}
			if line[0] == '(' { // Viit.
				// update puhul asukohti sisse ei loe; need kirjutame alati üle.
				if readLocs {
					h := logStmntDescs[cLS]
					h.locations = append(h.locations, line)
					logStmntDescs[cLS] = h
					// fmt.Println("Lisatud: viit")
				}
				prev = LOCS
				break
			}
			// Jääb üle ainult kommentaar.
			h := logStmntDescs[cLS]
			h.comments = append(h.comments, line)
			logStmntDescs[cLS] = h
			// fmt.Println("Lisatud: kommentaar")
			prev = CMNT
		case LOCS:
			if line == "" {
				break
			}
			if line[0] == '(' {
				if readLocs {
					h := logStmntDescs[cLS]
					h.locations = append(h.locations, line)
					logStmntDescs[cLS] = h
					// fmt.Println("Lisatud: viit")
				}
				break
			}
			if line == "----" {
				prev = LOGB
				break
			}
			h := logStmntDescs[cLS]
			h.comments = append(h.comments, line)
			logStmntDescs[cLS] = h
			// fmt.Println("Lisatud: kommentaar")
			prev = CMNT
		case CMNT:
			if line == "----" {
				prev = LOGB
				break
			}
			if line != "" {
				h := logStmntDescs[cLS]
				h.comments = append(h.comments, line)
				logStmntDescs[cLS] = h
				// fmt.Println("Lisatud: kommentaar")
			}
		}
	}
}

// walkFile otsib failist fname logilaused.
func walkFile(fname string) {
	var buf bytes.Buffer // Logilause AST "pretty-print" kuju koostamiseks.
	// Parsi fail fname AST puuks.
	fset := token.NewFileSet()
	// Tagastab ast.File tipu. https://golang.org/pkg/go/parser/#ParseFile
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
					loc := fmt.Sprintf("(%v, %v, %v)", fname, funcName, fset.Position(treeStack[esp].Pos()).Line)

					// Kas on juba mäpis?
					_, ok := logStmntDescs[k]
					if ok {
						// Lisa asukoht.
						h := logStmntDescs[k]
						h.locations = append(h.locations, loc)
						logStmntDescs[k] = h
					} else {
						logStmntDescs[k] = logStmntDesc{
							locations: []string{loc},
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

// clearLogDocFile eemaldab LogDoc failist logilaused, mille seosed koodibaasiga on kadunud.
func clearLogDocFile() {
	nd := 0
	for key, elem := range logStmntDescs {
		if len(elem.locations) == 0 {
			delete(logStmntDescs, key)
			fmt.Printf("Eemaldasin logilause %v.\n", key)
			nd++
		}
	}
	writeLogStmntDescs(logDocFileName)
	fmt.Printf("Eemaldasin %v logilauset.\n", nd)
}

// statistics väljastab statistikat LogDoc-faili kohta.
func statistics() {
	fmt.Printf("Logilauseid: %v\n", len(logStmntDescs))
	nd := 0
	for _, elem := range logStmntDescs {
		if len(elem.locations) == 0 {
			nd++
		}
	}
	fmt.Printf("Koodibaasiga seose kaotanud logilauseid: %v\n", nd)
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
