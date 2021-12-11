import json 
from random import randrange 

from locust import HttpUser, between, task 

debug = False

def printDebug(msg):
    if debug:
        print(msg)

class Reader():

    def __init__(self):
        self.array = []

    def pickRandom(self):
        length = len(self.array)

        if length > 0:
            random_index = randrange(0, length - 1) if length > 1 else 0
            return self.array.pop(random_index)
        else:
            print(">> Reader: not founded value in the file")
            return None 

    def load(self):
        print(">> Reader: Initializing read file")
        try:
            with open("traffic.json", "r") as data_file:
                self.array = json.loads(data_file.read())
        except Exception as e:
            print(f">> Reader: {e}")


class MessageTraffic(HttpUser):
    wait_time = between(0.1, 0.9)
    reader = Reader() 
    reader.load()

    def on_start(self):
        print(">> MessageTraffic: Initialized send traffic")

    @task 
    def PostMessage(self):
        random_data = self.reader.pickRandom()

        if random_data:
            data_to_send = json.dumps(random_data)
            printDebug(data_to_send)

            self.client.post("/", json=random_data)
        else :
            print(">> MessageTraffic: Finalized send traffic")
            self.stop(True)

    @task 
    def GetMessages(self):
        self.client.get("/")