# Check which ports the server is listening to
# redirect requests (not ssh) from port 2222 to 172.24.4.94:22
# with the previous command, IP-Address is changed to own private IP to for assigning purpose. To change it back to original IP address to answer, we need this command. 
sudo ss -tuln | grep LISTEN
sudo iptables -t nat -A PREROUTING -p tcp --dport $port -j DNAT --to-destination <>:30485
sudo iptables -t nat -A POSTROUTING -j MASQUERADE
sudo iptables -t nat -L -v -n --line-numbers
sudo iptables -t nat -D PREROUTING 1


