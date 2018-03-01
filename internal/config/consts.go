package config

const (
	DefaultManDir = "./docs/man"
	ManSect       = 8

	DefaultMarkdownDir       = "./docs/md"
	DefaultMarkdownURLPrefix = "/command"

	KeyDocManDir            = "doc.mandir"
	KeyDocMarkdownDir       = "doc.markdown-dir"
	KeyDocMarkdownURLPrefix = "doc.markdown-url-prefix"

	KeyEthLinkDestroyID  = "ethlink.destroy.ethlink-id"
	KeyEthLinkListSortBy = "ethlink.list.sort-by"

	KeyListObjCounts = "list.obj-counts"
	KeyListObjSortBy = "list.sort-by"
	KeyListObjType   = "list.type"

	KeyLogFormat    = "log.format"
	KeyLogLevel     = "log.level"
	KeyLogStats     = "log.stats"
	KeyLogTermColor = "log.use-color"

	KeySWPortAddID              = "switch.port.add.id"
	KeySWPortAddEthLinkID       = "switch.port.add.ethlink-id"
	KeySWPortAddL2Name          = "switch.port.add.l2-name"
	KeySWPortAddMAC             = "switch.port.add.mac"
	KeySWPortAddSwitchID        = "switch.port.add.switch-id"
	KeySWPortAddUplink          = "switch.port.add.uplink"
	KeySWPortConnectInterfaceID = "switch.port.connect.interface-id"
	KeySWPortConnectPortID      = "switch.port.connect.port-id"

	KeySWPortRemovePortID   = "switch.port.remove.port-id"
	KeySWPortRemoveSwitchID = "switch.port.remove.switch-id"

	KeyShellAutoCompBashDir = "shell.autocomplete.bash-dir"

	KeySWCreateSwitchID  = "switch.create.switch-id"
	KeySWCreateSwitchMAC = "switch.create.switch-mac"
	KeySWCreateVNI       = "switch.create.vni"
	KeySWDestroySwitchID = "switch.destroy.switch-id"

	KeyUsePager = "general.use-pager"
	KeyUseUTC   = "general.utc"

	KeyVMNICCreateID    = "vmnic.create.id"
	KeyVMNICCreateMAC   = "vmnic.create.mac"
	KeyVMNICDestroyID   = "vmnic.destroy.id"
	KeyVMNICGetNQueues  = "vmnic.get.num-queues"
	KeyVMNICGetVMNICID  = "vmnic.get.vmnic-id"
	KeyVMNICSetFreeze   = "vmnic.set.freeze"
	KeyVMNICSetNQueues  = "vmnic.set.num-queues"
	KeyVMNICSetUnfreeze = "vmnic.set.unfreeze"
	KeyVMNICSetVMNICID  = "vmnic.set.vmnic-id"
)
