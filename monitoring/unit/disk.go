package monitoring

import (
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskInfo struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
}

func Disk() DiskInfo {
	diskinfo := DiskInfo{}
	usage, err := disk.Partitions(false) // 使用 false 只获取物理分区
	if err != nil {
		diskinfo.Total = 0
		diskinfo.Used = 0
	} else {
		for _, part := range usage {
			// 排除临时文件系统和网络驱动器
			if isPhysicalDisk(part) {
				u, err := disk.Usage(part.Mountpoint)
				if err != nil {
					continue
				} else {
					diskinfo.Total += u.Total
					diskinfo.Used += u.Used
				}
			}
		}
	}
	return diskinfo
}

// isPhysicalDisk 判断分区是否为物理磁盘
func isPhysicalDisk(part disk.PartitionStat) bool {
    mountpoint := strings.ToLower(part.Mountpoint)
    fstype := strings.ToLower(part.Fstype)

    // 强制包含根分区
    // if mountpoint == "/" {
    //     return true
    // }

    // 排除临时文件系统
    if mountpoint == "/tmp" || mountpoint == "/var/tmp" || mountpoint == "/dev/shm" ||
        mountpoint == "/run" || mountpoint == "/run/lock" {
        return false
    }

    // 排除 Docker 和 k3s 挂载点
    if strings.Contains(mountpoint, "/run/k3s") || strings.Contains(mountpoint, "/var/lib/docker") {
        return false
    }

    // 排除网络文件系统
    if strings.HasPrefix(fstype, "nfs") || strings.HasPrefix(fstype, "cifs") ||
        strings.HasPrefix(fstype, "smb") || fstype == "vboxsf" || fstype == "9p" ||
        strings.Contains(fstype, "fuse") {
        return false
    }

    // 排除 overlay 文件系统
    if fstype == "overlay" {
        return false
    }

    // 排除虚拟内存和 loop 设备
    if strings.HasPrefix(part.Device, "/dev/loop") || fstype == "devtmpfs" || fstype == "tmpfs" {
        return false
    }

    // 排除网络相关选项
    optsStr := strings.ToLower(strings.Join(part.Opts, ","))
    if strings.Contains(optsStr, "remote") || strings.Contains(optsStr, "network") {
        return false
    }

    return true
}
