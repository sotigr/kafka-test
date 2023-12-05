from confluent_kafka import Consumer 
from threading import Thread, Lock
import time

lock = Lock()
  
def thread_callback(self, cb, msg): 
  if msg == None:
    self.cur_tasks -= 1
    return 
  cb(msg)
  with lock:
    self.cur_tasks -= 1
 
def runner(self, cb): 
  while self.running:  
    msg = self.c.poll(1.0)

    if self.max_tasks >= 0:
       
      if self.cur_tasks >= self.max_tasks:
        if self.delay > 0: 
          time.sleep(self.delay) 
        continue
        
      self.cur_tasks += 1
      t = Thread(target=thread_callback, args=(self, cb, msg)) 
      t.start()
    else:
      cb(msg)  
      if self.delay > 0:
        time.sleep(self.delay) 
        

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
    
    self.max_tasks = -1
    self.cur_tasks = 0
    self.delay = 0
    
    self.running = False
    self.t = Thread(target=runner, args=(self, on_message))
    self.on_message = on_message
    
  def set_max_tasks(self, max_tasks):
    self.max_tasks = max_tasks
  
  def set_task_delay_seconds(self, delay):
    self.delay = delay

  def start(self, block):
    if self.running:
      return 
     
    self.running = True
    if block: 
      runner(self, self.on_message)
    else: 
      self.t.start()
    
  def stop(self):
    self.running = False 
    self.c.close()
