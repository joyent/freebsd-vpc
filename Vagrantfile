# -*- mode: ruby -*-
# vi: set ft=ruby :

freebsd_box = 'jen20/FreeBSD-12.0-CURRENT-VPC'
guest_disk_path = "#{File.dirname(__FILE__)}/guest_disks"

Vagrant.configure("2") do |config|
	config.ssh.extra_args = ["-e", "%"]

	config.vm.define "cn1", autostart: true, primary: true do |vmCfg|
		vmCfg.vm.box = freebsd_box
		vmCfg.vm.hostname = "freebsd-cn1"
		vmCfg = configureFreeBSDProvisioners(vmCfg)
		vmCfg = ensure_disk(vmCfg, guest_disk_path, 'cn1_guests.vmdk')

		vmCfg.vm.network "private_network", ip: "172.27.10.10"

		vmCfg.vm.provider "vmware_fusion" do |v|
			v.vmx["memsize"] = "4096"
			v.vmx["numvcpus"] = "2"
			v.vmx["ethernet1.virtualDev"] = "vmxnet3"
		end
	end
	
	config.vm.define "cn2", autostart: true do |vmCfg|
		vmCfg.vm.box = freebsd_box
		vmCfg.vm.hostname = "freebsd-cn2"
		vmCfg = configureFreeBSDProvisioners(vmCfg)
		vmCfg = ensure_disk(vmCfg, guest_disk_path, 'cn2_guests.vmdk')

		vmCfg.vm.network "private_network", ip: "172.27.10.11"

		vmCfg.vm.provider "vmware_fusion" do |v|
			v.vmx["memsize"] = "4096"
			v.vmx["numvcpus"] = "2"
			v.vmx["ethernet1.virtualDev"] = "vmxnet3"
		end
	end
end

def configureFreeBSDProvisioners(vmCfg)
	vmCfg.vm.provision "shell",
		path: './scripts/vagrant-freebsd-priv-zpool.sh',
		privileged: true

	vmCfg.vm.provision "shell",
		path: './scripts/vagrant-freebsd-priv-packages.sh',
		privileged: true

	vmCfg.vm.provision "shell",
		path: './scripts/vagrant-freebsd-priv-avahi.sh',
		privileged: true

	vmCfg.vm.provision "shell",
		path: './scripts/vagrant-freebsd-priv-bhyve.sh',
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
