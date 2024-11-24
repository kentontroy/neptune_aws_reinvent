#### Identify the top 50 customers by purchase amount
```
MATCH (c:customer)-[:ordered]->(o:order)-[r:has_item]->(p:product)
WITH c.customer_id AS customer_id, ROUND(SUM(r.price) * 100) / 100 as purchase_amount
ORDER BY purchase_amount DESC
LIMIT 50
WITH COLLECT(customer_id) AS top_customers
UNWIND top_customers AS customer_id

MERGE (t:top_customer {customer_id: customer_id})
```

```
MATCH (c:top_customer), (o:order)
WHERE c.customer_id = o.customer_id
MERGE (c)-[:placed]->(o)
RETURN c, o
```

#### List what products those top 50 customers have purchased
```
MATCH (c:top_customer)-[i:placed]->(o:order)-[r:has_item]->(p:product)
RETURN c.customer_id, 
    COLLECT({
        product: p.product_category_name, amount: ROUND(r.price * 100) / 100
    }) AS purchased_items
```

#### Create a random sample of 50 customers
```
MATCH (c:customer)
WITH c, rand() AS randomValue
ORDER BY randomValue
LIMIT 50
WITH COLLECT(c.customer_id) AS sample_customers
UNWIND sample_customers AS customer_id

MERGE (s:sample_customer {customer_id: customer_id})
```

```
MATCH (c:sample_customer), (o:order)
WHERE c.customer_id = o.customer_id
MERGE (c)-[:placed]->(o)
RETURN c, o
```
