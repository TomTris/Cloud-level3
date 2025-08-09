./stack.sh
ssh-keygen -t rsa  -b 4096



source ~/devstack/openrc admin admin
wget https://cloud-images.ubuntu.com/jammy/20250619/jammy-server-cloudimg-amd64.img
openstack image create jammy_cloudimg --file jammy-server-cloudimg-amd64.img  --disk-format qcow2 --container-format bare
openstack flavor create costum --ram 16384 --disk 50 --vcpus 8 --public
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
# openstack security group rule create --proto tcp --dst-port 2379:2380 k3s_rule
# openstack security group rule create --proto tcp --dst-port 10250 k3s_rule
openstack security group rule create --proto icmp k3s_rule


for port in {30000..30010}; do
  echo PREROUTING add $port
  sudo iptables -t nat -A PREROUTING -p tcp --dport $port -j DNAT --to-destination 172.24.4.85:$port
done

for port in {30010..30020}; do
  echo $port
  sudo iptables -t nat -D PREROUTING 6
done