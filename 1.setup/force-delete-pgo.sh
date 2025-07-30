kubectl get namespace postgres-operator -o json > ns.json
nano ns.json
"finalizers": [
  "kubernetes"
]
"spec": {
  "finalizers": []
}
kubectl replace --raw "/api/v1/namespaces/postgres-operator/finalize" -f ./ns.json
