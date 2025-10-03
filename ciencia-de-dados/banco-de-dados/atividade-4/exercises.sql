-- a. 	Find all pizzerias frequented by at least one person under the age of 18.
SELECT name, pizzeria FROM Frequents f WHERE f.name IN (SELECT name FROM Person WHERE age < 18);
/*
|name|pizzeria      |
|----|--------------|
|Amy |Pizza Hut     |
|Dan |Straw Hat     |
|Dan |New York Pizza|
*/

-- b. 	Find the names of all females who eat either mushroom or pepperoni pizza (or both).
SELECT DISTINCT(e.name) FROM Eats e WHERE (e.pizza = "pepperoni" OR e.pizza = "mushroom") 
	AND e.name IN (SELECT p.name FROM Person p WHERE p.gender = "female");
--------------------------
SELECT e.name FROM Eats e, Person p WHERE (e.pizza = "pepperoni" OR e.pizza = "mushroom") AND e.name = p.name AND p.gender = "female";

/*
|name|
|----|
|Amy |
|Fay |
*/

-- c. 	Find the names of all females who eat both mushroom and pepperoni pizza.
SELECT e.name FROM Eats e WHERE e.pizza = "mushroom" AND e.name IN (SELECT name FROM Eats WHERE pizza = "pepperoni")
	AND e.name IN (SELECT p.name FROM Person p WHERE p.gender = "female")
/*
|name|
|----|
|Amy |
*/	
-- d. 	Find all pizzerias that serve at least one pizza that Amy eats for less than $10.00.
SELECT s.pizzeria FROM Serves s WHERE s.price < 10 AND s.pizzeria IN (SELECT f.pizzeria FROM Frequents f WHERE f.name = "Amy");
/*
|pizzeria |
|---------|
|Pizza Hut|
 */

-- e. 	Find all pizzerias that are frequented by only females or only males.
SELECT pizzeria FROM Frequents f, Person p WHERE p.gender = "female" AND f.name = p.name AND f.pizzeria NOT IN 
	(SELECT f.pizzeria FROM Frequents f, Person p WHERE p.gender = "male" AND f.name = p.name)
UNION
SELECT pizzeria FROM Frequents f, Person p WHERE p.gender = "male" AND f.name = p.name AND f.pizzeria NOT IN 
	(SELECT f.pizzeria FROM Frequents f, Person p WHERE p.gender = "female" AND f.name = p.name);
/*
|pizzeria      |
|--------------|
|Chicago Pizza |
|Little Caesars|
|New York Pizza|
 */

-- f. 	For each person, find all pizzas the person eats that are not served by any pizzeria the person frequents. Return all such person (name) / pizza pairs.

-- g. 	Find the names of all people who frequent only pizzerias serving at least one pizza they eat.

-- h. 	Find the names of all people who frequent every pizzeria serving at least one pizza they eat.

-- i. 	Find the pizzeria serving the cheapest pepperoni pizza. In the case of ties, return all of the cheapest-pepperoni pizzerias.
SELECT pizzeria FROM Serves WHERE pizza = "pepperoni" AND price = (SELECT s.price FROM Serves s WHERE s.pizza = "pepperoni" ORDER BY s.price ASC LIMIT 1);

-- |pizzeria      |
-- |--------------|
-- |Straw Hat     |
-- |New York Pizza|
