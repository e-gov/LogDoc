# LogDoc

LogDoc on koodianalüüsi ja logimise kavandamise ja dokumenteerimise vahendite kogum.

`Korjaja` on rakendus, mis korjab koodibaasist kokku logilaused ja kirjutab väljundvoogu.

Edaspidi lisandub teisi vahendeid.

Logilause on koodibaasi lause, mis koostab ja kirjutab logikirje. Logilaused on äratuntavad logimisteegi kasutamise järgi (meetodiväljakutsed `Info()`, `Error()`, `Debug()`).

Logilaused on laiali üle kogu koodibaasi. Logimise katvuse hindamiseks, aga samuti logi mõistmiseks on vaja ülevaadet, mida logitakse ja arusaamist, mida logikirjed tähendavad. `Korjaja` aitab neid vajadusi rahuldada - sellega, et koostab täieliku, kogu koodibaasi hõlmava nimekirja logimistest ja aitab inimesel siduda logimislausetega inimloetavaid kommentaare. 

Inimene saab logilausete kirjeldustele failis lisada kommentaare tavalise tekstiredaktoriga.

Logilaused on on väljundfailis Go moodulite, pakkide ja funktsioonide kaupa. Igal logilause juures näidatakse funktsioon, kuhu ta kuulub ja reanumber koodifailis.

## Ehitamine

`cd Korjaja`

`go build .`

## Käivitamine

`Korjaja -dir <koodibaasi kaust>`

`<koodibaasi kaust>` on kaust, millest ja mille alamkaustadest logilauseid kogutakse.

## Näide

Kaustas `Example` on logilauseid sisaldav näiterakendus. Sellest logilausete korjamiseks ja faili salvestamiseks:

`Korjaja -dir Example > Example.txt`