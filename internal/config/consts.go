package config

const (
	DefaultManDir = "./docs/man"
	ManSect       = 8

	DefaultMarkdownDir       = "./docs/md"
	DefaultMarkdownURLPrefix = "/command"

	KeyDocManDir            = "doc.mandir"
	KeyDocMarkdownDir       = "doc.markdown-dir"
	KeyDocMarkdownURLPrefix = "doc.markdown-url-prefix"

	KeyLogFormat    = "log.format"
	KeyLogLevel     = "log.level"
	KeyLogStats     = "log.stats"
	KeyLogTermColor = "log.use-color"

	KeyShellAutoCompBashDir = "shell.autocomplete.bash-dir"

	KeySWCreateID  = "switch.create.id"
	KeySWCreateMAC = "switch.create.mac"
	KeySWCreateVNI = "switch.create.vni"
	KeySWDestroyID = "switch.destroy.id"

	KeyUsePager = "general.use-pager"
	KeyUseUTC   = "general.utc"

	KeyVMNICCreateID  = "vmnic.create.id"
	KeyVMNICCreateMAC = "vmnic.create.mac"
	KeyVMNICDestroyID = "vmnic.destroy.id"
)
