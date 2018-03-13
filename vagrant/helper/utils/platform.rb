# Open the eigenclass for Vagrant::Util::Platform and override its OS and architecture detection
# for windows hosts. This code gives more reliable and correct results.
module Vagrant
    module Util
        class Platform
            class << self
                # *nix/posix and mac meta-platforms
                def posix?
                	platform =~ /darwin|bsd|linux|solaris/ || Process.respond_to?(:fork) ? true : false
                end
                def mac?
                	platform.include?('darwin')
                end

                # Windows override
                def windows?
                    platform =~ /mswin|mingw|cygwin/ ? true : false
                end

                # bit64 override
                def bit64?
                    if windows?
                        require 'win32/registry'
                        Win32::Registry::HKEY_LOCAL_MACHINE.open(
                            'HARDWARE\DESCRIPTION\System\CentralProcessor\0'
                        )['Platform ID'] == 32 ? false : true
                    else
                    	[' '].pack('P').size * 8 == 64 ? true : false
                    end
                end
            end
        end
    end
end