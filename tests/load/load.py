import os
import time
import subprocess
import requests
from data import text
import locust

source = os.path.dirname(__file__)
parent = os.path.join(source, "../../")
script_directory = os.path.join(parent, "scripts")


def run(command):
    subprocess.run(f"{script_directory}/{command}", shell=True, executable="/bin/bash")


def spin_up_instances():
    run("spinup_bg.sh")


def teardown_instances():
    run("teardown.sh")


def run_load_test():
    run("load.sh")


spin_up_instances()
run_load_test()
teardown_instances()
