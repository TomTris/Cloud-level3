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
