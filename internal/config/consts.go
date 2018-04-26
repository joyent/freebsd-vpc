// Copyright (c) 2018 Joyent, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

package config

const (
	DefaultManDir = "./docs/man"
	ManSect       = 8

	DefaultMarkdownDir       = "./docs/md"
	DefaultMarkdownURLPrefix = "/command"

	KeyDocManDir            = "doc.mandir"
	KeyDocMarkdownDir       = "doc.markdown-dir"
	KeyDocMarkdownURLPrefix = "doc.markdown-url-prefix"

	KeyEthLinkConnectID     = "ethlink.connect.id"
	KeyEthLinkConnectL2Name = "ethlink.connect.l2-name"
	KeyEthLinkCreateID      = "ethlink.create.id"
	KeyEthLinkDestroyID     = "ethlink.destroy.ethlink-id"
	KeyEthLinkListSortBy    = "ethlink.list.sort-by"
	KeyEthLinkVTagID        = "ethlink.vtag.ethlink-id"
	KeyEthLinkGetVTag       = "ethlink.vtag.get-vtag"
	KeyEthLinkSetVTag       = "ethlink.vtag.set-vtag"

	KeyListObjCounts = "list.obj-counts"
	KeyListObjSortBy = "list.sort-by"
	KeyListObjType   = "list.type"

	KeyLogFormat    = "log.format"
	KeyLogLevel     = "log.level"
	KeyLogStats     = "log.stats"
	KeyLogTermColor = "log.use-color"

	KeyPGDatabase = "db.name"
	KeyPGUser     = "db.username"
	KeyPGPassword = "db.password"
	KeyPGHost     = "db.host"
	KeyPGPort     = "db.port"

	KeyHostifCreateID  = "hostif.create.id"
	KeyHostifDestroyID = "hostif.destroy.id"

	KeyMuxConnectInterfaceID = "mux.connect.interface-id"
	KeyMuxConnectMuxID       = "mux.connect.mux-id"
	KeyMuxCreateMuxID        = "mux.create.mux-id"
	KeyMuxDestroyMuxID       = "mux.destroy.mux-id"
	KeyMuxDisconnectMuxID    = "mux.disconnect.mux-id"
	KeyMuxListenAddr         = "mux.listen.addr"
	KeyMuxListenMuxID        = "mux.listen.mux-id"
	KeyMuxShowMuxID          = "mux.show.mux-id"

	KeySWPortAddEthLinkID          = "switch.port.add.ethlink-id"
	KeySWPortAddID                 = "switch.port.add.id"
	KeySWPortAddMAC                = "switch.port.add.mac"
	KeySWPortAddSwitchID           = "switch.port.add.switch-id"
	KeySWPortAddUplink             = "switch.port.add.uplink"
	KeySWPortConnectInterfaceID    = "switch.port.connect.interface-id"
	KeySWPortConnectPortID         = "switch.port.connect.port-id"
	KeySWPortDisconnectInterfaceID = "switch.port.disconnect.interface-id"
	KeySWPortDisconnectPortID      = "switch.port.disconnect.port-id"
	KeySWPortSetPortID             = "switch.port.set.port-id"
	KeySWPortSetVNI                = "switch.port.set.vni"
	KeySWPortUplinkPortID          = "switch.port.uplink.port-id"
	KeySWPortUplinkSwitchID        = "switch.port.uplink.switch-id"

	KeySWPortRemovePortID   = "switch.port.remove.port-id"
	KeySWPortRemoveSwitchID = "switch.port.remove.switch-id"

	KeyShellAutoCompBashDir = "shell.autocomplete.bash-dir"

	KeySWCreateSwitchID  = "switch.create.switch-id"
	KeySWCreateSwitchMAC = "switch.create.switch-mac"
	KeySWCreateVNI       = "switch.create.vni"
	KeySWDestroySwitchID = "switch.destroy.switch-id"

	KeyUseGoogleAgent = "general.enable-agent"
	KeyUsePager       = "general.use-pager"
	KeyUseUTC         = "general.utc"

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
