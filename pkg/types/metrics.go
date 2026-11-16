package types

import "time"

// CPURawTicks holds raw jiffies / clock ticks read directly from /proc/stat.
// In Linux, these ticks increment monotonically at USER_HZ (typically 100Hz = 10ms per tick).
type CPURawTicks struct {
	User      uint64 // Normal processes executing in user mode
	Nice      uint64 // Niced processes executing in user mode
	System    uint64 // Processes executing in kernel mode
	Idle      uint64 // Twiddling thumbs
	IOWait    uint64 // Waiting for I/O to complete
	IRQ       uint64 // Servicing interrupts
	SoftIRQ   uint64 // Servicing softirqs
	Steal     uint64 // Involuntary wait stolen by hypervisor
	Guest     uint64 // Running a normal guest OS virtual CPU
	GuestNice uint64 // Running a niced guest OS virtual CPU
}

// TotalTicks calculates the grand total of all CPU ticks accounted for.
// Note: Guest and GuestNice are already included in User and Nice respectively
// in modern Linux kernels, so they are not double-counted.
func (c CPURawTicks) TotalTicks() uint64 {
	return c.User + c.Nice + c.System + c.Idle + c.IOWait + c.IRQ + c.SoftIRQ + c.Steal
}

// IdleTicks returns the sum of idle time and iowait time.
func (c CPURawTicks) IdleTicks() uint64 {
	return c.Idle + c.IOWait
}

// ActiveTicks returns total ticks spent doing real work.
func (c CPURawTicks) ActiveTicks() uint64 {
	return c.TotalTicks() - c.IdleTicks()
}

// CPUStats represents calculated CPU usage metrics over a time delta.
type CPUStats struct {
	CoreID        int     // -1 for aggregated/overall CPU, 0..N-1 for specific cores
	UsagePercent  float64 // Total active CPU percentage [0.0, 100.0]
	UserPercent   float64 // User space percentage
	SystemPercent float64 // Kernel space percentage
	IOWaitPercent float64 // IO wait percentage
}

// MemStats contains memory metrics parsed from /proc/meminfo, stored in bytes.
type MemStats struct {
	Total       uint64  // Total usable RAM (MemTotal)
	Available   uint64  // Memory available for starting new applications without swapping (MemAvailable)
	Used        uint64  // Actually used memory (MemTotal - MemAvailable)
	Free        uint64  // Unused memory (MemFree)
	Buffers     uint64  // Relatively temporary storage for raw disk blocks (Buffers)
	Cached      uint64  // In-memory cache for files read from the disk (Cached)
	UsedPercent float64 // (Used / Total) * 100.0
}

// SwapStats contains swap space metrics parsed from /proc/meminfo, stored in bytes.
type SwapStats struct {
	Total       uint64  // Total swap space available (SwapTotal)
	Free        uint64  // Unused swap space (SwapFree)
	Used        uint64  // Used swap space (SwapTotal - SwapFree)
	UsedPercent float64 // (Used / Total) * 100.0
}

// SystemSnapshot packages hardware telemetry at a specific sampling instant.
type SystemSnapshot struct {
	Timestamp  time.Time
	Uptime     time.Duration
	OverallCPU CPUStats
	Cores      []CPUStats
	Memory     MemStats
	Swap       SwapStats
}
