package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// Tõstetud globaalseteks, kuna astInspector vormistamisega iseseisvateks f-deks ei ole
// nende edastamine lihtne/võimalik.
var (
	// treeStack on AST puu läbimise pinu.
	treeStack []ast.Node
	// FileSet
	fset *token.FileSet
)

// walkFn läbib kausta v faili, otsides logilauseid.
func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	// Töötle kaust.
	if info.IsDir() {
		// Kontrolli, kas kaustas in moodul. Prindi mooduli nimi.
		modFile := filepath.Join(path, "go.mod")
		dat, err := ioutil.ReadFile(modFile)
		if err == nil {
			re := regexp.MustCompile("[^\n]*")
			moduleLine = re.FindString(string(dat))
			fmt.Printf("** %s **\n\n", moduleLine)
			noModules++
		}
		// Töötle kausta failid, kogumina, parser.ParseDir abil.
		fset = token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, path, isGoFile, 0)
		if err != nil {
			fmt.Printf("Viga kausta lugemisel: %s\n", err)
			os.Exit(1)
		}
		for name, pkg := range pkgs {
			fmt.Printf("    ** package %s **\n\n", name)
			// Läbi paki AST-puu, otsides logilauseid.
			ast.Inspect(pkg, astInspector)
			noPackages++
		}
		// Kausta alamkaustade töötlemise tagab filepath.Walk.
	}
	return nil
}

// isGoFile sõelub välja (jätab välja) Go testifailid.
func isGoFile(f os.FileInfo) bool {
	n := f.Name()
	return len(n) < 9 || n[len(n)-8:] != "_test.go"
}

// walkFile otsib failist fname logilaused.
/* func walkFile(fname string) {
	// Parsi fail fname AST puuks.
	fset = token.NewFileSet()
	// Tagastab ast.File tipu. https://golang.org/pkg/go/parser/#ParseFile
	node, err := parser.ParseFile(fset, fname, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Viga faili parsimisel: %v", err)
		os.Exit(1)
	}
	// Jaluta kood läbi
	ast.Inspect(node, astInspector)
} */

// astInspector uurib ja töötleb AST-puu tippu.
func astInspector(n ast.Node) bool {
	var buf bytes.Buffer // Logilause AST "pretty-print" kuju koostamiseks.
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
	// Kas tipp on pakk?
	/* pkg, isPkg := n.(*ast.Package)
	if isPkg {
		pkgName = pkg.Name
		fmt.Printf("Pakk: %s\n", pkgName)
		return true
	} */
	// Kas tipp on identifikaator?
	id, isID := n.(*ast.Ident)
	if isID {
		// fmt.Printf("id-tipp: nimi: %s\n", id.Name)
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
				// Valmista väljastatav logilause
				buf.Reset()
				printer.Fprint(&buf, fset, treeStack[esp])
				k := logStmnt(buf.String())
				// Valmista logilauseasukoht.
				loc := fmt.Sprintf("(%v, %v)", funcName, fset.Position(treeStack[esp].Pos()).Line)
				// Prindi logilause ja selle asukoht.
				fmt.Printf("        %s\n", loc)
				fmt.Printf("            %s\n\n", k)

				noLogStmts = noLogStmts + 1
			}
		}
	}
	return true
}

// Märkmed
// go/printer Fprint - AST-tipu "pretty print".
// https://pkg.go.dev/go/printer#pkg-functions
