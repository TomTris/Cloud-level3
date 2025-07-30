#PostgresOperator
# original link https://access.crunchydata.com/documentation/postgres-operator/latest/quickstart
git clone https://github.com/CrunchyData/postgres-operator-examples.git
cd postgres-operator-examples
kubectl apply -k kustomize/install/namespace
kubectl apply --server-side -k kustomize/install/default
kubectl -n postgres-operator get pods --selector=postgres-operator.crunchydata.com/control-plane=postgres-operator --field-selector=status.phase=Running
# wait until pod is done, then
sleep 150 # wait until operator is ready

kubectl apply -k kustomize/postgres
kubectl get secret -n postgres-operator hippo-pguser-hippo -o yaml # <clusterName>-pguser-<userName>

SECRET_NAME=postgres-secret
NAMESPACE=pgo-go-backend

kubectl create secret generic $SECRET_NAME \
  -n $NAMESPACE \
  --from-literal=POSTGRES_USER=$(kubectl get secret hippo-pguser-hippo -n postgres-operator -o jsonpath="{.data.user}" | base64 -d) \
  --from-literal=POSTGRES_HOST=hippo-ha.postgres-operator.svc.cluster.local \
  --from-literal=POSTGRES_PASSWORD=$(kubectl get secret hippo-pguser-hippo -n postgres-operator -o jsonpath="{.data.password}" | base64 -d) \
  --from-literal=POSTGRES_NAME=$(kubectl get secret hippo-pguser-hippo -n postgres-operator -o jsonpath="{.data.dbname}" | base64 -d) \
  --from-literal=POSTGRES_PORT=$(kubectl get secret hippo-pguser-hippo -n postgres-operator -o jsonpath="{.data.port}" | base64 -d)

#Need to use  hippo-ha. it will route to the primary pod
  # --from-literal=POSTGRES_HOST=$(kubectl get secret hippo-pguser-hippo -n postgres-operator -o jsonpath="{.data.host}" | base64 -d) \

# host: hippo-primary.postgres-operator.svc
#   jdbc-uri: jdbc:postgresql://hippo-primary.postgres-operator.svc:5432/zoo?password=w6E4%3Bi%5B_%2AdnBplvvKv.ohQzU&user=hippo
#   password: w6E4;i[_*dnBplvvKv.ohQzU
#   port: 5432
#   uri: postgresql://hippo:w6E4;i%5B_%2AdnBplvvKv.ohQzU@hippo-primary.postgres-operator.svc:5432/zoo
#   user: hippo
#   verifier: SCRAM-SHA-256$4096:xngngdkq3Gfh82F7hZYi4Q==$uPY8UZpZuS9Us8dtPun5VloUBZmea2Iu1wCE9CXy8KE=:UycFhXLmteD6kS8Nz03u3EH/r0G6l/xRRYVZlwAR+DE=