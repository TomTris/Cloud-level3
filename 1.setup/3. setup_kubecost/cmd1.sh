# helm upgrade --install kubecost \
#  --repo https://kubecost.github.io/cost-analyzer/ cost-analyzer \
# --namespace kubecost --create-namespace \
# --kubeconfig /etc/rancher/k3s/k3s.yaml
kubectl create namespace kubecost
kubectl apply -f https://raw.githubusercontent.com/kubecost/cost-analyzer-helm-chart/v2.8/kubecost.yaml
sleep 300 # wait 5 min for kubecost to be ready
kubectl patch svc kubecost-cost-analyzer -n kubecost --type='merge' -p '{
  "spec": {
    "type": "NodePort",
    "ports": [
      {
        "name": "tcp-model",
        "port": 9003,
        "protocol": "TCP",
        "targetPort": 9003
      },
      {
        "name": "tcp-frontend",
        "port": 9090,
        "protocol": "TCP",
        "targetPort": 9090,
        "nodePort": 30003
      }
    ]
  }
}'
# To Edit:
# kubectl edit service kubecost-cost-analyzer -n kubecost
# To change type to NodePort
# kubectl patch service kubecost-cost-analyzer -n kubecost -p '{"spec": {"type": "NodePort"}}'