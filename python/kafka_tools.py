from confluent_kafka import Consumer 
from threading import Thread
  
def runner(self,  cb): 
  while self.running: 
    msg = self.c.poll(1.0)
    if msg is None:
        continue
    else:
      cb(msg)

class ConsumerWorker: 
  def __init__(self, on_message, servers, group_id, instance_id, topics): 
    c = Consumer({
        "bootstrap.servers": servers,
        "group.id": group_id, 
        "auto.offset.reset": "earliest",
        "group.instance.id": instance_id,
    })  
    c.subscribe(topics) 
    self.c = c
    self.running = False
    self.t = Thread(target=runner, args=(self, on_message))
    self.on_message = on_message

  def start(self):
    if self.running:
      return 
    self.running = True
    self.t.start()
    
  def stop(self):
    self.running = False 
    self.c.close()
