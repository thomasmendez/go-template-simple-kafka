# go-template-simple-kafka

## Running locally

Run `make kafka-cluster` in a terminal window
Open another terminal window and run `make docker`

## Kafka Commands

After running `make kafka-cluster` and making sure container is running you can run the following commands if needed for debugging

Create Kafka Topic

`docker exec go-template-simple-kafka-kafka-1 kafka-topics --create --topic test-topic --bootstrap-server localhost:29092 --replication-factor 1 --partitions 1`

Create Kafka Message (No Key)

`docker exec go-template-simple-kafka-kafka-1 bash -c 'echo "{\"id\":2, \"value\":\"my message with no key\"}" | kafka-console-producer --bootstrap-server localhost:9092 --topic test-topic --property "parse.key=false" --property "key.separator=:"'`

Create Kafka Message (With Key)

`docker exec go-template-simple-kafka-kafka-1 bash -c 'echo "key1:{\"id\":3, \"value\":\"my message with key\"}" | kafka-console-producer --bootstrap-server localhost:9092 --topic test-topic --property "parse.key=true" --property "key.separator=:"'`
