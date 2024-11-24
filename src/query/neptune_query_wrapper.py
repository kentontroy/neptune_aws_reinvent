import logging
import neptune_query_strings as q_strings
import os
import random

from neo4j import GraphDatabase
from typing import Dict, List

logging.basicConfig(
    level=logging.INFO, 
    format="%(asctime)s - %(levelname)s \n %(message)s" 
)
logger = logging.getLogger(__name__)

USERNAME = os.getenv("AWS_NEPTUNE_USERNAME")
PASSWORD = os.getenv("AWS_NEPTUNE_PASSWORD")
ENDPOINT = os.getenv("AWS_NEPTUNE_ENDPOINT")
PORT = os.getenv("AWS_NEPTUNE_ENDPOINT_PORT")
URI = "bolt://{0}:{1}".format(ENDPOINT, PORT)

def getSampleCustomerIds()->List[str]:
    customer_ids = []
    try:
        with GraphDatabase.driver(URI, auth=(USERNAME, PASSWORD), encrypted=True) as driver:
            driver.verify_connectivity()
            drs = driver.session()
            res = drs.run(q_strings.GET_DEMO_CUSTOMER_IDS)
            for rec in res:
                customer_ids.append(rec["id"])
            drs.close()

    except Exception as e:
        raise Exception(f"A Neptune error occurred in getSampleCustomerIds: {e}")
    
    return customer_ids

def getTierForCustomerId(customer_id: str)->str:
    tier = ""
    try:
        with GraphDatabase.driver(URI, auth=(USERNAME, PASSWORD), encrypted=True) as driver:
            driver.verify_connectivity()
            drs = driver.session()
            res = drs.run(q_strings.GET_TIER, customer_id=customer_id)
            for rec in res:
                tier = rec["tier"]            
                break
            drs.close()

    except Exception as e:
        raise Exception(f"A Neptune error occurred in getTierForCustomerId: {e}")

    return tier

def getTiersForAllSampleCustomers()->Dict[str, str]:
    tiers = {}
    try:
        with GraphDatabase.driver(URI, auth=(USERNAME, PASSWORD), encrypted=True) as driver:
            driver.verify_connectivity()
            drs = driver.session()
            res = drs.run(q_strings.GET_TIERS_FOR_ALL_SAMPLES)
            for rec in res:
                tiers[rec["customer_id"]] = rec["tier"]
            drs.close()

    except Exception as e:
        raise Exception(f"A Neptune error occurred in getTiersForAllSampleCustomers: {e}")

    return tiers

if __name__ == "__main__":
    customer_ids = getSampleCustomerIds()
    logger.info(customer_ids)

    tier = getTierForCustomerId(random.choice(customer_ids))
    logger.info(tier)

    tiers = getTiersForAllSampleCustomers()
    logger.info(tiers)

