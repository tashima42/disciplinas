# Atividade 2

| A  | B  | C  | D  | E  |
|----|----|----|----|----|
| a1 | b1 | c1 | d1 | e1 |
| a1 | b2 | c2 | d2 | e1 |
| a2 | b1 | c3 | d3 | e1 |
| a2 | b1 | c4 | d3 | e1 |
| a3 | b2 | c5 | d1 | e1 |

#### Analise as dependencias funcionais

As dependencias válidas para esta relação são as seguintes:

* **AB -> D**: Sempre que A e B se repetem, D tem o mesmo valor
* **C -> BDE**: Todos os valores de C são diferentes
* **A -> E**: Sempre que os valores de A se repetem, E tem o mesmo valor (e1)
* **CD -> B**: Nenhuma combinação de CD se repete
