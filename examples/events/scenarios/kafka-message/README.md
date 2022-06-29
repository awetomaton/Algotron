# Kafka Trigger example

##Install event-source and sensor

```
kubectl apply -f ../../event-sources/kafka-eventsource.yaml
kubectl apply -f ../../sensors/kafka-sensor.yaml
```

## Exec into kafka pod to run kafka sh scripts
```
kubectl exec -n algnext-ns -it algnext-kafka-0 -- bash
```

## Create kafka topic pointing to the bootstrap server 
```
kafka-topics.sh --create --topic quickstart-events --partitions 10 --bootstrap-server algnext-kafka:9092
```

## Produce messages
Run the simple-kafka-message-loop.sh from inside the kafka container

## Verify Workflows Trigger
# Watch Event Flow in Argo UI to see workflows running

