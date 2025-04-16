package scheme

// Scheme is representation of the schemes that correspond
// to Loader types
type Scheme string

const (
	// SchemeKubernetes represents the scheme used to
	// signal that a Loader should load from Kubernetes
	SchemeKubernetes = "kube"

	// SchemeGit represents the scheme used to signal
	// that a Loader should load from a git repository
	SchemeGit = "git"

	// SchemeFile represents the scheme used to signal
	// that a Loader should load from a file
	SchemeFile = "file"
)
