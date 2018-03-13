# Advanced Object Methods
class Object
  def blank?
    respond_to?(:empty?) ? empty? : !self
  end

  def present?
    !blank?
  end

  def presence
    self if present?
  end
end

class NilClass
  def blank?
    true
  end
end

class FalseClass
  def blank?
    true
  end
end

class TrueClass
  def blank?
    false
  end
end

class Array
  alias_method :blank?, :empty?
end

class Hash
  alias_method :blank?, :empty?
end

class String
  def blank?
    if defined?(Encoding) && "".respond_to?(:encode)
      self !~ /[^[:space:]]/
    else
      self !~ %r![^\s#{[0x3000].pack("U")}]!
    end
  end
end

class Numeric
  def blank?
    false
  end
end