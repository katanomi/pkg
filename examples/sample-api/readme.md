# Try It

1. build and deploy `ko apply -L -f pod.yaml`
2. execute script to request api of `/v1/hello`

```bash
PORT=`kubectl get svc sample-api -o go-template --template="{{ (index .spec.ports 0).nodePort }}"`
SECRET=`kubectl get sa sample-api -o go-template --template="{{ (index .secrets 0).name }}"`
TOKEN=`kubectl get secret ${SECRET} --template "{{.data.token}}" | base64 -d`

curl -H "Authorization: Bearer ${TOKEN}" -v "http://localhost:${PORT}/v1/hello"
```
3. ensure that you will get all namespaces list count like this:

```json
{
 "client": 13,
 "dynamicClient": 13
}
```
