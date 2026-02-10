---
title: Atividade 4
author: Pedro H M Tashima
date: 2025-10-09
geometry: "margin=1in"
---

## 1. Escolha ao menos 3 consultas abaixo, e envie 1) o sql e 2) o resultado da consulta no C3SL

### a. Find all pizzerias frequented by at least one person under the age of 18

```sql
SELECT name, pizzeria 
FROM Frequents f 
WHERE f.name IN (SELECT name FROM Person WHERE age < 18);
```

|name|pizzeria      |
|----|--------------|
|Amy |Pizza Hut     |
|Dan |Straw Hat     |
|Dan |New York Pizza|

### b. Find the names of all females who eat either mushroom or pepperoni pizza (or both)

```sql
SELECT DISTINCT(e.name)
FROM Eats e WHERE (e.pizza = "pepperoni" OR e.pizza = "mushroom")
AND e.name IN (SELECT p.name FROM Person p WHERE p.gender = "female");
--------------------------
SELECT e.name FROM Eats e, Person p
WHERE (e.pizza = "pepperoni" OR e.pizza = "mushroom")
AND e.name = p.name AND p.gender = "female";
```

|name|
|----|
|Amy |
|Fay |

### c.  Find the names of all females who eat both mushroom and pepperoni pizza

```sql
SELECT e.name FROM Eats e
WHERE e.pizza = "mushroom"
AND e.name IN (SELECT name FROM Eats WHERE pizza = "pepperoni")
AND e.name IN (SELECT p.name FROM Person p WHERE p.gender = "female");
```

|name|
|----|
|Amy |

### d.  Find all pizzerias that serve at least one pizza that Amy eats for less than $10.00

```sql
SELECT s.pizzeria
FROM Serves s
WHERE s.price < 10
AND s.pizzeria IN (SELECT f.pizzeria FROM Frequents f WHERE f.name = "Amy");
```

|pizzeria |
|---------|
|Pizza Hut|

### e.  Find all pizzerias that are frequented by only females or only males

```sql
SELECT pizzeria
FROM Frequents f, Person p
WHERE p.gender = "female" AND f.name = p.name
AND f.pizzeria NOT IN
 (SELECT f.pizzeria FROM Frequents f, Person p WHERE p.gender = "male" AND f.name = p.name)
UNION
SELECT pizzeria
FROM Frequents f, Person p
WHERE p.gender = "male"
AND f.name = p.name AND f.pizzeria
NOT IN (SELECT f.pizzeria FROM Frequents f, Person p WHERE p.gender = "female" AND f.name = p.name);
```

|pizzeria      |
|--------------|
|Chicago Pizza |
|Little Caesars|
|New York Pizza|

### i.  Find the pizzeria serving the cheapest pepperoni pizza. In the case of ties, return all of the cheapest-pepperoni pizzerias

```sql
SELECT pizzeria
FROM Serves
WHERE pizza = "pepperoni"
AND price = (SELECT s.price
  FROM Serves s
  WHERE s.pizza = "pepperoni"
  ORDER BY s.price ASC
  LIMIT 1);
```

|pizzeria      |
|--------------|
|Straw Hat     |
|New York Pizza|

## 2) Utilizando a Tarefa 3, escreva um parágrafo sobre a qualidade dos dados: quantos arquivos foram utilizados, qual o formato, se havia muitos dados com  'null', se estavam  no formato do dicionário de dados, se o cabeçário possuía palavras em inglês, português ou abreviações, se haviam muitos dados duplicados, se houveram desafios para a inserção, etc

R: O conjunto de dados escolhido com registros de aves no Brasil foi gerado a partir de webscrapping, portanto algumas colunas tem muitos dados, ou todos, com 'null'. O dicionário de dados estava correto e o nome das colunas não possuía abreviações. Todas as colunas tem nomes em português. Não haviam dados duplicados, no entanto, existem registros da mesma ave realizados no mesmo momento por observadores diferentes.

## 3) Utilizando a Tarefa 3, pense em  ao menos 3 perguntas básicas (como a categorização por mês, por tipos) que descrevem características importantes  e curiosas sobre o dado.  Transforme as perguntas para SQL, execute o código no postgres. Escreva um parágrafo com este resultado. Pense também em qual seria a melhor maneira (tabela, gráfico, texto) de incluir este resultado no texto. Como exemplo, considere a seção 3 do artigo [1]

### 3.1 - Qual a ocorrência de registros por ano?

```sql
SELECT date_trunc('year', data_registro) AS ano, COUNT(id) AS total
FROM public.especializacao_pedro_tashima_registros_aves
GROUP BY ano
ORDER BY ano;
```

|ano|total|
|----|-----|
|2019|896  |
|2020|1,192|
|2021|1,697|
|2022|2,399|
|2023|2,161|

![Ocorrencias por ano](anos.eps)

### 3.2 - Quais são as 10 espécies de aves com mais registros?

```sql
select nome_vernaculo, count(id) as total
from public.especializacao_pedro_tashima_registros_aves
group by nome_vernaculo
order by total desc limit 10;
```

|nome_vernaculo      |total|
|--------------------|-----|
|sabiá-laranjeira    |176  |
|tapicuru            |144  |
|sanhaço-papa-laranja|141  |
|periquito-rico      |141  |
|galinha-d'água      |140  |
|marreca-ananaí      |133  |
|quero-quero         |125  |
|caraúna             |115  |
|biguá               |115  |
|carcará             |105  |

![Top 10 especies com mais ocorrencias](top-10-especies.eps)

### 3.3 - Qual a distribuição entre as ações observadas?

```sql
select acao, count(id) as total
from public.especializacao_pedro_tashima_registros_aves
group by acao
order by total desc;
```

|acao                           |total|
|-------------------------------|-----|
|Nenhuma                        |5,983|
| Alimentando-se/Caçando        |1,069|
| Voando                        |342  |
| Nadando                       |302  |
| Outra Ação                    |251  |
| Cantando                      |154  |
| Cuidando/Alimentando Filhote(s|75   |
| Fazendo Ninho                 |57   |
| Chocando                      |28   |
| Bebendo Água                  |19   |
| Cortejando                    |19   |
| Brigando                      |17   |
| Dormindo                      |16   |
| Acasalando                    |9    |
| Regurgitando                  |3    |
| Parasitando outra ave         |1    |

![Ocorrencias de acoes](acoes.eps)
