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

var oneOnOneMailPlaceholderTitle = "Copia mail 1:1 negli appunti"
var oneOnOneMailPlaceholderPlain = `### MM/DD

#### Ostacoli

* *Il rilascio di X è fermo da martedì: aspetto risposta dal partner, intanto valuto un workaround.*

#### Decisioni ferme

* *Migrazione della parte Y: parto adesso o dopo il rilascio di novembre?*

#### Dubbi tecnici

* *Non ho un'idea chiara di come testiamo le nuove API.*

#### Attriti con altri reparti

* *Le richieste del customer care mi arrivano in chat e ne perdo traccia.*

#### Crescita

* *Mi piacerebbe imparare a fare analisi di performance sul database.*


---
`
var oneOnOneMailPlaceholderHTML = `<h3>MM/DD</h3>
<h4>Ostacoli</h4>
<ul>
  <li><p><em>Il rilascio di X è fermo da martedì: aspetto risposta dal partner, intanto valuto un workaround.</em></p></li>
</ul>
<h4>Decisioni ferme</h4>
<ul>
  <li><p><em>Migrazione della parte Y: parto adesso o dopo il rilascio di novembre?</em></p></li>
</ul>
<h4>Dubbi tecnici</h4>
<ul>
  <li><p><em>Non ho un'idea chiara di come testiamo le nuove API.</em></p></li>
</ul>
<h4>Attriti con altri reparti</h4>
<ul>
  <li><p><em>Le richieste del customer care mi arrivano in chat e ne perdo traccia.</em></p></li>
</ul>
<h4>Crescita</h4>
<ul>
  <li><p><em>Mi piacerebbe imparare a fare analisi di performance sul database.</em></p></li>
</ul>
<hr>`
