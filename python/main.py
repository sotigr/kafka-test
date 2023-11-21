from kafka_tools import ConsumerWorker 
from confluent_kafka import Producer 

KAFKA_BROKER = "kafka"

def on_message(msg):
  if msg.error():
      print("Consumer error: {}".format(msg.error()))
      return
  print(msg.value().decode('utf-8'))
  
  # On Successful message reception produce a response message
  p = Producer({'bootstrap.servers': 'kafka'}) 
  p.produce('myTopicResp', "imporant response msg".encode('utf-8'))
  p.flush()
  
c1 = ConsumerWorker(on_message, KAFKA_BROKER, "myGroup3", "cp1", ["myTopic"])  
c1.start()