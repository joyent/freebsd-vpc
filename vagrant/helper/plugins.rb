#
# Plugin loader - shamelessly copied and modified from vagrant (lib/vagrant.rb)
load_plugin = lambda do |dir|
	next false if !dir.directory?
	plugin_file = dir.join('plugin.rb')
	if plugin_file.file?
		load(plugin_file)
		next true
	end
end

#
# Automatically load all plugins in the plugins directory
Pathname.new(File.expand_path('../plugins', __FILE__)).children(true).each do |dir|
	next if !dir.directory?
	next if load_plugin.call(dir)
	dir.children(true).each(&load_plugin)
end