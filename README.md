mackerel-plugin-proxysql
=====================

ProxySQL custom metrics plugin for mackerel.io agent.

## Install

You need:

- [mkr](https://github.com/mackerelio/mkr)

```
$ sudo mkr plugin install livesense-inc/mackerel-plugin-proxysql
```

## Synopsis

```shell
mackerel-plugin-proxysql [-connection-pool] [-monitor-stats] [-host=<host> [-port=<port>]|-socket=<sockfile>] [-username=<username>] [-password=<password>] [-tempfile=<tempfile>] [-metric-key-prefix=<prefix>]
```

Options:

- `-connection-pool`: Enable connection pool metrics (default: disable)
- `-monitor-stats` : Enable monitor stats metrics (default: disable)
- `-host`: ProxySQL admin intereface hostname or IP address (default: `127.0.0.1`)
- `-port`: ProxySQL admin interface port number (default: `6032`)
- `-socket`: Unix socket path. if defined, use Unix socket (default: undefined)
- `-username`: ProxySQL admin interface user (default: `stats`)
- `-password`: ProxySQL admin interface password (default: `stats`)
- `-tempfile` : Override tempfile path (default: mackerel default)
- `-metric-key-prefix` : Override metrics key prefix (default: `proxysql`)

## Example of mackerel-agent.conf

```
[plugin.metrics.proxysql]
command = "/path/to/mackerel-plugin-proxysql -connection-pool"
```

## Graphs and Metrics 

- Prefix word `ProxySQL` or `proxysql` can be replaced with `-metric-key-prefix` option
- The desctiptions of each metrics are in [ProxySQL Official wiki](https://github.com/sysown/proxysql/wiki/STATS-(statistics))

### ProxySQL Transactions (proxysql.transactions)

- Active (proxysql.transactions.Active_Transactions)

### ProxySQL Connections (proxysql.connections)

- Client aborted (proxysql.connections.Client_Connections_aborted)
- Client connected (proxysql.connections.Client_Connections_connected)
- Client created (proxysql.connections.Client_Connections_created)
- Client non idle (proxysql.connections.Client_Connections_non_idle)
- Server aborted (proxysql.connections.Server_Connections_aborted)
- Server connected (proxysql.connections.Server_Connections_connected)
- Server created (proxysql.connections.Server_Connections_created)
- Server delayed (proxysql.connections.Server_Connections_delayed)

### ProxySQL Traffic (proxysql.traffic)

- Sent to Backend Bytes (proxysql.traffic.Queries_backends_bytes_sent)
- Received from Backend Bytes (proxysql.traffic.Queries_backends_bytes_recv)
- Sent to Frontend Bytes (proxysql.traffic.Queries_frontends_bytes_sent)
- Received from Frontend Bytes (proxysql.traffic.Queries_frontends_bytes_recv)

### ProxySQL Processor Time (sec) (proxysql.processor)

- Query Processor Time (proxysql.processor.Query_Processor_time_sec)
- Backend Query Time (proxysql.processor.Backend_query_time_sec)

### ProxySQL Buffers (proxysql.buffers)

- Query Cache Memory (proxysql.buffers.Query_Cache_Memory_bytes)
- Connection Pool (proxysql.buffers.ConnPool_memory_bytes)
- SQLite (proxysql.buffers.SQLite3_memory_bytes)
- Backend (proxysql.buffers.mysql_backend_buffers_bytes)
- Frontend (proxysql.buffers.mysql_frontend_buffers_bytes)
- Session Internal (proxysql.buffers.mysql_session_internal_bytes)

### ProxySQL Commands (proxysql.commands)

- Autocommit (proxysql.commands.Com_autocommit)
- Autocommit Filtered (proxysql.commands.Com_autocommit_filtered)
- Commit (proxysql.commands.Com_commit)
- Commit Filtered (proxysql.commands.Com_commit_filtered)
- Rollback (proxysql.commands.Com_rollback)
- Rollback Filtered (proxysql.commands.Com_rollback_filtered)
- Backend Change User (proxysql.commands.Com_backend_change_user)
- Backend Init DB (proxysql.commands.Com_backend_init_db)
- Backend Set Names (proxysql.commands.Com_backend_set_names)
- Frontend Init DB (proxysql.commands.Com_frontend_init_db)
- Frontend Set Names (proxysql.commands.Com_frontend_set_names)
- Frontend Use DB (proxysql.commands.Com_frontend_use_db)
- Backend STMT Prepare (proxysql.commands.Com_backend_stmt_prepare)
- Backend STMT Execute (proxysql.commands.Com_backend_stmt_execute)
- Backend STMT Close (proxysql.commands.Com_backend_stmt_close)
- Frontend STMT Prepare (proxysql.commands.Com_frontend_stmt_prepare)
- Frontend STMT Execute (proxysql.commands.Com_frontend_stmt_execute)
- Frontend STMT Close (proxysql.commands.Com_frontend_stmt_close)
- Questions (proxysql.commands.Questions)

### ProxySQL Slow Queries (proxysql.slow_queries)

- total (proxysql.slow_queries.Slow_queries)

### ProxySQL GTID (proxysql.gtids)

- Consistent Queries (proxysql.gtids.GTID_consistent_queries)
- Session Collected (proxysql.gtids.GTID_session_collected)

### ProxySQL Access Errors (proxysql.access_errors)

- Access Denied Wrong Password (proxysql.access_errors.Access_Denied_Wrong_Password)
- Access Denied Max Connections (proxysql.access_errors.Access_Denied_Max_Connections)
- Access Denied Max User Connections (proxysql.access_errors.Access_Denied_Max_User_Connections)

### ProxySQL Threads (proxysql.threads)

- Thread (proxysql.threads.MySQL_Thread_Workers)
- Monitor (proxysql.threads.MySQL_Monitor_Workers)

### ProxySQL Connection Pool Summary (proxysql.connpool_summary)

- Get Immediate (proxysql.connpool_summary.ConnPool_get_conn_immediate)
- Get Success (proxysql.connpool_summary.ConnPool_get_conn_success)
- Get Failure (proxysql.connpool_summary.ConnPool_get_conn_failure)

### ProxySQL Statement (proxysql.statement)

- Client Active Total (proxysql.statement.Stmt_Client_Active_Total)
- Client Active Unique (proxysql.statement.Stmt_Client_Active_Unique)
- Server Active Total (proxysql.statement.Stmt_Server_Active_Total)
- Server Active Unique (proxysql.statement.Stmt_Server_Active_Unique)
- Cached (proxysql.statement.Stmt_Cached)

### ProxySQL Query Cache Count (proxysql.querycache_count)

- Get (proxysql.querycache_count.Query_Cache_count_GET)
- Get OK (proxysql.querycache_count.Query_Cache_count_GET_OK)
- Set (proxysql.querycache_count.Query_Cache_count_SET)

### ProxySQL Query Cache Bytes (proxysql.querycache_bytes)

- IN (proxysql.querycache_bytes.Query_Cache_bytes_IN)
- OUT (proxysql.querycache_bytes.Query_Cache_bytes_OUT)

### ProxySQL Query Cache Entries/Purged (proxysql.querycache_entries)

- Entries (proxysql.querycache_bytes.Query_Cache_Entries)
- Purged (proxysql.querycache_bytes.Query_Cache_Purged)


### ProxySQL Monitor Workers (proxysql.monitor.workers) (Optional)

- workers (proxysql.monitor.workers.MySQL_Monitor_Workers)

### ProxySQL Monitor Checks (proxysql.monitor.checks) (Optional)

- Ping OK (proxysql.monitor.checks.MySQL_Monitor_ping_check_OK)
- Ping Error (proxysql.monitor.checks.MySQL_Monitor_ping_check_ERR)
- Readonly check OK (proxysql.monitor.checks.MySQL_Monitor_read_only_check_OK)
- Readonly check Error (proxysql.monitor.checks.MySQL_Monitor_read_only_check_ERR)
- Replication lag check OK (proxysql.monitor.checks.MySQL_Monitor_replication_lag_check_OK)
- Replication lag check Error (proxysql.monitor.checks.MySQL_Monitor_replication_lag_check_ERR)

### ProxySQL ConnectionPool Connections (proxysql.connpool.conns.#) (Optional)

- Used (proxysql.connpool.conns.#.ConnUsed)
- Free (proxysql.connpool.conns.#.ConnFree)

### ProxySQL ConnectionPool Stats (proxysql.connpool.stats.#) (Optional)

- OK (proxysql.connpool.stats.#.ConnOK)
- Error (proxysql.connpool.stats.#.ConnErr)

### ProxySQL ConnectionPool Transfers (proxysql.connpool.transfers.#) (Optional)

- Sent (proxysql.connpool.transfers.#.Bytes_data_sent)
- Received (proxysql.connpool.transfers.#.Bytes_data_recv)

### ProxySQL ConnectionPool Queries (proxysql.connpool.queries.#) (Optional)

- Queries (proxysql.connpool.queries.#.Queries)

### ProxySQL ConnectionPool Performance (proxysql.connpool.perf.#) (Optional)

- Latency sec (proxysql.connpool.perf.#.Latency_sec)

## Build

This repository is using Go 1.11 and Go modules function.
