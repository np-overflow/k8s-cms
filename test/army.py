from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.keys import Keys
import os

def main():
    # Load Browser
    browser = webdriver.Firefox()
    browser.get("http://demo-npio.np-overflow.club/OverloadingServer")

    # Log In
    browser.find_element_by_id("username").send_keys("Army")
    browser.find_element_by_id("password").send_keys("army" + Keys.ENTER)

    # Load Testing Page 
    try:
        WebDriverWait(browser, 10).until(EC.presence_of_element_located((By.LINK_TEXT, "Submissions"))).click()
    except:
        browser.quit()
        print("Page not loading")

    # Submit
    browser.find_element_by_id("input0").send_keys(str(os.getcwd()) + "/num.c")
    try:
        WebDriverWait(browser, 10).until(EC.presence_of_element_located((By.CLASS_NAME, "btn-success"))).click()
    except:
        browser.quit()
        print("Page not loading")