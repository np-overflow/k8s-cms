#
# k8s-cms
# Contest Web Server Testing 
# Submission Script
#

import os
import time
import socket

from random import random
from multiprocessing import Process

from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities

# waits for the given tcp port to start accepting connections by polling every
# interval seconds
# waits for a maxmium of timeout seconds
def wait_for_port(host, port, timeout=30.0, interval=0.1):
    begin_timestamp = time.time()
    while (time.time() - begin_timestamp) <= timeout:
        try:
            sock = socket.create_connection((host, port), timeout=timeout)
            sock.close()
        except OSError:
            # still no respose, continue waiting
            time.sleep(interval)
            continue
        # successfully connected to the port
        return
    # timeout
    raise TimeoutError(f"Wait for port {host}:{port} timed out.")

# test if we logged in 
def is_login(browser):
    cookies = browser.get_cookies()
    login_cookies = filter((lambda cookie: "_login"  in cookie["name"]), cookies)
    return True if len(list(login_cookies)) >= 1 else False

def wait(browser, timeout, predicate):
    try:
        WebDriverWait(browser, timeout).until(predicate)
    except TimeoutError:
        browser.quit()
        raise TimeoutError

# perform login to the contest site
def login(browser):
    if is_login(browser): return # alreadly logged in

    # get contest page
    browser.get("https://demo-npio.np-overflow.club/test")

    # Log In
    wait(browser, 10, EC.presence_of_element_located((By.ID, "username")))
    browser.find_element_by_id("username").send_keys("test" + Keys.TAB)
    browser.find_element_by_id("password").send_keys("test" + Keys.ENTER)

    wait(browser, 10, EC.presence_of_element_located((By.ID, "countdown_box")))
    assert is_login(browser)

# submission test
def submit(browser):
    # get contest page
    browser.get("https://demo-npio.np-overflow.club/test")

    # if we are not already login, perform login first
    if not is_login(browser): login(browser)

    # Load submission page
    browser.get("https://demo-npio.np-overflow.club/test/tasks/test/submissions")
    wait(browser, 10, EC.presence_of_element_located((By.CLASS_NAME, "task_submissions")))

    # Submit
    browser.find_element_by_id("input0").send_keys("/home/seluser/project/test.c")
    browser.find_element_by_class_name("btn-success").click()

def main():
    # wait for selenium service to become available
    wait_for_port("selenium", 4444)


    # Load Browser
    success = True
    while True:
        if success == True:
            time.sleep(random() * 60)

        try:
            browser = webdriver.Remote(command_executor='http://selenium:4444/wd/hub',
                                       desired_capabilities=DesiredCapabilities.FIREFOX)
            submit(browser)
            browser.quit()
        except:
            print("error")
            success = False
        finally:
            print("submitted.")
            success = True

if __name__ == "__main__":
    processes = [ Process(target=main) for i in range(16) ]
    for process in processes:
        process.start()

    try:
        while True:
            time.sleep(0.1)
    except KeyboardInterrupt:
        for process in processes:
            process.kill()

