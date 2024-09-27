import os
import time
import subprocess
import requests
from data import text

source = os.path.dirname(__file__)
parent = os.path.join(source, "../../")
script_directory = os.path.join(parent, "scripts")


def run(command):
    subprocess.run(f"{script_directory}/{command}", shell=True, executable="/bin/bash")


def spin_up_instances():
    run("spinup_bg.sh")


def teardown_instances():
    run("teardown.sh")


def test_replication():
    master_url = "http://127.0.0.1:8100"
    slave_url = "http://127.0.0.1:8101"

    _ = requests.post(f"{master_url}/analysis", json={"text": text.s24u_text})

    master_response = requests.get(f"{master_url}/counts?words={text.s24u_words}")
    slave_response = requests.get(f"{slave_url}/counts?words={text.s24u_words}")

    master_data = master_response.json()
    slave_data = slave_response.json()

    assert master_data != slave_data
    print("Sync has not happened yet. Sleeping for 5 seconds...")
    time.sleep(5)

    slave_response = requests.get(f"{slave_url}/counts?words={text.s24u_words}")
    slave_data = slave_response.json()

    assert master_data == slave_data
    print("Sync has happened!")


spin_up_instances()

print("Testing replication...")
try:
    test_replication()
except AssertionError:
    print("Replication test FAILED")
print("Replication test PASSED")

teardown_instances()
