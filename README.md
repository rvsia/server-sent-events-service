# Hackaton project - kafka to frontend

This project connect to kafka based on specific topics described in topic.json it then sends event notifications trough server side events to registered clients.

Client needs to register with account_number in path, we'll then look into kafka message and send the message to clients with matching account_number.

## Running

First spawn local zookeeper and kafka by running `scripts/start-kafka.sh` you can then build this application. If you want to receive messages simply call `curl "localhost:3000/api/notifier/v1/connect?room=inventory&account_number=55"`.

Once kafka receives message on any topic listed in topics.json you will be notified over SSE.

To test locally producing of messages you can build and run producer_example with `./producer_example localhost:9092 platform.inventory.events`
