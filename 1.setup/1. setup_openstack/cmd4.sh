./stack.sh
ssh-keygen -t rsa  -b 4096



source ~/devstack/openrc admin admin
wget https://cloud-images.ubuntu.com/jammy/20250619/jammy-server-cloudimg-amd64.img
openstack image create jammy_cloudimg --file jammy-server-cloudimg-amd64.img  --disk-format qcow2 --container-format bareopenstack flavor create costum --ram 16384 --disk 50 --vcpus 8 --public
openstack keypair create --public-key ~/.ssh/id_rsa.pub server-key
openstack server create --flavor costum --image jammy_cloudimg --network private --key-name server-key vm-master
openstack server create --flavor costum --image jammy_cloudimg --network private --key-name server-key vm-worker-1
floating_ip=$(openstack floating ip create public -f value -c floating_ip_address)
openstack server add floating ip vm-master $floating_ip
floating_ip=$(openstack floating ip create public -f value -c floating_ip_address)
openstack server add floating ip vm-worker-1 $floating_ip
openstack security group create k3s_rule --description "for my cluster"
openstack server add security group vm-master k3s_rule
openstack server add security group vm-worker-1 k3s_rule
openstack security group rule create --proto tcp --dst-port 22 k3s_rule
openstack security group rule create --proto tcp --dst-port 30000:32767 k3s_rule
openstack security group rule create --proto tcp --dst-port 6443:6444 k3s_rule
openstack security group rule create --proto icmp k3s_rule

openstack server show vm-master -f json | jq -r '.addresses' | grep -oE '[0-9]+\.[0-9]+\.[0-9]+\
.[0-9]+' | while read ip; do   if openstack floating ip list --floating-ip-address $ip -f value -c "Floating IP Addres
s" | grep -q $ip; then     echo $ip;   fi; done

for port in {30000..32767}; do
  echo PREROUTING add $port
  sudo iptables -t nat -A PREROUTING -p tcp --dport $port -j DNAT --to-destination 172.24.4.32:$port
done

# for port in {30000..32767}; do
#   echo $port
#   sudo iptables -t nat -D PREROUTING 1
# done