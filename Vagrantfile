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

		vmCfg = ensure_disk(vmCfg, guest_disk_path, 'guests.vmdk')

		vmCfg = configureFreeBSDDevProvisioners(vmCfg)

		vmCfg = configureSyncedDir(vmCfg, '.',
			'/opt/gopath/src/github.com/joyent/freebsd-vpc')

		vmCfg = configureMachineSize(vmCfg, 4, 8192)
	end
end

def configureMachineSize(vmCfg, vcpuCount, memSize)
	["vmware_fusion", "vmware_workstation"].each do |p|
		vmCfg.vm.provider p do |v|
			v.vmx["memsize"] = memSize
			v.vmx["numvcpus"] = vcpuCount
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
	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-zpool.sh',
		privileged: true

	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-dev-packages.sh',
		privileged: true
	
	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-db-configure.sh',
		privileged: true
	
	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-priv-bhyve.sh',
		privileged: true

	vmCfg.vm.provision "shell",
		path: './vagrant/scripts/vagrant-freebsd-unpriv-dev-make.sh',
		privileged: false

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
