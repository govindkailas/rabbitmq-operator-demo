# RabbitMQ Operator Demo
We are going to try a simple RabbitMQ cluster created using RabbitMQ operator. Once the cluster is up, will use the go producer to push some messages.
We will later consume those messages using the consumers.


# Pre-req
1. Expect that you already have a working kubernetes cluster with the needed permissions to run everything.
2. Go must be installed

## Install RabbitMQ Cluster Operator
```
kubectl apply -f "https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml"

```
For more details refer the official documentation here https://www.rabbitmq.com/kubernetes/operator/quickstart-operator.html

1. Create a minimal RabbitMQ Cluster
```
kubectl apply -f rmq.yaml

 kubectl get po # run this to check the status of the pod
```
Wait for a few minutes for the pod to be up. 

Set the below env variables
```
username="$(kubectl get secret rmq-demo-default-user -o jsonpath='{.data.username}' | base64 --decode)"
password="$(kubectl get secret rmq-demo-default-user -o jsonpath='{.data.password}' | base64 --decode)"
rmqLBIP=$(k get svc rmq-demo -o jsonpath='{.status.loadBalancer.ingress[].ip}')
## rmqLBIP=$(k get svc rmq-demo -o jsonpath='{.status.clusterIP') ## if clusterIP user this
RMQ_SERVER_URL="amqp://${username}:${password}@${rmqLBIP}:5672/" 
export RMQ_SERVER_URL=$RMQ_SERVER_URL

```

We are exposing the service as a LoadBalancer, if its not suppored on your cluster, you can always leave it as default(ClusterIP)  portforward to local.

To verify the connection, you can connect to the API or the Web UI
```
curl -u${username}:${password} ${rmqLBIP}:15672/api/overview | jq
```



2. Run the Producer app:
```
cd rabbitmq/
go run .
```
It will wait for user input from console. Press `enter` when you want to send a message and print `exit` for graceful shutting down.

3. Run the Consumer app:
```
cd rabbitmq/consumer1
go run .

cd rabbitmq/consumer2
go run .
```
It will connect to the RabbitMQ queue and print all the messages awaiting; each consumer also imitates some work going on by sleeping for 10 seconds after receiving each message. The RabbitMQ acknowledgement is done manually by marking `d.Ack(false)` after every message processing.
Consumer1 also receives no more than 2 messages at once, Consumer2 - no more than 1


4. Let us do a perf test of the RabbitMQ 
```
kubectl run perf-test --image=pivotalrabbitmq/perf-test -- --uri amqp://${username}:${password}@${rmqLBIP}
```

To check the logs of perf test 
```
kubectl logs -f pod/perf-test
```


Referances and credits to - https://github.com/koterin/broker