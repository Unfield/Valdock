package models

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

type ConfigModel struct {
	ID         string `json:"id"`
	InstanceID string `json:"instance_id"`

	GeneratedBy string    `json:"generated_by"`
	GeneratedAt time.Time `json:"generated_at"`

	IncludeFiles []string `json:"include_files,omitempty"`
	LoadModules  []string `json:"load_modules,omitempty"`

	BindAddresses          []string      `json:"bind_addresses,omitempty"`
	BindSourceAddr         string        `json:"bind_source_addr,omitempty"`
	ProtectedMode          bool          `json:"protected_mode"`
	EnableProtectedConfigs string        `json:"enable_protected_configs,omitempty"`
	EnableDebugCommand     string        `json:"enable_debug_command,omitempty"`
	EnableModuleCommand    string        `json:"enable_module_command,omitempty"`
	Port                   int           `json:"port"`
	TCPBacklog             int           `json:"tcp_backlog"`
	MPTCPEnabled           bool          `json:"mptcp"`
	UnixSocket             string        `json:"unix_socket,omitempty"`
	UnixSocketGroup        string        `json:"unix_socket_group,omitempty"`
	UnixSocketPerm         int           `json:"unix_socket_perm,omitempty"`
	Timeout                time.Duration `json:"timeout"`
	TCPKeepalive           time.Duration `json:"tcp_keepalive"`

	TLSEnabled             bool     `json:"tls_enabled"`
	TLSPort                int      `json:"tls_port,omitempty"`
	TLSCertFile            string   `json:"tls_cert_file,omitempty"`
	TLSKeyFile             string   `json:"tls_key_file,omitempty"`
	TLSKeyPass             string   `json:"tls_key_pass,omitempty"`
	TLSCACertFile          string   `json:"tls_ca_cert_file,omitempty"`
	TLSCACertDir           string   `json:"tls_ca_cert_dir,omitempty"`
	TLSProtocols           []string `json:"tls_protocols,omitempty"`
	TLSPreferServerCiphers bool     `json:"tls_prefer_server_ciphers"`
	TLSSessionCacheSize    int      `json:"tls_session_cache_size"`
	TLSSessionTimeout      int      `json:"tls_session_timeout"`

	Daemonize           bool   `json:"daemonize"`
	SupervisedMode      string `json:"supervised_mode"`
	PIDFile             string `json:"pid_file"`
	LogLevel            string `json:"log_level"`
	LogFormat           string `json:"log_format"`
	LogTimestampFormat  string `json:"log_timestamp_format"`
	LogFile             string `json:"log_file"`
	SyslogEnabled       bool   `json:"syslog_enabled"`
	SyslogIdent         string `json:"syslog_ident"`
	SyslogFacility      string `json:"syslog_facility"`
	Databases           int    `json:"databases"`
	AlwaysShowLogo      bool   `json:"always_show_logo"`
	HideUserDataFromLog bool   `json:"hide_user_data_from_log"`
	SetProcTitle        bool   `json:"set_proc_title"`
	ProcTitleTemplate   string `json:"proc_title_template"`

	SavePoints              [][]int `json:"save_points,omitempty"`
	StopWritesOnBgsaveError bool    `json:"stop_writes_on_bgsave_error"`
	RDBCompression          bool    `json:"rdb_compression"`
	RDBChecksum             bool    `json:"rdb_checksum"`
	RDBVersionCheck         string  `json:"rdb_version_check"`
	DBFilename              string  `json:"db_filename"`
	RDBDelSyncFiles         bool    `json:"rdb_del_sync_files"`
	DataDir                 string  `json:"data_dir"`

	ReplicaOf             string `json:"replica_of,omitempty"`
	PrimaryAuth           string `json:"primary_auth,omitempty"`
	PrimaryUser           string `json:"primary_user,omitempty"`
	ReplicaServeStaleData bool   `json:"replica_serve_stale_data"`
	ReplicaReadOnly       bool   `json:"replica_read_only"`
	ReplDisklessSync      bool   `json:"repl_diskless_sync"`
	ReplDisklessDelay     int    `json:"repl_diskless_sync_delay"`
	ReplicaPriority       int    `json:"replica_priority"`
	ReplDisableTcpNodelay bool   `json:"repl_disable_tcp_nodelay"`

	ACLFile            string            `json:"acl_file,omitempty"`
	RequirePass        string            `json:"require_pass,omitempty"`
	DefaultChannelPerm string            `json:"default_channel_perm,omitempty"`
	Users              map[string]string `json:"users,omitempty"`

	AppendOnly             bool   `json:"append_only"`
	AppendFilename         string `json:"append_filename"`
	AppendDirname          string `json:"append_dirname"`
	AppendFsyncPolicy      string `json:"append_fsync_policy"`
	NoAppendFsyncOnRewrite bool   `json:"no_append_fsync_on_rewrite"`
	AutoAOFRewritePct      int    `json:"auto_aof_rewrite_pct"`
	AutoAOFRewriteMinSize  string `json:"auto_aof_rewrite_min_size"`
	AOFLoadTruncated       bool   `json:"aof_load_truncated"`
	AOFUseRDBPreamble      bool   `json:"aof_use_rdb_preamble"`

	MaxMemory              string `json:"max_memory,omitempty"`
	MaxMemoryPolicy        string `json:"max_memory_policy,omitempty"`
	MaxMemorySamples       int    `json:"max_memory_samples"`
	EvictionTenacity       int    `json:"eviction_tenacity"`
	ReplicaIgnoreMaxMemory bool   `json:"replica_ignore_max_memory"`

	LazyFreeEviction bool `json:"lazy_free_eviction"`
	LazyFreeExpire   bool `json:"lazy_free_expire"`
	LazyServerDel    bool `json:"lazy_server_del"`
	LazyUserDel      bool `json:"lazy_user_del"`
	LazyUserFlush    bool `json:"lazy_user_flush"`

	IOThreads         int    `json:"io_threads"`
	OOMScoreAdj       string `json:"oom_score_adj"`
	OOMScoreAdjValues []int  `json:"oom_score_adj_values"`
	DisableTHP        bool   `json:"disable_thp"`

	ClusterEnabled               bool   `json:"cluster_enabled"`
	ClusterConfigFile            string `json:"cluster_config_file"`
	ClusterNodeTimeout           int    `json:"cluster_node_timeout"`
	ClusterRequireFullCoverage   bool   `json:"cluster_require_full_coverage"`
	ClusterAllowReadsWhenDown    bool   `json:"cluster_allow_reads_when_down"`
	ClusterReplicaNoFailover     bool   `json:"cluster_replica_no_failover"`
	ClusterMigrationBarrier      int    `json:"cluster_migration_barrier"`
	ClusterAllowReplicaMigration bool   `json:"cluster_allow_replica_migration"`

	CommandLogExecutionSlowMicros int `json:"commandlog_execution_slower_than"`
	CommandLogSlowMaxLen          int `json:"commandlog_slow_max_len"`
	CommandLogReqLargerThan       int `json:"commandlog_request_larger_than"`
	CommandLogLargeReqMaxLen      int `json:"commandlog_large_req_max_len"`
	CommandLogReplyLargerThan     int `json:"commandlog_reply_larger_than"`
	CommandLogLargeReplyMaxLen    int `json:"commandlog_large_reply_max_len"`

	LatencyMonitorThreshold int       `json:"latency_monitor_threshold"`
	LatencyTracking         bool      `json:"latency_tracking"`
	LatencyInfoPercentiles  []float64 `json:"latency_info_percentiles,omitempty"`

	NotifyKeyspaceEvents string `json:"notify_keyspace_events"`

	HashMaxListpackEntries int `json:"hash_max_listpack_entries"`
	HashMaxListpackValue   int `json:"hash_max_listpack_value"`
	ListMaxListpackSize    int `json:"list_max_listpack_size"`
	ListCompressDepth      int `json:"list_compress_depth"`
	SetMaxIntsetEntries    int `json:"set_max_intset_entries"`
	SetMaxListpackEntries  int `json:"set_max_listpack_entries"`
	SetMaxListpackValue    int `json:"set_max_listpack_value"`
	ZSetMaxListpackEntries int `json:"zset_max_listpack_entries"`
	ZSetMaxListpackValue   int `json:"zset_max_listpack_value"`

	MaxClients               int                 `json:"max_clients"`
	ClientOutputBufferLimits map[string][]string `json:"client_output_buffer_limits"`
	ClientQueryBufferLimit   string              `json:"client_query_buffer_limit"`
	ClientEvictionThreshold  string              `json:"maxmemory_clients"`

	HLLSparseMaxBytes    int `json:"hll_sparse_max_bytes"`
	StreamNodeMaxBytes   int `json:"stream_node_max_bytes"`
	StreamNodeMaxEntries int `json:"stream_node_max_entries"`

	ActiveRehashing  bool `json:"active_rehashing"`
	ActiveDefrag     bool `json:"active_defrag"`
	JemallocBGThread bool `json:"jemalloc_bg_thread"`

	HZ               int    `json:"hz"`
	AvailabilityZone string `json:"availability_zone"`
}

func yesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func kv(key, value string) string {
	if value == "" {
		return ""
	}
	return fmt.Sprintf("%s %s\n", key, value)
}

func (c *ConfigModel) ToConf() string {
	var sb strings.Builder

	sb.WriteString("# ---- Generated by Valdock ----\n\n")

	sb.WriteString("# Persistence\n")
	sb.WriteString(kv("appendonly", yesNo(c.AppendOnly)))
	sb.WriteString(kv("appendfsync", c.AppendFsyncPolicy))
	sb.WriteString(kv("aof-use-rdb-preamble", yesNo(c.AOFUseRDBPreamble)))
	sb.WriteString(kv("stop-writes-on-bgsave-error", yesNo(c.StopWritesOnBgsaveError)))
	sb.WriteString(kv("rdbcompression", yesNo(c.RDBCompression)))
	sb.WriteString(kv("rdbchecksum", yesNo(c.RDBChecksum)))
	sb.WriteString(kv("auto-aof-rewrite-percentage", fmt.Sprint(c.AutoAOFRewritePct)))
	sb.WriteString(kv("auto-aof-rewrite-min-size", c.AutoAOFRewriteMinSize))
	sb.WriteString(kv("dir", c.DataDir))
	sb.WriteString(kv("dbfilename", c.DBFilename))
	sb.WriteString(kv("appenddirname", c.AppendDirname))
	sb.WriteString("\n")

	sb.WriteString("# Memory & Performance\n")
	sb.WriteString(kv("maxmemory", c.MaxMemory))
	sb.WriteString(kv("maxmemory-policy", c.MaxMemoryPolicy))
	sb.WriteString(kv("activedefrag", yesNo(c.ActiveDefrag)))
	sb.WriteString(kv("activerehashing", yesNo(c.ActiveRehashing)))
	sb.WriteString(kv("io-threads", fmt.Sprint(c.IOThreads)))
	sb.WriteString(kv("jemalloc-bg-thread", yesNo(c.JemallocBGThread)))
	sb.WriteString(kv("hz", fmt.Sprint(c.HZ)))
	sb.WriteString("\n")

	sb.WriteString("# Security\n")
	sb.WriteString(kv("requirepass", c.RequirePass))
	sb.WriteString(kv("hide-user-data-from-log", yesNo(c.HideUserDataFromLog)))
	if c.ACLFile != "" {
		sb.WriteString(kv("aclfile", c.ACLFile))
	}
	sb.WriteString("\n")

	sb.WriteString("# Logging / Clients\n")
	sb.WriteString(kv("loglevel", c.LogLevel))
	sb.WriteString(kv("logfile", c.LogFile))
	sb.WriteString(kv("syslog-enabled", yesNo(c.SyslogEnabled)))
	sb.WriteString(kv("maxclients", fmt.Sprint(c.MaxClients)))

	if len(c.ClientOutputBufferLimits) > 0 {
		classes := make([]string, 0, len(c.ClientOutputBufferLimits))
		for k := range c.ClientOutputBufferLimits {
			classes = append(classes, k)
		}
		sort.Strings(classes)
		for _, class := range classes {
			v := c.ClientOutputBufferLimits[class]
			if len(v) == 3 {
				sb.WriteString(fmt.Sprintf(
					"client-output-buffer-limit %s %s %s %s\n",
					class, v[0], v[1], v[2],
				))
			}
		}
	}
	sb.WriteString("\n")

	sb.WriteString("# Snapshotting (optional)\n")
	if len(c.SavePoints) > 0 {
		var parts []string
		for _, p := range c.SavePoints {
			if len(p) == 2 {
				parts = append(parts, fmt.Sprintf("%d %d", p[0], p[1]))
			}
		}
		sb.WriteString(fmt.Sprintf("save %s\n", strings.Join(parts, " ")))
	}
	sb.WriteString("\n")

	return sb.String()
}

func (c *ConfigModel) DebugDump() string {
	v := reflect.ValueOf(*c)
	t := v.Type()
	var lines []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := t.Field(i).Name
		if isZero(field) {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s = %v", name, field.Interface()))
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int64:
		return v.Int() == 0
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}
	return false
}
