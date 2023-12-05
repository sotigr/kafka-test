from kafka_tools import ConsumerWorker 
from confluent_kafka import Producer 
import time
import random
KAFKA_BROKER = "kafka" 

def on_message(msg):
  if msg.error():
      print("Consumer error: {}".format(msg.error()))
      return 
  time.sleep(1)
  num = random.random()
  print(msg.value().decode('utf-8'), num)
  
  
  # On Successful message reception produce a response message
  
  # p = Producer({'bootstrap.servers': 'kafka'}) 
  # p.produce('myTopicResp', "imporant response msg".encode('utf-8'))
  # p.flush()
  
c1 = ConsumerWorker(on_message, KAFKA_BROKER, "myGroup3", "cp1", ["myTopic"])  
c1.set_task_delay_seconds(3)
c1.set_max_tasks(3)
c1.start(True) 

print("lwol")