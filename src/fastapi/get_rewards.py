from neo4j import GraphDatabase

ENDPOINT = "reinvent-db-neptune-1.cluster-cnlncb8narog.us-east-2.neptune.amazonaws.com"
PORT = 8182
URI = "bolt://{0}:{1}".format(ENDPOINT, PORT)

with GraphDatabase.driver(uri, auth=("username", "password"), encrypted=True) as driver:
    driver.verify_connectivity()
    drs = driver.session()
    res = drs.run("MATCH (t:demo_set_customer) RETURN t.customer_id AS id")
    for rec in res:
        print(rec)

if __name__ == "__main__":
    main()
