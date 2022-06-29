for i in {1..5}; do echo {'"cust_id"': $i, '"month"': 9, '"amount_paid"':457.78} | kafka-console-producer.sh --topic quickstart-events-2 --bootstrap-server algnext-kafka:9092; done
