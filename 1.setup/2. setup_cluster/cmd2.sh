# Do this in vm-worker(s)
# Take K3S_TOKEN from vm-master

#Write token to here
K3S_TOKEN_MASTER=
vm-master-ip=

# change hostname
sudo sed -i '1s/.*/127.0.0.1 localhost vm-worker-1/' /etc/hosts
# For DNS Server
echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf
curl -sfL https://get.k3s.io | K3S_URL=https://$(vm-master-ip):6443 K3S_TOKEN=K3S_TOKEN_MASTER sh -



# Commands to check
# sudo systemctl status k3s-agent

# In vm-master
# sudo kubectl get nodes
# sudo kubectl get node vm-1 --show-labels


#--- ----setup all-worker-cluster--------
curl -sfL https://get.k3s.io | K3S_TOKEN=K1086b211d8dc7b6a1aa18490960337d4115d86a45640c636a248df423578478cbc::server:b2cfde3314560a19cb53a6ff44645c5b sh -s - server --server https://168.119.243.127:6443 --tls-san 168.119.243.127 --tls-san 172.24.4.85


sudo sed -i '1s/.*/127.0.0.1 localhost vm-worker-2/' /etc/hosts
# For DNS Server
echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf
curl -sfL https://get.k3s.io | K3S_TOKEN=K10e2ac55a788d5128368aecf84b80a9e28a930bdbf162b76b4ccefbbde3fe093c4::server:6b21dd3cafd2aa06930df2da1892839b sh -s - server --server https://91.99.116.140:6443 --tls-san 91.99.116.140