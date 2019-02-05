package mpproxysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var proxysql ProxySQLPlugin

	graphdef := proxysql.GraphDefinition()

	assert.Equal(t, len(graphdef), 15)

	assert.Equal(t, len(graphdef["transactions"].Metrics), 1)
	assert.Equal(t, len(graphdef["connections"].Metrics), 8)
	assert.Equal(t, len(graphdef["traffic"].Metrics), 4)
	assert.Equal(t, len(graphdef["processor"].Metrics), 2)
	assert.Equal(t, len(graphdef["buffers"].Metrics), 6)
	assert.Equal(t, len(graphdef["commands"].Metrics), 19)
	assert.Equal(t, len(graphdef["slow_queries"].Metrics), 1)
	assert.Equal(t, len(graphdef["gtids"].Metrics), 2)
	assert.Equal(t, len(graphdef["access_errors"].Metrics), 3)
	assert.Equal(t, len(graphdef["threads"].Metrics), 2)
	assert.Equal(t, len(graphdef["connpool_summary"].Metrics), 3)
	assert.Equal(t, len(graphdef["statement"].Metrics), 5)
	assert.Equal(t, len(graphdef["querycache_count"].Metrics), 3)
	assert.Equal(t, len(graphdef["querycache_bytes"].Metrics), 2)
	assert.Equal(t, len(graphdef["querycache_entries"].Metrics), 2)

}

func TestGraphDefinitionWithMonitorStats(t *testing.T) {
	var proxysql ProxySQLPlugin

	proxysql.EnableMonitorStats = true
	graphdef := proxysql.GraphDefinition()

	assert.Equal(t, len(graphdef), 15+2)

	// ConnectionPool addtional metrics
	assert.Equal(t, len(graphdef["monitor.workers"].Metrics), 1)
	assert.Equal(t, len(graphdef["monitor.checks"].Metrics), 6)
}

func TestGraphDefinitionWithConnectionPool(t *testing.T) {
	var proxysql ProxySQLPlugin

	proxysql.EnableConnectionPool = true
	graphdef := proxysql.GraphDefinition()

	assert.Equal(t, len(graphdef), 15+5)

	// ConnectionPool addtional metrics
	assert.Equal(t, len(graphdef["connpool.conns.#"].Metrics), 2)
	assert.Equal(t, len(graphdef["connpool.stats.#"].Metrics), 2)
	assert.Equal(t, len(graphdef["connpool.transfers.#"].Metrics), 2)
	assert.Equal(t, len(graphdef["connpool.queries.#"].Metrics), 1)
	assert.Equal(t, len(graphdef["connpool.perf.#"].Metrics), 1)
}

func TestGraphDefinitionAllStats(t *testing.T) {
	var proxysql ProxySQLPlugin

	proxysql.EnableMonitorStats = true
	proxysql.EnableConnectionPool = true
	graphdef := proxysql.GraphDefinition()

	assert.Equal(t, len(graphdef), 15+2+5)
}
