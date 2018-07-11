package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseKvpDynoLoadMetrics(t *testing.T) {
	msg := "source=web.4 dyno=heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b sample#load_avg_1m=0.27 sample#load_avg_5m=0.22 sample#load_avg_15m=0.19"

	metrics, err := parseKvp(msg)

	assert.Nil(t, err)

	assert.Equal(
		t,
		map[string]interface{}{
			"source":       "web.4",
			"dyno":         "heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b",
			"load_avg_1m":  0.27,
			"load_avg_5m":  0.22,
			"load_avg_15m": 0.19,
		},
		metrics,
	)
}

func TestParseKvpDynoMemoryMetrics(t *testing.T) {
	msg := "source=web.4 dyno=heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b sample#memory_total=7020.55MB sample#memory_rss=7013.48MB sample#memory_cache=7.08MB sample#memory_swap=0.00MB sample#memory_pgpgin=1949515pages sample#memory_pgpgout=779250pages sample#memory_quota=14336.00MB"

	metrics, err := parseKvp(msg)

	assert.Nil(t, err)

	assert.Equal(
		t,
		map[string]interface{}{
			"source":          "web.4",
			"dyno":            "heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b",
			"memory_total_MB": 7020.55,
			"memory_rss_MB":   7013.48,
			"memory_cache_MB": 7.08,
			"memory_swap_MB":  0.0,
			"memory_pgpgin":   1949515.0,
			"memory_pgpgout":  779250.0,
			"memory_quota_MB": 14336.0,
		},
		metrics,
	)
}

func TestParseKvpPostgresMetrics(t *testing.T) {
	msg := "source=HEROKU_POSTGRESQL_ONYX sample#current_transaction=153674642 sample#db_size=53666224664bytes sample#tables=105 sample#active-connections=20 sample#waiting-connections=0 sample#index-cache-hit-rate=0.99832 sample#table-cache-hit-rate=0.9657 sample#load-avg-1m=0.01 sample#load-avg-5m=0.045 sample#load-avg-15m=0.025 sample#read-iops=0.58264 sample#write-iops=3.8988 sample#memory-total=8173656kB sample#memory-free=1124196kB sample#memory-cached=5935596kB sample#memory-postgres=111996kB"

	metrics, err := parseKvp(msg)

	assert.Nil(t, err)

	assert.Equal(
		t,
		map[string]interface{}{
			"source":               "HEROKU_POSTGRESQL_ONYX",
			"current_transaction":  153674642,
			"db_size_bytes":        53666224664.0,
			"db_size_MB":           53666224664.0 / 1024 / 1024,
			"tables":               105,
			"active_connections":   20,
			"waiting_connections":  0,
			"index_cache_hit_rate": 0.99832,
			"table_cache_hit_rate": 0.9657,
			"load_avg_1m":          0.01,
			"load_avg_5m":          0.045,
			"load_avg_15m":         0.025,
			"read_iops":            0.58264,
			"write_iops":           3.8988,
			"memory_total_kB":      8173656.0,
			"memory_total_MB":      8173656.0 / 1024,
			"memory_free_kB":       1124196.0,
			"memory_free_MB":       1124196.0 / 1024,
			"memory_cached_kB":     5935596.0,
			"memory_cached_MB":     5935596.0 / 1024,
			"memory_postgres_kB":   111996.0,
			"memory_postgres_MB":   111996.0 / 1024,
		},
		metrics,
	)
}

func TestParseKvpRedisMetrics(t *testing.T) {
	msg := "source=HEROKU_REDIS_BLUE sample#active-connections=15 sample#load-avg-1m=0.16 sample#load-avg-5m=0.1 sample#load-avg-15m=0.075 sample#read-iops=0 sample#write-iops=0.034583 sample#memory-total=15664360kB sample#memory-free=11254388kB sample#memory-cached=2726820kB sample#memory-redis=2087784bytes sample#hit-rate=0.094117 sample#evicted-keys=0"

	metrics, err := parseKvp(msg)

	assert.Nil(t, err)

	assert.Equal(
		t,
		map[string]interface{}{
			"source":             "HEROKU_REDIS_BLUE",
			"active_connections": 15,
			"load_avg_1m":        0.16,
			"load_avg_5m":        0.1,
			"load_avg_15m":       0.075,
			"read_iops":          0,
			"write_iops":         0.034583,
			"memory_total_kB":    15664360.0,
			"memory_total_MB":    15664360.0 / 1024,
			"memory_free_kB":     11254388.0,
			"memory_free_MB":     11254388.0 / 1024,
			"memory_cached_kB":   2726820.0,
			"memory_cached_MB":   2726820.0 / 1024,
			"memory_redis_bytes": 2087784.0,
			"memory_redis_MB":    2087784.0 / 1024 / 1024,
			"hit_rate":           0.094117,
			"evicted_keys":       0,
		},
		metrics,
	)
}

func TestParseKvpHerokuError(t *testing.T) {
	msg := `at=error code=H11 desc="Backlog too deep" method=GET path="/?k=p" host=myapp.herokuapp.com fwd=17.17.17.17 dyno= connect= service= status=503 bytes=`

	metrics, err := parseKvp(msg)

	assert.Nil(t, err)

	assert.Equal(
		t,
		map[string]interface{}{
			"at":      "error",
			"code":    "H11",
			"desc":    "Backlog too deep",
			"method":  "GET",
			"path":    "/?k=p",
			"host":    "myapp.herokuapp.com",
			"fwd":     "17.17.17.17",
			"dyno":    nil,
			"connect": nil,
			"service": nil,
			"status":  503,
			"bytes":   nil,
		},
		metrics,
	)
}

func TestParseKvpInvalid1(t *testing.T) {
	msg := `looks=like a=kvp but=really NOTTTTTT sorta=notreally`

	_, err := parseKvp(msg)

	assert.NotNil(t, err)
}
