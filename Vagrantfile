# -*- mode: ruby -*-
# vi: set ft=ruby :

freebsd_box = 'jen20/FreeBSD-12.0-CURRENT-VPC'
guest_disk_path = "#{File.dirname(__FILE__)}/vagrant/guest_disks"

Vagrant.configure("2") do |config|
	config.ssh.extra_args = ["-e", "%"]
	
	config.vm.define "compile", autostart: true, primary: true do |vmCfg|
		vmCfg.vm.box = freebsd_box
		vmCfg.vm.hostname = "freebsd-compile"
		vmCfg = configureFreeBSDDevProvisioners(vmCfg)

		vmCfg.vm.synced_folder '.',
			'/opt/gopath/src/github.com/sean-/vpc',
			type: "nfs",
			bsd__nfs_options: ['noatime']

		vmCfg.vm.network "private_network", ip: "172.27.10.5"

		vmCfg.vm.provider "vmware_fusion" do |v|
			v.vmx["memsize"] = "1024"
			v.vmx["numvcpus"] = "2"
			v.vmx["ethernet1.virtualDev"] = "vmxnet3"
		end
	end

	3.times do |n|
		hostname = "crdb%d" % [n + 1]
		ip = "172.27.10.%d" % [n + 11]

		config.vm.define hostname, autostart: false do | vmCfg|
			vmCfg.vm.box = freebsd_box
			vmCfg.vm.hostname = hostname
			vmCfg = configureFreeBSDDBProvisioners(vmCfg, hostname, ip)
		
			vmCfg.vm.network "private_network", ip: ip
		
			vmCfg.vm.provider "vmware_fusion" do |v|
				v.vmx["memsize"] = "1024"
				v.vmx["numvcpus"] = "2"
				v.vmx["ethernet1.virtualDev"] = "vmxnet3"
			end
		end
	end

	config.vm.define "cn1", autostart: false do |vmCfg|
		vmCfg.vm.box = freebsd_box
		vmCfg.vm.hostname = "freebsd-cn1"
		vmCfg = configureFreeBSDProvisioners(vmCfg)
		vmCfg = ensure_disk(vmCfg, guest_disk_path, 'cn1_guests.vmdk')

		vmCfg.vm.network "private_network", ip: "172.27.10.20"

		vmCfg.vm.provider "vmware_fusion" do |v|
			v.vmx["memsize"] = "4096"
			v.vmx["numvcpus"] = "2"
			v.vmx["ethernet1.virtualDev"] = "vmxnet3"
		end
	end
	
	config.vm.define "cn2", autostart: false do |vmCfg|
		vmCfg.vm.box = freebsd_box
		vmCfg.vm.hostname = "freebsd-cn2"
		vmCfg = configureFreeBSDProvisioners(vmCfg)
		vmCfg = ensure_disk(vmCfg, guest_disk_path, 'cn2_guests.vmdk')

		vmCfg.vm.network "private_network", ip: "172.27.10.21"

		vmCfg.vm.provider "vmware_fusion" do |v|
			v.vmx["memsize"] = "4096"
			v.vmx["numvcpus"] = "2"
			v.vmx["ethernet1.virtualDev"] = "vmxnet3"
		end
	end
end

def configureFreeBSDDevProvisioners(vmCfg)
	vmCfg.vm.provision "file",
		source: './vagrant/certs/ca/ca.crt',
		destination: "/home/vagrant/.cockroach-certs/ca.crt"

	vmCfg.vm.provision "file",
		source: "./vagrant/certs/client/client.root.crt",
		destination: "/home/vagrant/.cockroach-certs/client.root.crt"

	vmCfg.vm.provision "file",
		source: "./vagrant/certs/client/client.root.key",
		destination: "/home/vagrant/.cockroach-certs/client.root.key"

	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-dev-packages.sh',
		privileged: true

	return vmCfg
end

def configureFreeBSDDBProvisioners(vmCfg, hostname, ip)
	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-db-packages.sh',
		privileged: true

	vmCfg.vm.provision "file",
		source: './vagrant/certs/ca/ca.crt',
		destination: "/home/vagrant/.cockroach-certs/ca.crt"
	
	vmCfg.vm.provision "file",
		source: "./vagrant/certs/client/client.root.crt",
		destination: "/home/vagrant/.cockroach-certs/client.root.crt"
	
	vmCfg.vm.provision "file",
		source: "./vagrant/certs/client/client.root.key",
		destination: "/home/vagrant/.cockroach-certs/client.root.key"
	
	vmCfg.vm.provision "file",
		source: './vagrant/certs/ca/ca.crt',
		destination: "/secrets/crdb/ca.crt"
	
	vmCfg.vm.provision "file",
		source: "./vagrant/certs/#{hostname}/node.crt",
		destination: "/secrets/crdb/node.crt"
	
	vmCfg.vm.provision "file",
		source: "./vagrant/certs/#{hostname}/node.key",
		destination: "/secrets/crdb/node.key"

	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-db-configure.sh',
		privileged: true

	if hostname == "crdb3"
		vmCfg.vm.provision "shell",
			path: './vagrant/scripts/vagrant-freebsd-unpriv-db-init.sh',
			privileged: false
	end

	return vmCfg
end

def configureFreeBSDProvisioners(vmCfg)
	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-zpool.sh',
		privileged: true

	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-packages.sh',
		privileged: true

	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-avahi.sh',
		privileged: true

	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-bhyve.sh',
		privileged: true

	return vmCfg
end

def ensure_disk(vmCfg, dirname, filename)
	vdiskmanager = '/Applications/VMware\ Fusion.app/Contents/Library/vmware-vdiskmanager'

	unless Dir.exists?(dirname)
		Dir.mkdir dirname
	end

	completePath = File.join(dirname, filename)

	unless File.exists?(completePath)
		`#{vdiskmanager} -c -s 30GB -a lsilogic -t 1 #{completePath}`
	end

	vmCfg.vm.provider "vmware_fusion" do |v|
		v.vmx["scsi0:1.filename"] = File.expand_path(completePath)
		v.vmx["scsi0:1.present"] = 'TRUE'
		v.vmx["scsi0:1.redo"] = ''
	end

	return vmCfg
end
