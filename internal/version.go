package internal

import (
	"bytes"
	"fmt"
)

// Program version information
var (
	buildArch     string
	buildBranch   string
	buildCompiler string
	buildDate     string
	buildHash     string
	buildUser     string
	buildVersion  string
)

// VersionInfo for client consumption
type VersionInfo struct {
	Arch     string
	Branch   string
	Compiler string
	Date     string
	Hash     string
	User     string
	Version  string
}

func (vi VersionInfo) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Version:\t\t%s\n", vi.Version))
	buffer.WriteString(fmt.Sprintf("Build time:\t\t%s\n", vi.Date))
	buffer.WriteString(fmt.Sprintf("Build:\t\t\t%s@%s:%s\n", vi.User, vi.Branch, vi.Hash))
	buffer.WriteString(fmt.Sprintf("Compiler:\t\t%s\n", vi.Compiler))
	buffer.WriteString(fmt.Sprintf("Arch:\t\t\t%s\n", vi.Arch))
	return buffer.String()
}

// GetVersionInfo to set expectations
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Arch:     buildArch,
		Branch:   buildBranch,
		Compiler: buildCompiler,
		Hash:     buildHash,
		Date:     buildDate,
		User:     buildUser,
		Version:  buildVersion,
	}
}

