package configtemplates

import (
	"fmt"
	"time"

	"github.com/Unfield/Valdock/models"
	"github.com/Unfield/Valdock/utils"
	"github.com/Unfield/Valdock/version"
)

type DefaultConfigOptions struct {
	MaxMemoryMB int
	MaxClients  int
}

func NewDefaultConfig(instanceID string, options DefaultConfigOptions) *models.ConfigModel {
	configID, _ := utils.GenerateID()

	if options.MaxMemoryMB <= 0 {
		options.MaxMemoryMB = 256
	}
	options.MaxMemoryMB = (options.MaxMemoryMB * 75) / 100

	if options.MaxClients < 1 {
		options.MaxClients = 10000
	}

	return &models.ConfigModel{
		ID:         configID,
		InstanceID: instanceID,

		GeneratedBy: fmt.Sprintf("Valdock %s", version.FullVersion()),
		GeneratedAt: time.Now(),

		AppendOnly:              true,
		AppendFsyncPolicy:       "everysec",
		AOFUseRDBPreamble:       true,
		StopWritesOnBgsaveError: true,
		RDBCompression:          true,
		RDBChecksum:             true,
		AutoAOFRewritePct:       100,
		AutoAOFRewriteMinSize:   "64mb",
		DataDir:                 "/data",
		DBFilename:              "dump.rdb",
		AppendDirname:           "appendonlydir",

		MaxMemory:        formatMemory(options.MaxMemoryMB),
		MaxMemoryPolicy:  "allkeys-lru",
		ActiveDefrag:     true,
		ActiveRehashing:  true,
		IOThreads:        2,
		JemallocBGThread: true,
		HZ:               10,

		ACLFile: "/data/users.acl",

		LogFile:       "",
		LogLevel:      "notice",
		SyslogEnabled: false,

		MaxClients: options.MaxClients,
		ClientOutputBufferLimits: map[string][]string{
			"pubsub":  {"32mb", "8mb", "60"},
			"replica": {"256mb", "64mb", "60"},
		},
	}
}
