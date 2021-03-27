# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.box = "bento/ubuntu-18.04"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # NOTE: This will enable public access to the opened port
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine and only allow access
  # via 127.0.0.1 to disable public access
  config.vm.network "forwarded_port", guest: 8081, host: 8081, host_ip: "127.0.0.1"
  config.vm.network "forwarded_port", guest: 8082, host: 8082, host_ip: "127.0.0.1"
  config.vm.network "forwarded_port", guest: 8083, host: 8083, host_ip: "127.0.0.1"

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  # config.vm.provider "virtualbox" do |vb|
  #   # Display the VirtualBox GUI when booting the machine
  #   vb.gui = true
  #
  #   # Customize the amount of memory on the VM:
  #   vb.memory = "1024"
  # end
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  # Enable provisioning with a shell script. Additional provisioners such as
  # Ansible, Chef, Docker, Puppet and Salt are also available. Please see the
  # documentation for more information about their specific syntax and use.
  config.vm.provision "shell", inline: <<-SHELL
    cd /tmp
    go_url=https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz
    go_ark=$(basename $go_url)
    helm_url=https://get.helm.sh/helm-v3.5.3-linux-amd64.tar.gz

    export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

    sudo add-apt-repository ppa:dqlite/stable -y
    sudo apt-get update -y
    sudo apt-get install -y libdqlite-dev build-essential docker docker-compose
    sudo usermod -aG docker vagrant
    
    curl -sL $go_url -o $go_ark 
    ls go || tar -xvf $go_ark
    ls /usr/local/go &> /dev/null || sudo mv go /usr/local
    
    if [ -z "$GOROOT" ];then
      echo export GOROOT=/usr/local/go >> /home/vagrant/.bashrc
      echo export GOPATH=/home/vagrant/go >> /home/vagrant/.bashrc
      echo export PATH=$GOPATH/bin:$GOROOT/bin:$PATH >> /home/vagrant/.bashrc
    else
      echo "\$env vars setted up"
    fi

    which k3s || curl -sfL https://get.k3s.io | sh -

    which helm || (curl -sL $helm_url -o $(basename $helm_url) && \
      tar -xvf helm-v3.5.3-linux-amd64.tar.gz && \
      sudo cp ./linux-amd64/helm /usr/local/bin/ && \
      sudo chmod 775 /usr/local/bin/helm)

    helm list &> /dev/null || sudo chown vagrant:root /etc/rancher/k3s/k3s.yaml
    kubectl get pods -l app.kubernetes.io/name=ingress-nginx | grep nginx &> /dev/null

    if [ "$?" -ne 0 ];then
      helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
      helm repo update
      helm install ingress-nginx ingress-nginx/ingress-nginx
    fi

  SHELL
end
