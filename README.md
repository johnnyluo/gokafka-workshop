# gokafka-workshop

This is a talk about golang and kafka prepared for the golang sydney meetup 

* Please install latest golang, refer to [Getting Started](https://golang.org/doc/install)
* We will be running a one node kafka in docker composer on your local machine, if you don't have docker composer setup , please follow [Install Docker Compose](https://docs.docker.com/compose/install/) , to set it up.

You will need the following libraries in this workshop
* github.com/Shopify/sarama
```bash
go get -u github.com/Shopify/sarama
```
* github.com/bsm/sarama-cluster
```bash
go get -u github.com/bsm/sarama-cluster
```

# Run kafka in docker
* open a terminal
* go to kafka-local folder
```bash
docker-composer up -d
```
If you are comfortable to cli tools, then you can go [here](http://kafka.apache.org/downloads.html)

If you are more comfortable to user UI,then you could use the following steps to setup your kafka manager. 

# setup kafka manager
* Open a browser  , go to `http://localhost:9000`
* You should see a screen like this
![](assets/Kafka_Manager.png)
* Click Cluster->Add Cluster
![](assets/Add_Cluster.png)
* What is my zookeeper ip?
```
docker network inspect kafka-local_default
```
```json
[
    {
        "Name": "kafka-local_default",
        "Id": "baf73f678775f5a3903f024ecebe5123c9de83220a4db89e44e06d585a579203",
        "Created": "2018-06-25T11:01:24.434971611Z",
        "Scope": "local",
        "Driver": "bridge",
        "EnableIPv6": false,
        "IPAM": {
            "Driver": "default",
            "Options": null,
            "Config": [
                {
                    "Subnet": "172.21.0.0/16",
                    "Gateway": "172.21.0.1"
                }
            ]
        },
        "Internal": false,
        "Attachable": true,
        "Ingress": false,
        "ConfigFrom": {
            "Network": ""
        },
        "ConfigOnly": false,
        "Containers": {
            "20070212faccd1a3b47067032e4bfa2a7de6b8645827cb1f89d28f2aa9715314": {
                "Name": "kafka-local_kafka_1",
                "EndpointID": "cd697e1434e539e8685ecc877b240b59d0b39eb23ed1ad9ba51df36f4969cc8a",
                "MacAddress": "02:42:ac:15:00:04",
                "IPv4Address": "172.21.0.4/16",
                "IPv6Address": ""
            },
            "7ea6b5848f45e205c46a77faabc5845248ca3f736be58639eda7760f5c590b0f": {
                "Name": "kafka-local_zookeeper_1",
                "EndpointID": "747ebf17af4e1aa2de8a4c7a3c673af7f1f60bca4c9016959f5bfb92c25ac9bc",
                "MacAddress": "02:42:ac:15:00:03",
                "IPv4Address": "172.21.0.3/16",
                "IPv6Address": ""
            },
            "cbe0d113e2d1a972aee79f7ab33f1154f353de80ae407d2e0996be45187d6ad4": {
                "Name": "kafka-local_kafkamgr_1",
                "EndpointID": "6f09ec62a3fe15b3c03d42bb9f5e725855f8e7d3396f10c6993dcd7d4d67e445",
                "MacAddress": "02:42:ac:15:00:02",
                "IPv4Address": "172.21.0.2/16",
                "IPv6Address": ""
            }
        },
        "Options": {},
        "Labels": {
            "com.docker.compose.network": "default",
            "com.docker.compose.project": "kafka-local",
            "com.docker.compose.version": "1.21.1"
        }
    }
]
```
* Click "Save" button
![](assets/clusterinformation.png)
* Add topic
![](assets/Create_Topic.png)

