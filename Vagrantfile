Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu-18.04-amd64"

  config.vm.hostname = "ldap.example.com"

  config.vm.network "private_network", ip: "192.168.33.253", libvirt__forward_mode: "route", libvirt__dhcp_enabled: false

  config.vm.provider "libvirt" do |lv, config|
    lv.memory = 2048
    lv.cpus = 2
    lv.cpu_mode = "host-passthrough"
    lv.keymap = "pt"
    config.vm.synced_folder ".", "/vagrant", type: "nfs"
  end

  config.vm.provider "virtualbox" do |vb|
    vb.linked_clone = true
    vb.memory = 1024
    vb.cpus = 2
  end

  config.vm.provision "shell", path: "provision.sh"
end
