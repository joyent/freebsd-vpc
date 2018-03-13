# Advanced Hash Features
class Hash
	#
	# Special method to convert keys, valid values for op would be :to_s or :to_sym
	def deep_convert_keys(op)
		return self if !op.is_a?(Symbol)
		inject({}) do |result, (key, value)|
			result[(key.send(op) rescue key) || key] = value.is_a?(Hash) ? value.deep_convert_keys(op) : value
			result
		end
	end

	def deep_symbolize_keys
		deep_convert_keys(:to_sym)
	end unless Hash.method_defined?(:deep_symbolize_keys)

	def deep_stringify_keys
		deep_convert_keys(:to_s)
	end unless Hash.method_defined?(:deep_stringify_keys)

	#
	# Deep merging of hashes (also merges arrays)
	def deep_merge(other)
		dup.deep_merge!(other)
	end unless Hash.method_defined?(:deep_merge)

	def deep_merge!(other)
		other.each_pair do |k, v|
			t = self[k]
			self[k] = t.is_a?(Hash) && v.is_a?(Hash) ? t.deep_merge(v) : t.is_a?(Array) && v.is_a?(Array) ? t | v : v
		end
		self
	end unless Hash.method_defined?(:deep_merge!)
end