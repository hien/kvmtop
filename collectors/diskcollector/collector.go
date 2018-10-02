package diskcollector

import (
	"strings"

	"github.com/cha87de/kvmtop/collectors"
	"github.com/cha87de/kvmtop/models"
)

// Collector describes the disk collector
type Collector struct {
	models.Collector
}

// Lookup disk collector data
func (collector *Collector) Lookup() {
	hostDiskSources := ""

	models.Collection.Domains.Map.Range(func(key, value interface{}) bool {
		uuid := key.(string)
		domain := value.(models.Domain)
		libvirtDomain, _ := models.Collection.LibvirtDomains.Load(uuid)

		diskLookup(&domain, libvirtDomain)
		// merge sourcedir metrics from domains to one metric for host
		disksources := strings.Split(collectors.GetMetricString(domain.Measurable, "disk_sources", 0), ",")
		for _, disksource := range disksources {
			if !strings.Contains(hostDiskSources, disksource) {
				if hostDiskSources != "" {
					hostDiskSources += ","
				}
				hostDiskSources += disksource
			}
		}

		return true
	})

	models.Collection.Host.AddMetricMeasurement("disk_sources", models.CreateMeasurement(hostDiskSources))

	diskHostLookup(models.Collection.Host)
}

// Collect disk collector data
func (collector *Collector) Collect() {
	// lookup for each domain
	models.Collection.Domains.Map.Range(func(key, value interface{}) bool {
		// uuid := key.(string)
		domain := value.(models.Domain)
		diskCollect(&domain)
		return true
	})
	diskHostCollect(models.Collection.Host)
}

// Print returns the collectors measurements in a Printable struct
func (collector *Collector) Print() models.Printable {
	printable := models.Printable{
		HostFields: []string{
			"disk_device_reads",
			"disk_device_readsmerged",
			"disk_device_sectorsread",
			"disk_device_timereading",
			"disk_device_writes",
			"disk_device_writesmerged",
			"disk_device_sectorswritten",
			"disk_device_timewriting",
			"disk_device_currentops",
			"disk_device_timeforops",
			"disk_device_weightedtimeforops",
		},
		DomainFields: []string{
			"disk_size_capacity",
			"disk_size_allocation",
			"disk_size_physical",
			"disk_stats_errs",
			"disk_stats_flushreq",
			"disk_stats_flushtotaltimes",
			"disk_stats_rdbytes",
			"disk_stats_rdreq",
			"disk_stats_rdtotaltimes",
			"disk_stats_wrbytes",
			"disk_stats_wrreq",
			"disk_stats_wrtotaltimes",
			"disk_delayblkio",
		},
	}

	// lookup for each domain
	printable.DomainValues = make(map[string][]string)
	models.Collection.Domains.Map.Range(func(key, value interface{}) bool {
		uuid := key.(string)
		domain := value.(models.Domain)
		printable.DomainValues[uuid] = diskPrint(&domain)
		return true
	})

	// lookup for host
	printable.HostValues = diskPrintHost(models.Collection.Host)

	return printable
}

// CreateCollector creates a new disk collector
func CreateCollector() Collector {
	return Collector{}
}
