#
# k8s-cms
# Contest Web Server Testing 
# Submission Script to simulate load
#

import os
import sys
import time
import socket
import numpy as np
import traceback
import socket

from multiprocessing import Process
from argparse import ArgumentParser

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

# perform login to the contest site with given credentia
# waits for a maxmium of timeout seconds giving up trying to login
def login(browser, username, password, timeout=10):
    if is_login(browser): return # alreadly logged in

    # Log In
    wait(browser, timeout, EC.presence_of_element_located((By.ID, "username")))
    browser.find_element_by_id("username").send_keys(username + Keys.TAB)
    browser.find_element_by_id("password").send_keys(password + Keys.ENTER)

    wait(browser, timeout, EC.presence_of_element_located((By.ID, "countdown_box")))
    assert is_login(browser)

# perform a single submission to target_url using the given contest and task
# waits for a maxmium of timeout seconds giving up trying to submit
def submit(browser, target_url, contest, task, timeout=10):
    # get contest page
    browser.get(f"{target_url}/{contest}")
    # if we are not already authenticated, perform login first
    if not is_login(browser): login(browser, args.username, args.password, timeout)

    # Load submission page
    browser.get(f"{target_url}/{contest}/tasks/{task}/submissions")
    wait(browser, timeout, EC.presence_of_element_located((By.CLASS_NAME, "task_submissions")))

    # Submit
    browser.find_element_by_id("input0").send_keys("/home/seluser/project/test.c")
    browser.find_element_by_class_name("btn-success").click()

# compute wait time for given parameters:
# hit_mean - average wait time when submitting
# hit_deviation - standard deviation of wait time before submitting
# returns the duration of time in seconds waited
def random_wait(mean, deviation, verbose=False):
    # wait for a random normal distribution before first submission
    wait_time = np.random.normal(mean, deviation)
    wait_time = max([wait_time, 0])

    return wait_time

# perform the submission test with the given arguments:
# process_id - id of the process running
# seed - random no. genertor seed.
# args.hit_mean - average wait time when submitting
# args.hit_deviation - standard deviation of wait time before submitting
# args.target_url - url of the contest web server to target
# args.selenium_host - the host that exposes a selenium servicef
# args.selenium_port - port of the seleniums service
def main(process_id, seed, args):
    prefix = f"{socket.getfqdn()}: Process {process_id}: "

    # seed random no. generator
    np.random.seed(seed)

    # wait for a random duration before starting to submit
    while True:
        # wait for a random duration before submitting again
        wait_time = random_wait(args.hit_mean, args.hit_deviation, args.verbose)
        if args.verbose:
            print("{} waiting for {:0.2f}s".format(prefix, wait_time), flush=True)
        time.sleep(wait_time)

        # perform submission hit
        browser = webdriver.Remote(
            command_executor=f"http://{args.selenium_host}:{args.selenium_port}/wd/hub",
            desired_capabilities=DesiredCapabilities.FIREFOX)
        try:
            submit(browser, args.target_url, args.contest, args.task, args.timeout)
            if args.verbose: print(f"{prefix} sent submission", flush=True)
        except Exception as e:
            print(f"{prefix} failed to send submission", flush=True)
            print(traceback.format_exc(), flush=True)
        finally:
            browser.quit()

# parse command line arguments for submit.py
def parse_args():
    parser = ArgumentParser(description="Simulates a submission load test on CMS")

    # options
    parser.add_argument("-v", "--verbose", action="store_true",
                        help="Enable verbose debugging output")
    parser.add_argument("-u", "--username", dest="username",
                        help="Username to use when authenticating with contest server", default="test")
    parser.add_argument("-p", "--password", dest="password",
                        help="password to use when authenticating with contest server", default="test")
    parser.add_argument("-n", "--n-users", dest="processes", type=int,
                        help="No. of users to simulate", default=4)
    parser.add_argument("--hit-avg", type=float, dest="hit_mean",
                        help="Average secs to wait before submitting", default=60)
    parser.add_argument("--hit-stddev", type=float, dest="hit_deviation",
                        help="Standard deviation of secs to wait before submitting", default=30)
    parser.add_argument("-c", "--contest", help="The contest to select when testing",
                        dest="contest", default="test")
    parser.add_argument("-t", "--task", help="The task to select when testing",
                        dest="task", default="test")
    parser.add_argument("--selenium-port", type=int, dest="selenium_port",
                        help="Port to use when talking to selenium", default=4444)
    parser.add_argument("-w", "--timeout", type=int,
                        help="Maximum seconds to wait when testing before timing out",
                        dest="timeout", default=10)

    # required arguments
    parser.add_argument("selenium_host",
                        help="Selenium server to send requests to")
    parser.add_argument("target_url", help="The url of the contest server to send load.")

    return parser.parse_args()

# health check status path
# created when everything is functioning normally
HEALTH_CHECK_PATH="/tmp/healthz"

if __name__ == "__main__":
    args = parse_args()

    # wait for selenium service to become available
    if args.verbose:
        print(f"{socket.getfqdn()}: waiting for selenium service on: "
              f"{args.selenium_host}:{args.selenium_port}", flush=True)
    wait_for_port(args.selenium_host, args.selenium_port)

    # start proccesses to simulate uses
    n_processes = args.processes
    processes = []
    for i in range(1, n_processes + 1):
        seed = np.random.randint(0, 2**32 -1)
        process = Process(target=main, kwargs={ "process_id": i,
                                                "seed": seed,
                                                "args": args })
        process.start()
        if args.verbose:
            print(f"{socket.getfqdn()} started user process {i}", flush=True)
        processes.append(process)

    # notify that we are up and running
    open(HEALTH_CHECK_PATH, "w").close()

    is_healthy = True
    while is_healthy:
        time.sleep(0.5)
        # check the health of the given processes
        is_healthy = all([ process.is_alive() for process in processes ])
        heath_status = "healthy" if is_healthy else "unhealthy"


    if not is_healthy:
        # something bad happened - signal to health check
        os.remove(HEALTH_CHECK_PATH)

        # cleanup worker processes
        for process in processes:
            process.kill()

        sys.exit(1)
