select time_bucket('1 minute', ts) as time,
	host,
	max(usage) max_cpu_usage,
	min(usage) min_cpu_usage
from cpu_usage
where ts between $2
  and $3
  and host = $1
group by time, host
