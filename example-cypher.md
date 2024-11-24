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

#### Query what tier a sample customer belongs to based upon the lifetime_rewards_variable components
```
MATCH (c:sample_customer)-[:placed]->(o:order)-[r:has_item]->(p:product), (l:lifetime_rewards_variable)
WITH l, c, ROUND(SUM(r.price) * 100) / 100 AS purchase_amount
RETURN
CASE 
    WHEN purchase_amount > l.average_purchase_amount + (2 * l.stddev_purchase_amount) THEN "Diamond"
    WHEN purchase_amount > l.average_purchase_amount + l.stddev_purchase_amount THEN "Gold"
    WHEN purchase_amount >= l.average_purchase_amount THEN "Silver"
    ELSE "Member"
END AS tier

```
