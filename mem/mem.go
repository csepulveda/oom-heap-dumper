package mem

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

// LimitAndUsageForProc returns memory limit and usage for cgroup where proc
// is running.
func LimitAndUsageForProc(proc *os.Process) (uint64, uint64, error) {
	limit, err := LimitForProc(proc)
	if err != nil {
		return 0, 0, fmt.Errorf("error reading memory limit for pid %d: %w", proc.Pid, err)
	}
	usage, err := UsageForProc(proc)
	if err != nil {
		return 0, 0, fmt.Errorf("error reading memory usage for pid %d: %w", proc.Pid, err)
	}
	return limit, usage, nil
}

// LimitForProc returns the max memory on bytes that the process can use.
func LimitForProc(proc *os.Process) (uint64, error) {

	//search limit for cgroup v1
	limitFile := fmt.Sprintf("/proc/%d/root/sys/fs/cgroup/memory/memory.limit_in_bytes", proc.Pid)
	val, err := readBytesFromFile(limitFile)
	if err == nil {
		return val, nil
	}

	//search limit for cgroup v2
	limitFile = fmt.Sprintf("/proc/%d/root/sys/fs/cgroup/memory.max", proc.Pid)
	val, err = readBytesFromFile(limitFile)
	if err == nil {
		return val, nil
	}

	return 0, err

}

// UsageForProc returns the amount of memory currently in use.
func UsageForProc(proc *os.Process) (uint64, error) {

	//search usage for cgroup v1
	usageFile := fmt.Sprintf("/proc/%d/root/sys/fs/cgroup/memory/memory.usage_in_bytes", proc.Pid)
	val, err := readBytesFromFile(usageFile)
	if err == nil {
		return val, nil
	}

	//search usage for cgroup v2
	usageFile = fmt.Sprintf("/proc/%d/root/sys/fs/cgroup/memory.current", proc.Pid)
	val, err = readBytesFromFile(usageFile)
	if err == nil {
		return val, nil
	}
	return 0, err
}

// readBytesFromFile reads a file and returns its content as a uint64. if the string
// "max" is found, this returns 0.
func readBytesFromFile(fpath string) (uint64, error) {
	content, err := os.ReadFile(fpath)
	if err != nil {
		return 0, err
	}
	content = bytes.TrimSpace(content)
	if string(content) == "max" {
		return 0, nil
	}
	return strconv.ParseUint(string(content), 10, 64)
}
