# -*- coding: utf-8 -*-
import json
import time
from urllib.request import urlretrieve
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
from random import randrange

products = []

with open("./generator/products/all/products.json", 'r', encoding='UTF-8-sig') as f:
    products = json.load(f)

driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()))
driver.implicitly_wait(10)
driver.get('https://search.shopping.naver.com/best/category/click?categoryCategoryId=ALL&categoryDemo=F05&categoryRootCategoryId=ALL&chartRank=1&period=P7D')
time.sleep(10)

product_lists = driver.find_elements(
    by=By.CSS_SELECTOR, value='#__next > div > div > div > div > div > div.category_panel > div > ul > li')

affinity_target = 'F50'

for i, product_list in enumerate(product_lists):
    code = product_list.get_attribute('id')

    for p in products:
        if p.get('code') == code:
            p['affinity'].append(affinity_target)

with open("./generator/products/all/products.json", 'w', encoding='UTF-8-sig') as json_file:
    json.dump(products, json_file, ensure_ascii=False)
