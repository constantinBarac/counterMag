import time
from locust import HttpUser, task, between
from data import text


class QuickstartUser(HttpUser):
    wait_time = between(1, 1)

    @task(1)
    def analyze_text(self):
        self.client.post("/analysis", name="/analysis", json={"text": text.s24u_text})
