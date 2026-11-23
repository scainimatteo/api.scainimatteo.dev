package outline

var monthTablePlaceholderTitle = "Copia template del mese negli appunti"
var monthTablePlaceholderPlain = `Descrizione
Categoria
Importo
Eseguito
Bolletta internet "{MM/YY}"
Bollette
50
`
var monthTablePlaceholderHTML = `<table>
  <colgroup>
    <col style="width: 421px;">
    <col style="width: 123px;">
    <col style="width: 96px;">
    <col style="width: 96px;">
  </colgroup>
  <tbody>
    <tr>
      <th data-colwidth="421"><p dir="auto">Descrizione</p></th>
      <th data-colwidth="123"><p dir="auto">Categoria</p></th>
      <th data-colwidth="96"><p dir="auto">Importo</p></th>
      <th data-colwidth="96"><p dir="auto">Eseguito</p></th>
    </tr>
    <tr>
      <td data-colwidth="421"><p dir="auto">Bolletta internet "{MM/YY}"</p></td>
      <td data-colwidth="123"><p dir="auto">Bollette</p></td>
      <td data-colwidth="96"><p dir="auto">50</p></td>
      <td data-colwidth="96"><p dir="auto"></p></td>
    </tr>
  </tbody>
</table>`

var oneOnOneOdgPlaceholderTitle = "Copia OdG 1:1 negli appunti"

// true: l'1:1 cade il mercoledì delle settimane ISO dispari; false: delle settimane pari
var oneOnOneOddWeek = true

var oneOnOneOdgPlaceholderPlain = `## "{MM/DD}"

### Ostacoli

* *Il rilascio di X è fermo da martedì: aspetto risposta dal partner, intanto valuto un workaround.*

### Decisioni ferme

* *Migrazione della parte Y: parto adesso o dopo il rilascio di novembre?*

### Dubbi tecnici

* *Non ho un'idea chiara di come testiamo le nuove API.*

### Attriti con altri reparti

* *Le richieste del customer care mi arrivano in chat e ne perdo traccia.*

### Crescita

* *Mi piacerebbe imparare a fare analisi di performance sul database.*


---
`
var oneOnOneOdgPlaceholderHTML = `<h2>"{MM/DD}"</h2>
<h3>Ostacoli</h3>
<ul>
  <li><p><em>Il rilascio di X è fermo da martedì: aspetto risposta dal partner, intanto valuto un workaround.</em></p></li>
</ul>
<h3>Decisioni ferme</h3>
<ul>
  <li><p><em>Migrazione della parte Y: parto adesso o dopo il rilascio di novembre?</em></p></li>
</ul>
<h3>Dubbi tecnici</h3>
<ul>
  <li><p><em>Non ho un'idea chiara di come testiamo le nuove API.</em></p></li>
</ul>
<h3>Attriti con altri reparti</h3>
<ul>
  <li><p><em>Le richieste del customer care mi arrivano in chat e ne perdo traccia.</em></p></li>
</ul>
<h3>Crescita</h3>
<ul>
  <li><p><em>Mi piacerebbe imparare a fare analisi di performance sul database.</em></p></li>
</ul>
<hr>`

var oneOnOneMailPlaceholderTitle = "Copia mail 1:1 negli appunti"
var oneOnOneMailPlaceholderPlain = `OdG del meeting 1:1 del "{DD/MM}"
Ciao Andrea,
di seguito ti mando l'OdG del meeting 1:1 di mercoledì.

[OdG]

A mercoledì!
Matteo
`
var oneOnOneMailPlaceholderHTML = `<p>OdG del meeting 1:1 del "{DD/MM}"</p>
<p>Ciao Andrea,<br>di seguito ti mando l'OdG del meeting 1:1 di mercoledì.</p>
<p>[OdG]</p>
<p>A mercoledì!<br>Matteo</p>`

var weeklyActivitiesPlaceholderTitle = "Copia attività settimanali negli appunti"
var weeklyActivitiesPlaceholderPlain = `## "{MONDAY}" → "{FRIDAY}"

- [ ] _Lista di attività settimanali_

### Weekly Dev Sync "{MONDAY}"
- _Lista di note del Weekly Dev Sync_
`
var weeklyActivitiesPlaceholderHTML = `<h2>"{MONDAY}" → "{FRIDAY}"</h2>
<ul class="checkbox_list">
  <li data-type="checkbox_item"><p><em>Lista di attività settimanali</em></p></li>
</ul>
<h3>Weekly Dev Sync "{MONDAY}"</h3>
<ul>
  <li><p><em>Lista di note del Weekly Dev Sync</em></p></li>
</ul>`
