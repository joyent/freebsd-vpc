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

	KeySWPortAddID       = "switch.port.add.id"
	KeySWPortAddL2ID     = "switch.port.add.l2-id"
	KeySWPortAddL2Name   = "switch.port.add.l2-name"
	KeySWPortAddMAC      = "switch.port.add.mac"
	KeySWPortAddSwitchID = "switch.port.add.switch-id"
	KeySWPortAddUplink   = "switch.port.add.uplink"

	KeySWPortRemovePortID   = "switch.port.remove.port-id"
	KeySWPortRemoveSwitchID = "switch.port.remove.switch-id"

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
