# -*- mode: ruby -*-
# vi: set ft=ruby :

freebsd_box = 'joyent/FreeBSD-12.0-CURRENT-VPC'
guest_disk_path = "#{File.dirname(__FILE__)}/vagrant/guest_disks"

require './vagrant/helper/core'
require './vagrant/helper/utils'

Vagrant.configure("2") do |config|
	config.ssh.extra_args = ["-e", "%"]

	config.vm.define "compile", autostart: true, primary: true do |vmCfg|
		vmCfg.vm.box = freebsd_box
		vmCfg.vm.hostname = "freebsd-compile"
		vmCfg = configureFreeBSDDevProvisioners(vmCfg)

		vmCfg = configureSyncedDir(vmCfg, '.',
			'/opt/gopath/src/github.com/joyent/freebsd-vpc')

		vmCfg = addPrivateNICOptions(vmCfg, "172.27.10.5")
		vmCfg = configureMachineSize(vmCfg, 2, 1024)
	end

	3.times do |n|
		hostname = "crdb%d" % [n + 1]
		ip = "172.27.10.%d" % [n + 11]

		config.vm.define hostname, autostart: false do |vmCfg|
			vmCfg.vm.box = freebsd_box
			vmCfg.vm.hostname = hostname
			vmCfg = configureFreeBSDDBProvisioners(vmCfg, hostname, ip)

			vmCfg = addPrivateNICOptions(vmCfg, ip)
			vmCfg = configureMachineSize(vmCfg, 2, 1024)
		end
	end

	config.vm.define "cn1", autostart: false do |vmCfg|
		vmCfg.vm.box = freebsd_box
		vmCfg.vm.hostname = "freebsd-cn1"
		vmCfg = configureFreeBSDProvisioners(vmCfg)
		vmCfg = ensure_disk(vmCfg, guest_disk_path, 'cn1_guests.vmdk')

		vmCfg = addPrivateNICOptions(vmCfg, "172.27.10.20")
		vmCfg = configureMachineSize(vmCfg, 2, 4096)
	end

	config.vm.define "cn2", autostart: false do |vmCfg|
		vmCfg.vm.box = freebsd_box
		vmCfg.vm.hostname = "freebsd-cn2"
		vmCfg = configureFreeBSDProvisioners(vmCfg)
		vmCfg = ensure_disk(vmCfg, guest_disk_path, 'cn2_guests.vmdk')
		
		vmCfg = addPrivateNICOptions(vmCfg, "172.27.10.21")
		vmCfg = configureMachineSize(vmCfg, 2, 4096)
	end
end

def addPrivateNICOptions(vmCfg, ip)
	vmCfg.vm.network "private_network", ip: ip

	["vmware_fusion", "vmware_workstation"].each do |p|
		vmCfg.vm.provider p do |v|
			v.vmx["ethernet1.virtualdev"] = "vmxnet3"
			v.vmx["ethernet1.pcislotnumber"] = "192"
		end
	end

	return vmCfg
end

def configureMachineSize(vmCfg, vcpuCount, memSize)
	["vmware_fusion", "vmware_workstation"].each do |p|
		vmCfg.vm.provider p do |v|
			v.vmx["memsize"] = "1024"
			v.vmx["numvcpus"] = "2"
		end
	end

	return vmCfg
end

def configureSyncedDir(vmCfg, hostSource, guestTarget)
	if Vagrant::Util::Platform::windows?
		vmCfg.vm.synced_folder hostSource,
			guestTarget,
			type: "nfs",
			mount_options: ['nfsv3', 'mntudp', 'vers=3', 
				'udp', 'noatime']
	else
		vmCfg.vm.synced_folder hostSource,
			guestTarget,
			type: "nfs",
			bsd__nfs_options: ['noatime']
	end

	return vmCfg
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

	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-unpriv-dev-migrate.sh',
		privileged: false

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
	completePath = File.join(dirname, filename)
	if Vagrant::Util::Platform::mac?
		vdiskmanager = '/Applications/VMware Fusion.app/Contents/Library/vmware-vdiskmanager'
	elsif Vagrant::Util::Platform::windows?
		vdiskmanager = "C:\\Program Files (x86)\\VMWare\\VMWare Workstation\\vmware-vdiskmanager.exe"
	end

	unless Dir.exists?(dirname)
		Dir.mkdir dirname
	end

	unless File.exists?(completePath)
		system("cd \"#{dirname}\" && \"#{vdiskmanager}\" -c -s 30GB -a lsilogic -t 1 \"#{filename}\"")
	end

	["vmware_fusion", "vmware_workstation"].each do |p|
		vmCfg.vm.provider p do |v|
			v.vmx["scsi0:1.filename"] = File.expand_path(completePath)
			v.vmx["scsi0:1.present"] = 'TRUE'
			v.vmx["scsi0:1.redo"] = ''
		end
	end

	return vmCfg
end
