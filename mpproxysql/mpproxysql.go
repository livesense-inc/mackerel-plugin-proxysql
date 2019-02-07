package mpproxysql

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/ziutek/mymysql/mysql"

	// MySQL Driver
	_ "github.com/ziutek/mymysql/native"
)

// ProxySQLPlugin mackerel plugin for ProxySQL
type ProxySQLPlugin struct {
	Target               string
	Tempfile             string
	prefix               string
	Username             string
	Password             string
	isUnixSocket         bool
	EnableMonitorStats   bool
	EnableConnectionPool bool
}

const defaultMetricKeyPrefix = "proxysql"
const defaultMetricName = "ProxySQL"

// MetricKeyPrefix retruns the metrics key prefix
func (p ProxySQLPlugin) MetricKeyPrefix() string {
	if p.prefix == "" {
		p.prefix = defaultMetricKeyPrefix
	}
	return p.prefix
}

// LabelPrefix retruns the metrics label prefix
func (p ProxySQLPlugin) LabelPrefix() string {
	return strings.Title(strings.Replace(p.MetricKeyPrefix(), defaultMetricKeyPrefix, defaultMetricName, -1))
}

/*
StatsProxySQLGlobal(default) Metrics
*/
func (p ProxySQLPlugin) defaultGraphdef() map[string]mp.Graphs {
	labelPrefix := p.LabelPrefix()

	return map[string]mp.Graphs{
		"transactions": {
			Label: labelPrefix + " Transactions",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "Active_Transactions", Label: "Active", Diff: true, Stacked: false},
			},
		},
		"connections": {
			Label: labelPrefix + " Connections",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "Client_Connections_aborted", Label: "Client aborted", Diff: true, Stacked: false},
				{Name: "Client_Connections_connected", Label: "Client connected", Diff: false, Stacked: false},
				{Name: "Client_Connections_created", Label: "Client created", Diff: true, Stacked: false},
				{Name: "Client_Connections_non_idle", Label: "Client non idle", Diff: false, Stacked: false},
				{Name: "Server_Connections_aborted", Label: "Server aborted", Diff: true, Stacked: false},
				{Name: "Server_Connections_connected", Label: "Server connected", Diff: false, Stacked: false},
				{Name: "Server_Connections_created", Label: "Server created", Diff: true, Stacked: false},
				{Name: "Server_Connections_delayed", Label: "Server delayed", Diff: true, Stacked: false},
			},
		},
		"traffic": {
			Label: labelPrefix + " Traffic",
			Unit:  "bytes/sec",
			Metrics: []mp.Metrics{
				{Name: "Queries_backends_bytes_sent", Label: "Sent to Backend Bytes", Diff: true, Stacked: false},
				{Name: "Queries_backends_bytes_recv", Label: "Received from Backend Bytes", Diff: true, Stacked: false},
				{Name: "Queries_frontends_bytes_sent", Label: "Sent to Frontend Bytes", Diff: true, Stacked: false},
				{Name: "Queries_frontends_bytes_recv", Label: "Received from Frontend Bytes", Diff: true, Stacked: false},
			},
		},
		"processor": {
			Label: labelPrefix + " Processor Time (sec)",
			Unit:  "float",
			Metrics: []mp.Metrics{
				// convert from Query_Processor_time_nsec
				{Name: "Query_Processor_time_sec", Label: "Query Processor Time", Diff: true, Stacked: false},
				// convert from Backend_query_time_nsec
				{Name: "Backend_query_time_sec", Label: "Backend Query Time", Diff: true, Stacked: false},
			},
		},
		"buffers": {
			Label: labelPrefix + " Buffers",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "Query_Cache_Memory_bytes", Label: "Query Cache Memory", Diff: false, Stacked: true},
				{Name: "ConnPool_memory_bytes", Label: "Connection Pool", Diff: false, Stacked: true},
				{Name: "SQLite3_memory_bytes", Label: "SQLite", Diff: false, Stacked: true},
				{Name: "mysql_backend_buffers_bytes", Label: "Backend", Diff: false, Stacked: true},
				{Name: "mysql_frontend_buffers_bytes", Label: "Frontend", Diff: false, Stacked: true},
				{Name: "mysql_session_internal_bytes", Label: "Session Internal", Diff: false, Stacked: true},
			},
		},
		"commands": {
			Label: labelPrefix + " Commands",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "Com_autocommit", Label: "Autocommit", Diff: true, Stacked: true},
				{Name: "Com_autocommit_filtered", Label: "Autocommit Filtered", Diff: true, Stacked: true},
				{Name: "Com_commit", Label: "Commit", Diff: true, Stacked: true},
				{Name: "Com_commit_filtered", Label: "Commit Filtered", Diff: true, Stacked: true},
				{Name: "Com_rollback", Label: "Rollback", Diff: true, Stacked: true},
				{Name: "Com_rollback_filtered", Label: "Rollback Filtered", Diff: true, Stacked: true},
				{Name: "Com_backend_change_user", Label: "Backend Change User", Diff: true, Stacked: true},
				{Name: "Com_backend_init_db", Label: "Backend Init DB", Diff: true, Stacked: true},
				{Name: "Com_backend_set_names", Label: "Backend Set Names", Diff: true, Stacked: true},
				{Name: "Com_frontend_init_db", Label: "Frontend Init DB", Diff: true, Stacked: true},
				{Name: "Com_frontend_set_names", Label: "Frontend Set Names", Diff: true, Stacked: true},
				{Name: "Com_frontend_use_db", Label: "Frontend Use DB", Diff: true, Stacked: true},
				{Name: "Com_backend_stmt_prepare", Label: "Backend STMT Prepare", Diff: true, Stacked: true},
				{Name: "Com_backend_stmt_execute", Label: "Backend STMT Execute", Diff: true, Stacked: true},
				{Name: "Com_backend_stmt_close", Label: "Backend STMT Close", Diff: true, Stacked: true},
				{Name: "Com_frontend_stmt_prepare", Label: "Frontend STMT Prepare", Diff: true, Stacked: true},
				{Name: "Com_frontend_stmt_execute", Label: "Frontend STMT Execute", Diff: true, Stacked: true},
				{Name: "Com_frontend_stmt_close", Label: "Frontend STMT Close", Diff: true, Stacked: true},
				{Name: "Questions", Label: "Questions", Diff: true, Stacked: true},
			},
		},
		"slow_queries": {
			Label: labelPrefix + " Slow Queries",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "Slow_queries", Label: "total", Diff: true, Stacked: false},
			},
		},
		"gtids": {
			Label: labelPrefix + " GTID",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "GTID_consistent_queries", Label: "Consistent Queries", Diff: true, Stacked: false},
				{Name: "GTID_session_collected", Label: "Session Collected", Diff: true, Stacked: false},
			},
		},
		"access_errors": {
			Label: labelPrefix + " Access Errors",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "Access_Denied_Wrong_Password", Label: "Access Denied Wrong Password", Diff: true, Stacked: false},
				{Name: "Access_Denied_Max_Connections", Label: "Access Denied Max Connections", Diff: true, Stacked: false},
				{Name: "Access_Denied_Max_User_Connections", Label: "Access Denied Max User Connections", Diff: true, Stacked: false},
			},
		},
		"threads": {
			Label: labelPrefix + " Threads",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "MySQL_Thread_Workers", Label: "Thread", Diff: false, Stacked: false},
				{Name: "MySQL_Monitor_Workers", Label: "Monitor", Diff: false, Stacked: false},
			},
		},
		"connpool_summary": {
			Label: labelPrefix + " Connection Pool Summary",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "ConnPool_get_conn_immediate", Label: "Get Immediate", Diff: true, Stacked: false},
				{Name: "ConnPool_get_conn_success", Label: "Get Success", Diff: true, Stacked: false},
				{Name: "ConnPool_get_conn_failure", Label: "Get Failure", Diff: true, Stacked: false},
			},
		},
		"statement": {
			Label: labelPrefix + " Statement",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "Stmt_Client_Active_Total", Label: "Client Active Total", Diff: false, Stacked: false},
				{Name: "Stmt_Client_Active_Unique", Label: "Client Active Unique", Diff: false, Stacked: false},
				{Name: "Stmt_Server_Active_Total", Label: "Server Active Total", Diff: false, Stacked: false},
				{Name: "Stmt_Server_Active_Unique", Label: "Server Active Unique", Diff: false, Stacked: false},
				{Name: "Stmt_Cached", Label: "Cached", Diff: false, Stacked: false},
			},
		},
		"querycache_count": {
			Label: labelPrefix + " Query Cache Count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "Query_Cache_count_GET", Label: "Get", Diff: true, Stacked: false},
				{Name: "Query_Cache_count_GET_OK", Label: "Get OK", Diff: true, Stacked: false},
				{Name: "Query_Cache_count_SET", Label: "Set", Diff: true, Stacked: false},
			},
		},
		"querycache_bytes": {
			Label: labelPrefix + " Query Cache Bytes",
			Unit:  "bytes/sec",
			Metrics: []mp.Metrics{
				{Name: "Query_Cache_bytes_IN", Label: "IN", Diff: true, Stacked: false},
				{Name: "Query_Cache_bytes_OUT", Label: "OUT", Diff: true, Stacked: false},
			},
		},
		"querycache_entries": {
			Label: labelPrefix + " Query Cache Entries/Purged",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "Query_Cache_Entries", Label: "Entries", Diff: false, Stacked: true},
				{Name: "Query_Cache_Purged", Label: "Purged", Diff: true, Stacked: false},
			},
		},
	}
}

func (p ProxySQLPlugin) fetchStatsProxySQLGlobal(db mysql.Conn, stat map[string]float64) error {
	rows, _, err := db.Query("select * from stats_mysql_global")
	if err != nil {
		log.Fatalln("FetchMetrics (proxysql_global): ", err)
		return err
	}

	for _, row := range rows {
		if len(row) > 1 {
			variableName := string(row[0].([]byte))
			stat[variableName], _ = atof(string(row[1].([]byte)))
		} else {
			log.Fatalln("FetchMetrics (proxysql_global): row length is too small: ", len(row))
		}
	}
	return nil
}

/*
Monitor Metrics
*/
func (p ProxySQLPlugin) addGraphdefMonitorStats(graphdef map[string]mp.Graphs) map[string]mp.Graphs {
	labelPrefix := p.LabelPrefix()

	graphdef["monitor.workers"] = mp.Graphs{
		Label: labelPrefix + " Monitor Workers",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "MySQL_Monitor_Workers", Label: "workers", Diff: false, Stacked: false},
		},
	}
	graphdef["monitor.checks"] = mp.Graphs{
		Label: labelPrefix + " Monitor Checks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "MySQL_Monitor_ping_check_OK", Label: "Ping OK", Diff: true, Stacked: false},
			{Name: "MySQL_Monitor_ping_check_ERR", Label: "Ping Error", Diff: true, Stacked: false},
			{Name: "MySQL_Monitor_read_only_check_OK", Label: "Readonly check OK", Diff: true, Stacked: false},
			{Name: "MySQL_Monitor_read_only_check_ERR", Label: "Readonly check Error", Diff: true, Stacked: false},
			{Name: "MySQL_Monitor_replication_lag_check_OK", Label: "Replication lag check OK", Diff: true, Stacked: false},
			{Name: "MySQL_Monitor_replication_lag_check_ERR", Label: "Replication lag check Error", Diff: true, Stacked: false},
		},
	}
	return graphdef
}

/*
ConnectionPool Metrics
*/
func (p ProxySQLPlugin) addGraphdefConnectionPool(graphdef map[string]mp.Graphs) map[string]mp.Graphs {
	labelPrefix := p.LabelPrefix()

	graphdef["connpool.conns.#"] = mp.Graphs{
		Label: labelPrefix + " ConnectionPool Connections",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ConnUsed", Label: "Used", Diff: false, Stacked: false},
			{Name: "ConnFree", Label: "Free", Diff: false, Stacked: false},
		},
	}
	graphdef["connpool.stats.#"] = mp.Graphs{
		Label: labelPrefix + " ConnectionPool Stats",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ConnOK", Label: "OK", Diff: true, Stacked: false},
			{Name: "ConnERR", Label: "Error", Diff: true, Stacked: false},
		},
	}
	graphdef["connpool.transfers.#"] = mp.Graphs{
		Label: labelPrefix + " ConnectionPool Transfers",
		Unit:  "bytes/sec",
		Metrics: []mp.Metrics{
			{Name: "Bytes_data_sent", Label: "Sent", Diff: true, Stacked: false},
			{Name: "Bytes_data_recv", Label: "Received", Diff: true, Stacked: false},
		},
	}
	graphdef["connpool.queries.#"] = mp.Graphs{
		Label: labelPrefix + " ConnectionPool Queries",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "Queries", Label: "Queries", Diff: true, Stacked: false},
		},
	}
	graphdef["connpool.perf.#"] = mp.Graphs{
		Label: labelPrefix + " ConnectionPool Performance",
		Unit:  "float",
		Metrics: []mp.Metrics{
			// convert from Latency_us
			{Name: "Latency_sec", Label: "Latency sec", Diff: false, Stacked: false},
		},
	}
	return graphdef
}

func (p ProxySQLPlugin) fetchConnectionPool(db mysql.Conn, stat map[string]float64) error {
	rows, _, err := db.Query("select hostgroup, srv_host, srv_port, ConnUsed, ConnFree, ConnOK, ConnERR, Queries, Bytes_data_sent, Bytes_data_recv, Latency_us from stats_mysql_connection_pool")
	if err != nil {
		log.Fatalln("FetchMetrics (proxysql_connpool): ", err)
		return err
	}

	for _, row := range rows {
		if len(row) == 11 {
			hostgroup := string(row[0].([]byte))
			host := string(row[1].([]byte))
			port := string(row[2].([]byte))
			srvName := fmt.Sprintf("%s-%s-%s", hostgroup, strings.Replace(host, ".", "_", -1), port)

			stat[fmt.Sprintf("connpool.conns.%s.ConnUsed", srvName)], _ = atof(string(row[3].([]byte)))
			stat[fmt.Sprintf("connpool.conns.%s.ConnFree", srvName)], _ = atof(string(row[4].([]byte)))

			stat[fmt.Sprintf("connpool.stats.%s.ConnOK", srvName)], _ = atof(string(row[5].([]byte)))
			stat[fmt.Sprintf("connpool.stats.%s.ConnERR", srvName)], _ = atof(string(row[6].([]byte)))

			stat[fmt.Sprintf("connpool.queries.%s.Queries", srvName)], _ = atof(string(row[7].([]byte)))

			stat[fmt.Sprintf("connpool.transfers.%s.Bytes_data_sent", srvName)], _ = atof(string(row[8].([]byte)))
			stat[fmt.Sprintf("connpool.transfers.%s.Bytes_data_recv", srvName)], _ = atof(string(row[9].([]byte)))

			stat[fmt.Sprintf("connpool.perf.%s.Latency_usec", srvName)], _ = atof(string(row[10].([]byte)))

		} else {
			log.Fatalln("FetchMetrics (proxysql_connpool): row length is invalid: ", len(row))
		}
	}
	return nil
}

func (p ProxySQLPlugin) calculateSec(stat map[string]float64) {
	for k, v := range stat {
		if strings.HasSuffix(k, "_usec") {
			newKey := strings.TrimSuffix(k, "_usec") + "_sec"
			stat[newKey] = v / 1000 / 1000
			delete(stat, k)
		} else if strings.HasSuffix(k, "_nsec") {
			newKey := strings.TrimSuffix(k, "_nsec") + "_sec"
			stat[newKey] = v / 1000 / 1000 / 1000
			delete(stat, k)
		}
	}
}

// FetchMetrics interface for mackerelplugin
func (p ProxySQLPlugin) FetchMetrics() (map[string]interface{}, error) {
	proto := "tcp"
	if p.isUnixSocket {
		proto = "unix"
	}
	db := mysql.New(proto, "", p.Target, p.Username, p.Password, "")
	err := db.Connect()
	if err != nil {
		log.Fatalln("FetchMetrics (DB Connect): ", err)
		return nil, err
	}
	defer db.Close()

	stat := make(map[string]float64)
	p.fetchStatsProxySQLGlobal(db, stat)
	if p.EnableConnectionPool {
		p.fetchConnectionPool(db, stat)
	}

	p.calculateSec(stat)

	statRet := make(map[string]interface{})
	for key, value := range stat {
		statRet[key] = value
	}
	return statRet, err
}

// GraphDefinition interface for mackerelplugin
func (p ProxySQLPlugin) GraphDefinition() map[string]mp.Graphs {
	graphdef := p.defaultGraphdef()
	if p.EnableMonitorStats {
		graphdef = p.addGraphdefMonitorStats(graphdef)
	}
	if p.EnableConnectionPool {
		graphdef = p.addGraphdefConnectionPool(graphdef)
	}
	return graphdef
}

func atof(str string) (float64, error) {
	str = strings.Replace(str, ",", "", -1)
	str = strings.Replace(str, ";", "", -1)
	str = strings.Replace(str, "/s", "", -1)
	str = strings.Trim(str, " ")
	return strconv.ParseFloat(str, 64)
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "127.0.0.1", "Hostname")
	optPort := flag.String("port", "6032", "Port")
	optSocket := flag.String("socket", "", "Unix socket path")
	optUser := flag.String("username", "stats", "Username")
	optPass := flag.String("password", "stats", "Password")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optEnableMonitorStats := flag.Bool("monitor-stats", false, "Enable Monitor Stats metrics")
	optEnableConnectionPool := flag.Bool("connection-pool", false, "Enable Connection Pool metrics")
	optMetricKeyPrefix := flag.String("metric-key-prefix", "proxysql", "metric key prefix")
	flag.Parse()

	var proxysql ProxySQLPlugin

	if *optSocket != "" {
		proxysql.Target = *optSocket
		proxysql.isUnixSocket = true
	} else {
		proxysql.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	}
	proxysql.Username = *optUser
	proxysql.Password = *optPass
	proxysql.EnableMonitorStats = *optEnableMonitorStats
	proxysql.EnableConnectionPool = *optEnableConnectionPool
	proxysql.prefix = *optMetricKeyPrefix
	helper := mp.NewMackerelPlugin(proxysql)
	helper.Tempfile = *optTempfile
	helper.Run()
}
