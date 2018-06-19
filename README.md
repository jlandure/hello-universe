# hello-universe

## Démos

- Show the code from handler.go
- Launch the code and go `localhost:8080`
```
go run main.go handler.go -http 0.0.0.0:8080
go build -o hello main.go handler.go && ./hello -http 0.0.0.0:8080
```
- Show the Dockerfile
- Create the image
```
docker build -t jlandure/hello-universe:0.0.1 .
```
- Launch using Docker and go to `http://192.168.99.100/`
```
docker run -p 8080:8080 jlandure/hello-universe:0.0.1 --http 0.0.0.0:8080
```
- Push the image
```
docker push jlandure/hello-universe:0.0.1
```
- get credentials for your K8S/GKE 
```
gcloud --project matinale-k8s container clusters get-credentials test-z
#service account from kubectl
kubectl get serviceAccounts/default -o yaml
kubectl describe secrets/default-token-q73rg
```

- Manual deployment
```
kubectl create -f replicasets/hello-universe.yaml
kubectl create -f services/hello-universe.yaml
kubectl describe services/hello-universe
kubectl scale --replicas=3 -f replicasets/hello-universe.yaml
```
- go to `http://130.211.75.174/` and check the hostname
- use the script `iwget` for infinite loop calls
- check logs with `kubectl logs hello-universe-2me3w` 

- Autoscaling
```
kubectl apply -f autoscaling/hpa.yaml
watch -n 1 kubectl get pods
watch -n 1 kubectl get hpa
ab -c 100 -n 1000 -t 10000 http://35.230.74.10/
```

## Example Usage

```
$ hello-universe -h
```
```
Usage of hello-universe:
  -api-server string
    	Kubernetes API server (default "http://127.0.0.1:8080")
  -cpu-limit string
    	Max CPU in milicores (default "100m")
  -cpu-request string
    	Min CPU in milicores (default "100m")
  -http string
    	HTTP service address (default "127.0.0.1:80")
  -kubernetes
    	Deploy to Kubernetes.
  -memory-limit string
    	Max memory in MB (default "64M")
  -memory-request string
    	Min memory in MB (default "64M")
  -replicas int
    	Number of replicas (default 1)
```

### Example

```
$ export HELLO_UNIVERSE_TOKEN="7346053dafaf4a24825790f4389704f5"
```

```
$ hello-universe -kubernetes -replicas 3  -cpu-limit 200m -memory-limit 500M
```

```
$ kubectl get replicasets
```

```
NAME             DESIRED   CURRENT   READY     AGE
hello-universe   3         3         0         15s
```

```
$ kubectl get pods
```
```
NAME                   READY     STATUS    RESTARTS   AGE
hello-universe-2genp   0/1       Pending   0          43s
hello-universe-gg1eu   0/1       Pending   0          43s
hello-universe-zgpuy   0/1       Pending   0          43s
```
