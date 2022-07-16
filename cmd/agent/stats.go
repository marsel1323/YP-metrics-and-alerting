package main

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"runtime"
	"time"
)

func SetRuntimeStats(cache *AgentCache) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	cache.Set(Alloc, NewGaugeMetric(Alloc, float64(memStats.Alloc)))
	cache.Set(BuckHashSys, NewGaugeMetric(BuckHashSys, float64(memStats.BuckHashSys)))
	cache.Set(Frees, NewGaugeMetric(Frees, float64(memStats.Frees)))
	cache.Set(GCCPUFraction, NewGaugeMetric(GCCPUFraction, memStats.GCCPUFraction))
	cache.Set(GCSys, NewGaugeMetric(GCSys, float64(memStats.GCSys)))
	cache.Set(HeapAlloc, NewGaugeMetric(HeapAlloc, float64(memStats.HeapAlloc)))
	cache.Set(HeapIdle, NewGaugeMetric(HeapIdle, float64(memStats.HeapIdle)))
	cache.Set(HeapInuse, NewGaugeMetric(HeapInuse, float64(memStats.HeapInuse)))
	cache.Set(HeapObjects, NewGaugeMetric(HeapObjects, float64(memStats.HeapObjects)))
	cache.Set(HeapReleased, NewGaugeMetric(HeapReleased, float64(memStats.HeapReleased)))
	cache.Set(HeapSys, NewGaugeMetric(HeapSys, float64(memStats.HeapSys)))
	cache.Set(LastGC, NewGaugeMetric(LastGC, float64(memStats.LastGC)))
	cache.Set(Lookups, NewGaugeMetric(Lookups, float64(memStats.Lookups)))
	cache.Set(MCacheSys, NewGaugeMetric(MCacheSys, float64(memStats.MCacheSys)))
	cache.Set(MCacheInuse, NewGaugeMetric(MCacheInuse, float64(memStats.MCacheInuse)))
	cache.Set(MSpanInuse, NewGaugeMetric(MSpanInuse, float64(memStats.MSpanInuse)))
	cache.Set(MSpanSys, NewGaugeMetric(MSpanSys, float64(memStats.MSpanSys)))
	cache.Set(Mallocs, NewGaugeMetric(Mallocs, float64(memStats.Mallocs)))
	cache.Set(NextGC, NewGaugeMetric(NextGC, float64(memStats.NextGC)))
	cache.Set(NumForcedGC, NewGaugeMetric(NumForcedGC, float64(memStats.NumForcedGC)))
	cache.Set(NumGC, NewGaugeMetric(NumGC, float64(memStats.NumGC)))
	cache.Set(OtherSys, NewGaugeMetric(OtherSys, float64(memStats.OtherSys)))
	cache.Set(PauseTotalNs, NewGaugeMetric(PauseTotalNs, float64(memStats.PauseTotalNs)))
	cache.Set(StackInuse, NewGaugeMetric(StackInuse, float64(memStats.StackInuse)))
	cache.Set(StackSys, NewGaugeMetric(StackSys, float64(memStats.StackSys)))
	cache.Set(TotalAlloc, NewGaugeMetric(TotalAlloc, float64(memStats.TotalAlloc)))
	cache.Set(Sys, NewGaugeMetric(Sys, float64(memStats.Sys)))
}

func SetVirtualMemoryStats(cache *AgentCache) {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
	}
	cache.Set(TotalMemory, NewGaugeMetric(TotalMemory, float64(v.Total)))
	cache.Set(FreeMemory, NewGaugeMetric(FreeMemory, float64(v.Free)))
}

func SetCPUStats(cache *AgentCache) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println(err)
	}
	cache.Set(CPUutilization1, NewGaugeMetric(CPUutilization1, percent[0]))
}
