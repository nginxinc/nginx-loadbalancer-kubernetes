package buildinfo

var semVer string

// SemVer is the version number of this build as provided by build pipeline
func SemVer() string {
	return semVer
}

var shortHash string

// ShortHash is the 8 char git shorthash
func ShortHash() string {
	return shortHash
}
