# LogDoc

LogDoc on koodianalüüsi- ja dokumenteerimisvahend, mis koostab koodibaasis tehtavatest logimistest ülevaate ja abistab inimest logimiskohtadele inimloetavate selgituste lisamisel. 

LogDoc koosneb rakendusest ja failist.

Rakendus käib läbi kogu koodibaasi ja tuvastab kohad, kus
logitakse (logimislaused). Tuvastamine põhineb asjaolul, et teada on
koodibaasis kasutatav logimisteek.
Tuvastatud logimislausete kohta peab rakendus arvestust eraldi failis
(edaspidi - LogDoc fail).
Inimene lisab failis logimislausetele inimloetavad seletused.

Koodi muutudes lastakse LogDoc uuesti koodibaasi analüüsima.
Iga koodis tuvastatud logimislause kohta rakendus:
- kontrollib, kas logimislause juba on LogDoc failis
- kui ei ole, siis lisab ja markeerib inimesele selgituse lisamiseks
- kui ei ole, aga on sarnane lause, siis rakendus markeerib sarnasuse
ja lisab logimislause sarnase lause kõrvale.

Inimene avab faili ja kasutades taustateadmist ning vajadusel
uurides täiendavalt koodi, lisab igale logimislausele inimloetava
kommentaari, samuti lahendab sarnasuskonfliktid.

LogDoc kogub koodibaasist kokku logilaused ( `.Info()`, `.Debug()` või `.Error()` sisaldavad avaldislaused) ja esitab need inimloetava koondina (väljundina konsoolile). Testifaile ei analüüsita.

Kasutamine:

````
go run . -dir <kausta nimi> [-level <logitase>]
````

`-dir` on kaust, millest ja mille alamkaustadest logilauseid kogutakse.

`-logitase` väärtuseks anda `Info`, `Debug` või `Error`. Vaikimisi haaratakse
väljundisse kõik logitasemed. Väärtused on tõstutundlikud.