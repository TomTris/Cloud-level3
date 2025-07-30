Kubectl commands


Sudo kubectl get all

sudo kubectl get namespace
sudo kubectl get pods --all-namespaces
sudo kubectl delete namespaces db-platform
sudo kubectl create namespace db-platform

kubectl get all -n <namespace>

kubectl exec  -it simple-paas-b7bcdb859-z8txj  -n simple-paas  -- sh
kubectl get serviceaccount -n simple-paas
kubectl get pods -o wide —all-namespaces

kubectl top nodes
kubectl get nodes

###
Change current namespace
###

kubectl config current-context
kubectl config set-context --current --namespace=kubecost
# Then verify
kubectl config view --minify | grep namespace



kubectl get serviceaccount -n default --no-headers | awk '$1 != "default" {print $1}' | xargs -r -n1 kubectl delete serviceaccount -n default

#on which node
kubectl get pod -n postgres-operator -o wide

# postgres-operator - replicas
# get primary and replica pods
kubectl get pods -n postgres-operator -l postgres-operator.crunchydata.com/role=master
kubectl get pods -n postgres-operator -L postgres-operator.crunchydata.com/role

# get delay
kubectl exec -it hippo-instance1-4dnv-0 -n postgres-operator -- bash
psql -U postgres -d postgres
SELECT now() - pg_last_xact_replay_timestamp() AS replication_lag;
