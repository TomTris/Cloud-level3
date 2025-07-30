# Do this in vm-master
# change hostname
sudo sed -i '1s/.*/127.0.0.1 localhost vm-master/' /etc/hosts

# For DNS Server
echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf

# for creating master and nodes
curl -sfL https://get.k3s.io | sh -
sudo cat /var/lib/rancher/k3s/server/node-token
