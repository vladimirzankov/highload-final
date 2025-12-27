import time
import random

from locust import HttpUser, task, between


class APIUser(HttpUser):
    wait_time = between(0, 0)

    @task
    def post_sample(self):
        self.client.post(
            "/metrics",
            json={
                "timestamp": int(time.time()),
                "cpu": random.uniform(10, 90),
                "rps": random.uniform(100, 1000),
            },
            headers={"Content-Type": "application/json"},
        )
