import time
from locust import HttpUser, task, between
from data import text


class QuickstartUser(HttpUser):
    wait_time = between(1, 1)

    @task(1)
    def get_counts(self):
        self.client.get(f"/counts?words={text.s24u_words}", name="/counts?words=...")
