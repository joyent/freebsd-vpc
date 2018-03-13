vagrant-helper
==============

Common Vagrant Helpers for User Config, OS Detection, Etc.

## What is vagrant-helper?

At DDM, we're starting to use Vagrant a lot. Because it's just Ruby, we've found
that we are adding a little bit of programming to our ```Vagrantfile``` to make
things easier to share among our team. So we've made this repo intended to be a 
submodule in our vagrant projects.

## Installation

### Git

In your git project, add this as a submodule:

```
cd /path/to/vagrant/root/
git submodule add https://github.com/deseretdigital/vagrant-helper.git helper
```

### Other

You can [download the latest vagrant-helper](https://github.com/deseretdigital/vagrant-./helper/archive/master.zip), 
create a directory in your vagrant project called "helper" and then move the
unziped files there.

## Usage

### Setup

At the top of your ```vagrantfile``` you can include the helper libraries like
so:

```ruby
require './helper/core'		# Required
require './helper/utils'		# If you want to use the Utils helpers
require './helper/config'	# If you want to use the Config helpers
```

## User Configs

Sometimes it is helpful to allow some customization of a Vagrant project for things
like the location of files you would like to mount. Many times a project will
require several mounts, and these paths on the host machine can vary from
developer to developer.

You can create in your project the ```config``` directory and add a ```prefs.yml```
file. It can then look something like this:

```yml
vm:
  provider: virtualbox
  box:
    name: precise
  network: 10.13.37.10
  forward:
    http: 8080
    ssh:  2222
```

Then in your Vagrantfile you can load the User Preferences. You can even check to make
sure they were loaded.

```ruby
require './helper/core'		# Required
require './helper/config'	# If you want to use the Config helpers

Vagrant.configure("2") do |config|
    # Load the user preferences in config/prefs.yml
    if (prefs = UserConfig.load(:prefs)).empty?
        abort("Error: no user preferences were loaded, make sure you have created 'config/prefs.yml'")
    end

  	# [ ... continue Vagrantfile settings ... ]

end
```

Finally, to use a config setting you can just do this:

```ruby
config.vm.network :hostonly, ip: prefs['vm']['network']
```

## Platform Detection

We have developers on Linux, Mac, and Windows (32 and 64 bit). Sometimes there are 
certain settings that are more performant than others for a given Provider and OS.

Our Utils helper allows you to detect certain platform settings so you can only
apply certain settings to a particular environment. 

### 64bit vs 32bit

You can't have a 64bit guest on a 32bit host, so you can decide which basebox to use
like so:

```ruby
arch = Vagrant::Util::Platform::bit64? ? 64 : 32
config.vm.box = "precise#{arch}"
```

### Is Posix?

If you're on VirtualBox & a Posix based machine (i.e. Mac, Linux) with a large amount of files to share with the 
guest machine, you likely will run into performance problems. So its recommended
to use NFS. However, other providers, or other OSes, do not have this problem
so it would be nice not to require them to have NFS on. So here is a way to
do it with the Platform utils:


**Note, right now there I can't find a programatic way to detect which provider
is being used, so it's best to have that in your config file. See above.**

```ruby
provider = prefs['vm']['provider'].to_sym

config.vm.synced_folder(
    '/path/to/host/www', 	# Host
    "/var/www",	 			# Guest

    # NFS on *nix-based platforms to resolve performance issues
    :nfs => (Vagrant::Util::Platform::posix? and provider == :virtualbox)
)
```

### Other Detections

```ruby
Vagrant::Util::Platform::mac?
Vagrant::Util::Platform::windows?
```