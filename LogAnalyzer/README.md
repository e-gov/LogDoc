### LogAnalyzer

LogAnalyzer kogub koodibaasist kokku logilaused ( `.Info()`, `.Debug()` või `.Error()` sisaldavad avaldislaused) ja esitab need inimloetava koondina (väljundina konsoolile). Testifaile ei analüüsita.

Kasutamine:

````
go run . [-dir <kausta nimi>] [-level <logitase>]
````

`-dir` on kaust, millest ja mille alamkaustadest logilauseid kogutakse.

`-logitase` väärtuseks anda `Info`, `Debug` või `Error`. Vaikeväärtus on `Info`. NB: väärtused on tõstutundlikud.

Väljundi näited on failides `LogInfo.txt`, `LogDebug.txt` ja `LogError.txt`.
