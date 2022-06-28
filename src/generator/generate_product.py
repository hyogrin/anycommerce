# -*- coding: utf-8 -*-
import os
import json
import time
from urllib.request import urlretrieve
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
from random import randrange

# https://sa-na.tistory.com/entry/Selenium%EC%9D%84-%EC%82%AC%EC%9A%A9%ED%95%B4-%EB%84%A4%EC%9D%B4%EB%B2%84-%EC%87%BC%ED%95%91-%ED%81%AC%EB%A1%A4%EB%A7%81-%ED%95%98%EA%B8%B0

# Set Shop URL (Naver Shopping for example)
products = []
category_id = '100'
category_name = 'all'
shop_url = 'https://search.shopping.naver.com/best/category/click?categoryCategoryId=ALL&categoryDemo=A00&categoryRootCategoryId=ALL&chartRank=1&period=P7D%27'

os.makedirs("./products_data/" + category_id + "_" +
            category_name + "/images/", exist_ok=True)

# Install Chrome Driver
driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()))
driver.implicitly_wait(5)
driver.get(shop_url)

SCROLL_PAUSE_TIME = 0.5

# Get scroll height
last_height = driver.execute_script("return document.body.scrollHeight")

while True:
    # Scroll down to bottom
    driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")

    # Wait to load page
    time.sleep(SCROLL_PAUSE_TIME)

    # Calculate new scroll height and compare with last scroll height
    new_height = driver.execute_script("return document.body.scrollHeight")
    if new_height == last_height:
        break
    last_height = new_height

# time.sleep(5)

# Find product element
product_lists = driver.find_elements(
    by=By.CSS_SELECTOR, value='#__next > div > div > div > div > div > div.category_panel > div > ul > li')

for i, product_list in enumerate(product_lists):
    product = {}

    # Find Title and price
    product_title = product_list.find_element(
        by=By.XPATH, value='.//div[2]/div[2]')
    product_price = product_list.find_element(
        by=By.XPATH, value='.//div[2]/div[1]/strong')

    product['id'] = i + 1 + int(category_id)
    product['category'] = category_name
    product['name'] = product_title.text
    product['price'] = int(product_price.text.replace(',', ''))
    product['code'] = product_list.get_attribute('id')
    products.append(product)

    # Download the image
    product_img = product_list.find_element(
        by=By.XPATH, value='.//div[1]/div[2]/img')
    src = product_img.get_attribute('src')
    urlretrieve(src, "products_data/" +
                category_id + "_" + category_name + "/images/" + str(i + 1 + int(category_id)) + ".jpg")

with open("products_data/" + category_id + "_" + category_name + "/" + 'products_data.json', 'w', encoding='UTF-8') as json_file:
    json.dump(products, json_file, ensure_ascii=False)
