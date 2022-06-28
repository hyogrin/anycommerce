# -*- coding: utf-8 -*-
import json
import csv
import time
from urllib.request import urlretrieve
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
from random import randrange

users = []

with open("users.json", 'r', encoding='UTF-8-sig') as f:
    users = json.load(f)

# for user in users:
    # user['segment'] = user['gender'] + str(int(user['age']/10) * 10)
    # user['email'] = user['username'] + "@example.com"
    # user['first_name'] = user['name'][0]
    # user['last_name'] = user['name'][1:]
    # user['id'] = str(user['id'])

# with open("users.json", 'w', encoding='UTF-8-sig') as json_file:
#     json.dump(users, json_file, ensure_ascii=False)

with open("users.csv", 'w', encoding='UTF-8-sig') as csv_file:
    keys = [k for k, v in users[0].items()]
    writer = csv.DictWriter(csv_file, fieldnames=keys)
    writer.writeheader()
    writer.writerows(users)
