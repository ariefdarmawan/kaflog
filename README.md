# Kaflog
Kaflog is a hook for EACIIT Toolkit LogEngine to publish log information into Kafka cluster

## How to use
- Initiate Log Engine
  ```go
  log := toolkit.NewLog(true,false,"./logs","","")
  ```
- Hook kafka log
  ```go
  // attach kafka publisher hook to log
   log.AddHook(
        // Host = address of web application
        // Topic = name of Kafka topic
        // Brokers = list of Kafka brokers
        kaflog.Hook(host, topic, brokers...), 

        // type of log will be published to Kafka
        "ERROR", "WARNING")
  ```
  
## Running example
1. Run docker compose to bring up zookeeper and kafka service
   ```bash
   cd docker
   MY_IP={IP} docker-compose up
   ```
2. Use kafkacat, create a consumer to subscribe for topic
   ```sh
   kafkacat -C -b server1:host1,serverN:hostN -t topicname -p 0
   ```
3. Test by using kafkacat to see if respective consumer subscribtion work properly
   ```sh
   echo "Test publish kafka" | kafkacat -P -b server1:host1,serverN:hostN -t topicname -p 0
   ```
4. Attach kaflog hook to Knot Web Application
   ```go
   // prepare log
   log, _ := toolkit.NewLog(true, true, 
        "./logs", "$LOGTYPE_$DATE.log", 
        "yyyyMMdd")

   // attach kafka publisher hook to log
   log.AddHook(
        // Host = address of web application
        // Topic = name of Kafka topic
        // Brokers = list of Kafka brokers
        kaflog.Hook(host, topic, brokers...), 

        // type of log will be published to Kafka
        "ERROR", "WARNING")

    // attach log to knot app
    s := knot.NewServer().SetLogger(log)
   ```
5. Check and monitor