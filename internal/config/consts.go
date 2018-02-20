package config

const (
	DefaultManDir = "./docs/man"
	ManSect       = 8

	DefaultMarkdownDir       = "./docs/md"
	DefaultMarkdownURLPrefix = "/command"

	KeyDocManDir            = "doc.mandir"
	KeyDocMarkdownDir       = "doc.markdown-dir"
	KeyDocMarkdownURLPrefix = "doc.markdown-url-prefix"

	KeyShellAutoCompBashDir = "shell.autocomplete.bash-dir"

	KeySWCreateID  = "switch.create.id"
	KeySWCreateVNI = "switch.create.vni"
	KeySWDestroyID = "switch.destroy.id"

	KeyUsePager = "general.use-pager"
	KeyUseUTC   = "general.utc"

	KeyLogFormat    = "log.format"
	KeyLogLevel     = "log.level"
	KeyLogStats     = "log.stats"
	KeyLogTermColor = "log.use-color"
)
