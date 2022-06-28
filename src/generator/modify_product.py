# -*- coding: utf-8 -*-
import json
import time
import csv
from urllib.request import urlretrieve
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
from random import randrange

products = []

with open("./products.json", 'r', encoding='UTF-8-sig') as f:
    products = json.load(f)

# for product in products:
#     product['image'] = str(product['id']) + ".jpg"

# with open("./generator/products.json", 'w', encoding='UTF-8-sig') as json_file:
#     json.dump(products, json_file, ensure_ascii=False)

with open("products.csv", 'w', encoding='UTF-8-sig') as csv_file:
    keys = [k for k, v in products[7].items()]
    writer = csv.DictWriter(csv_file, fieldnames=keys)
    writer.writeheader()
    writer.writerows(products)
