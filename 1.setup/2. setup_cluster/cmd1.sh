# ----setup 1-master-cluster--------
# Do this in vm-master
# change hostname
sudo sed -i '1s/.*/127.0.0.1 localhost vm-master/' /etc/hosts

# For DNS Server
echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf

# for creating master and nodes
curl -sfL https://get.k3s.io | sh -
sudo cat /var/lib/rancher/k3s/server/node-token


# ----setup all-master-cluster--------

curl -sfL https://get.k3s.io | sh -s - server --cluster-init --tls-san <Floating_ip_First_Node>
sudo cat /var/lib/rancher/k3s/server/node-token

K1086b211d8dc7b6a1aa18490960337d4115d86a45640c636a248df423578478cbc::server:b2cfde3314560a19cb53a6ff44645c5b