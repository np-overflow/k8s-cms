#
# k8s-cms
# Contest Web Server Testing 
# Submission Script to simulate load
#

import os
import time
import socket
import numpy as np

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

# perform a single submission
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


# wait for random duration sampled from a normal distribtion governed by the args:
# hit_mean - average wait time when submitting
# hit_deviation - standard deviation of wait time before submitting
def random_wait(mean, deviation):
    # wait for a random normal distribution before first submission
    wait_time = np.random.normal(mean, deviation)
    wait_time = max([wait_time, 0])
    time.sleep(wait_time)

# perform the submission test with the given arguments;
# seed - random no. genertor seed.
# hit_mean - average wait time when submitting
# hit_deviation - standard deviation of wait time before submitting
# selenium - the host that exposes a selenium servicef
# port - port of the seleniums service
def main(seed, hit_mean=60, hit_deviation=30, selenium_host="selenium", port=4444):
    # wait for selenium service to become available
    wait_for_port(selenium_host, port)

    # seed random no. generator
    np.random.seed(seed)

    # wait for a random duration before starting to submit
    random_wait(hit_mean, hit_deviation)
    while True:

        try:
            # perform submission hit
            browser = webdriver.Remote(command_executor=f"http://{selenium_host}:{port}/wd/hub",
                                       desired_capabilities=DesiredCapabilities.FIREFOX)
            submit(browser)
            browser.quit()
            print(".", end="", flush=True)

            # wait for a random duration before submitting again
            random_wait(hit_mean, hit_deviation)
        except:
            print("E", end="", flush=True)
            # wait for up to 5s before retrying
            time.sleep(np.random.random() * 5)

if __name__ == "__main__":
    n_processes = 12
    processes = []
    for i in range(n_processes):
        seed = np.random.randint(0, 2**32 -1)
        process = Process(target=main, kwargs={ "seed": seed })
        process.start()
        processes.append(process)

    try:
        while True:
            time.sleep(0.1)
    except KeyboardInterrupt:
        for process in processes:
            process.kill()

