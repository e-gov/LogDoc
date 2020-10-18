# LogDoc

LogDoc on koodianalüüsi ja -dokumenteerimisvahend, mis koostab koodibaasis tehtavatest logimistest ülevaate ja abistab inimest logilausete kommenteerimisel.

Logilause on koodibaasi lause, mis koostab ja kirjutab logikirje. Logilaused on äratuntavad eelkõige logimisteegi kasutamise järgi. Logilaused on laiali üle kogu koodibaasi. Logimise katvuse hindamiseks, aga samuti logi mõistmiseks on vaja ülevaadet, mida logitakse ja arusaamist, mida logikirjed tähendavad. LogDoc aitab neid vajadusi rahuldada, sellega, et koostavb täieliku, kogu koodibaasi hõlmava nimekirja logimistest ja aitab inimesel siduda logimislausetega inimloetavaid kommentaare. 

LogDoc on mõeldud juhtudeks, kus arendaja oma koodi ei dokumenteeri s.t koodibaas ei sisalda kommentaare.

LogDoc koosneb rakendusest ja logilausete failist.

Logilausete fail on järgmise struktuuriga. Märkus: Süntaks on kirjeldatud EBNF abil, vt https://golang.org/ref/spec#Notation. 
 
````
logilausete_fail = { logilause_kirjeldus } .
logilause_kirjeldus =
    "----"
    logilause
    "----"
    [ viidad ]
    [ kommentaarid ] .
viidad = { viit } .
viit = "(" failitee "," reanumber "," funktsooninimi ")" .
kommentaarid = { kommentaaririda } .
````

`logilause`, `viit` ja `kommentaar` on eraldi real. `----` on eraldaja. Tühje ridu ei arvestata. 

`logilause` on koodibaasist kopeeritud logi kirjutav lause (tehniliselt: Go AST "pretty-print" kujul). Logilause on ühel või mitmel real.

`viit` näitab koodilause asukohta koodibaasis. Samakujuline logilause võib koodibaasis esineda mitmes kohas. Logilause kirjeldusse kogutakse viited kõigile esinemistele. Iga viit on eraldi real.

`kommentaarid` koosneb ühest või enamast tekstireast.

Logilausete faili genereerib LogDoc rakendus. Fail on lihtsa struktuuriga tekstifail. Inimene saab faili lugeda ja sinna kommentaare lisada. 

LogDoc rakendus korjab koodibaasist kokku logilaused ja salvestab logilausete faili. Logilausele lisatakse viidad kohtadele, kus lause koodis esineb. Rakendust käivitatakse perioodiliselt, hõivamaks arenduses toimunud muutusi. Kui logilause auskoht koodis on muutunud, siis LogDoc uuendab viita(sid).

Inimene saab logilausete kirjeldustele failis lisada kommentaare tavalise tekstiredaktoriga.

LogDoc analüüsib Go koodi. Testifaile ei analüüsita.

Kasutamine:

LogDoc käivitatakse käsurealt. LogDoc pakub 3 alamkäsku:

- `create` - kogub koodibaasist logilaused ja moodustab LogDoc-faili
- `update` - uuendab LogDoc-faili koodibaasi muutustega
- `clear` - eemaldab LogDoc-failist logilaused, mis on kaotanud seose koodibaasiga  
- `stat` - väljastab statistikat logilausete kohta.

Koodibaasi ja LogDoc-faili asukoht antakse lippudega:

````
logdoc create -dir <koodibaasi kaust> -logdocfile <logilausete fail>
````

````
logdoc update -dir <koodibaasi kaust> -logdocfile <logilausete fail>
````

````
logdoc clear -logdocfile <logilausete fail>
````

````
logdoc stat -logdocfile <logilausete fail>
````

`-dir` on kaust, millest ja mille alamkaustadest logilauseid kogutakse.

`-logdocfile` on logilausete faili nimi (failitee).
