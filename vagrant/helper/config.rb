require 'yaml'

module UserConfig
	def self.load(name)
		file = name.to_s << (name.is_a?(Symbol) ? '.yml' : '')

		# Get the resolved path
		if Pathname.new(file).relative?
			base = File.dirname(caller[0].split(':')[0])
			['', 'config'].each do |search|
				if File.exists?(f = File.join(base, search, file))
					file = f
					break
				end
			end
		end

		# Return the data or an empty hash if the file doesn't exist
		File.exists?(file) ?
			YAML.load(File.open(file)) :
			{}
	end
end