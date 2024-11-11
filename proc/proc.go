package proc

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/csepulveda/oom-heap-dumper/mem"
)

type Process interface {
	Pid() int
	MemoryUsagePercent() (uint64, error)
	PortsInUse() ([]int, error)
}

// Others return a list of all other processes running on the system, excluding
// the current one.
func Others() ([]*os.Process, error) {
	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	ps := make([]*os.Process, 0)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}

		if pid == os.Getpid() {
			continue
		}

		proccess, err := os.FindProcess(pid)
		if err != nil {
			return nil, err
		}

		ps = append(ps, proccess)
	}

	if len(ps) == 0 {
		return nil, fmt.Errorf("unable to find any process")
	}

	return ps, nil
}

type OsProcess struct {
	process *os.Process
}

func NewOsProcess(p *os.Process) OsProcess {
	return OsProcess{
		process: p,
	}
}

func (p OsProcess) Pid() int {
	return p.process.Pid
}

func (p OsProcess) MemoryUsagePercent() (uint64, error) {
	limit, usage, err := mem.LimitAndUsageForProc(p.process)
	if err != nil {
		return 0, err
	} else if limit == 0 {
		return 0, nil
	}
	return (usage * 100) / limit, nil
}

func (p OsProcess) PortsInUse() ([]int, error) {
	// Buscar puertos en tcp (IPv4)
	ports, err := readListeningPorts(fmt.Sprintf("/proc/%d/net/tcp", p.Pid()))
	if err != nil {
		return nil, err
	}

	// Buscar puertos en tcp6 (IPv6)
	ports6, err := readListeningPorts(fmt.Sprintf("/proc/%d/net/tcp6", p.Pid()))
	if err != nil {
		return nil, err
	}

	// Combinar puertos IPv4 e IPv6
	ports = append(ports, ports6...)
	return ports, nil
}

func readListeningPorts(filepath string) ([]int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", filepath, err)
	}
	defer file.Close()

	var ports []int
	scanner := bufio.NewScanner(file)
	// Saltar la línea de encabezado
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		// Extraer la dirección local
		localAddress := fields[1]
		ipPort := strings.Split(localAddress, ":")
		if len(ipPort) < 2 {
			continue
		}

		// Convertir el puerto de hexadecimal a decimal
		portHex := ipPort[len(ipPort)-1] // Obtener el último elemento que es el puerto
		port, err := strconv.ParseInt(portHex, 16, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse port %s: %v", portHex, err)
		}

		// Verificar el estado para asegurarse de que esté en escucha (estado 0x0A es LISTEN)
		state := fields[3]
		if state == "0A" {
			ports = append(ports, int(port))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading %s: %v", filepath, err)
	}

	return ports, nil
}
