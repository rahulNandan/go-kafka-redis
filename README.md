
# Go Redis Kafka Demo

A demo application that interacts with Kafka (producer & consumer) + Redis

This application has 2 parts
1. **http-server**
2. **kafka-consumer**

The http-server handles 2 routes

1. http://127.0.0.1:8080/produce - this produces an event on kafka of the shape `{ type: "number", number: 123333 }` the kafka-consumer listens to this event and determines whether the number is odd/even and based on that increments (`incr`) a key in redis
2. http://127.0.0.1:8080/ - gives you statistics like `{"requests":12,"even":5,"odd":7}` number of events/requests processed so far, number of odd numbers, number of even numbers

  
# To run this project
 
1. ```git clone git@gitlab.com:melwyn95/go-redis-kafka-demo.git```
2. [Install go](https://www.tecmint.com/install-go-in-linux/) & [setup your development environment](https://golang.org/doc/code.html)
3. Install Redis
4. Install Kafka & make a topic call `numbers` (you can name the topic something else also; but make sure the name of the topic is correct in `start.sh`)
5. RUN `./start.sh`

the `start.sh` contains enviroment variables you need to handle it according to you convinence