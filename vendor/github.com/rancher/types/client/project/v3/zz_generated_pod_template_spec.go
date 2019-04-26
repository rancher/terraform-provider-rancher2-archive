package client

const (
	PodTemplateSpecType                               = "podTemplateSpec"
	PodTemplateSpecFieldActiveDeadlineSeconds         = "activeDeadlineSeconds"
	PodTemplateSpecFieldAutomountServiceAccountToken  = "automountServiceAccountToken"
	PodTemplateSpecFieldContainers                    = "containers"
	PodTemplateSpecFieldDNSConfig                     = "dnsConfig"
	PodTemplateSpecFieldDNSPolicy                     = "dnsPolicy"
	PodTemplateSpecFieldFsgid                         = "fsgid"
	PodTemplateSpecFieldGids                          = "gids"
	PodTemplateSpecFieldHostAliases                   = "hostAliases"
	PodTemplateSpecFieldHostIPC                       = "hostIPC"
	PodTemplateSpecFieldHostNetwork                   = "hostNetwork"
	PodTemplateSpecFieldHostPID                       = "hostPID"
	PodTemplateSpecFieldHostname                      = "hostname"
	PodTemplateSpecFieldImagePullSecrets              = "imagePullSecrets"
	PodTemplateSpecFieldNodeID                        = "nodeId"
	PodTemplateSpecFieldObjectMeta                    = "metadata"
	PodTemplateSpecFieldPriority                      = "priority"
	PodTemplateSpecFieldPriorityClassName             = "priorityClassName"
	PodTemplateSpecFieldReadinessGates                = "readinessGates"
	PodTemplateSpecFieldRestartPolicy                 = "restartPolicy"
	PodTemplateSpecFieldRunAsGroup                    = "runAsGroup"
	PodTemplateSpecFieldRunAsNonRoot                  = "runAsNonRoot"
	PodTemplateSpecFieldRuntimeClassName              = "runtimeClassName"
	PodTemplateSpecFieldSchedulerName                 = "schedulerName"
	PodTemplateSpecFieldScheduling                    = "scheduling"
	PodTemplateSpecFieldServiceAccountName            = "serviceAccountName"
	PodTemplateSpecFieldShareProcessNamespace         = "shareProcessNamespace"
	PodTemplateSpecFieldSubdomain                     = "subdomain"
	PodTemplateSpecFieldSysctls                       = "sysctls"
	PodTemplateSpecFieldTerminationGracePeriodSeconds = "terminationGracePeriodSeconds"
	PodTemplateSpecFieldUid                           = "uid"
	PodTemplateSpecFieldVolumes                       = "volumes"
)

type PodTemplateSpec struct {
	ActiveDeadlineSeconds         *int64                 `json:"activeDeadlineSeconds,omitempty" yaml:"activeDeadlineSeconds,omitempty"`
	AutomountServiceAccountToken  *bool                  `json:"automountServiceAccountToken,omitempty" yaml:"automountServiceAccountToken,omitempty"`
	Containers                    []Container            `json:"containers,omitempty" yaml:"containers,omitempty"`
	DNSConfig                     *PodDNSConfig          `json:"dnsConfig,omitempty" yaml:"dnsConfig,omitempty"`
	DNSPolicy                     string                 `json:"dnsPolicy,omitempty" yaml:"dnsPolicy,omitempty"`
	Fsgid                         *int64                 `json:"fsgid,omitempty" yaml:"fsgid,omitempty"`
	Gids                          []int64                `json:"gids,omitempty" yaml:"gids,omitempty"`
	HostAliases                   []HostAlias            `json:"hostAliases,omitempty" yaml:"hostAliases,omitempty"`
	HostIPC                       bool                   `json:"hostIPC,omitempty" yaml:"hostIPC,omitempty"`
	HostNetwork                   bool                   `json:"hostNetwork,omitempty" yaml:"hostNetwork,omitempty"`
	HostPID                       bool                   `json:"hostPID,omitempty" yaml:"hostPID,omitempty"`
	Hostname                      string                 `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	ImagePullSecrets              []LocalObjectReference `json:"imagePullSecrets,omitempty" yaml:"imagePullSecrets,omitempty"`
	NodeID                        string                 `json:"nodeId,omitempty" yaml:"nodeId,omitempty"`
	ObjectMeta                    *ObjectMeta            `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Priority                      *int64                 `json:"priority,omitempty" yaml:"priority,omitempty"`
	PriorityClassName             string                 `json:"priorityClassName,omitempty" yaml:"priorityClassName,omitempty"`
	ReadinessGates                []PodReadinessGate     `json:"readinessGates,omitempty" yaml:"readinessGates,omitempty"`
	RestartPolicy                 string                 `json:"restartPolicy,omitempty" yaml:"restartPolicy,omitempty"`
	RunAsGroup                    *int64                 `json:"runAsGroup,omitempty" yaml:"runAsGroup,omitempty"`
	RunAsNonRoot                  *bool                  `json:"runAsNonRoot,omitempty" yaml:"runAsNonRoot,omitempty"`
	RuntimeClassName              string                 `json:"runtimeClassName,omitempty" yaml:"runtimeClassName,omitempty"`
	SchedulerName                 string                 `json:"schedulerName,omitempty" yaml:"schedulerName,omitempty"`
	Scheduling                    *Scheduling            `json:"scheduling,omitempty" yaml:"scheduling,omitempty"`
	ServiceAccountName            string                 `json:"serviceAccountName,omitempty" yaml:"serviceAccountName,omitempty"`
	ShareProcessNamespace         *bool                  `json:"shareProcessNamespace,omitempty" yaml:"shareProcessNamespace,omitempty"`
	Subdomain                     string                 `json:"subdomain,omitempty" yaml:"subdomain,omitempty"`
	Sysctls                       []Sysctl               `json:"sysctls,omitempty" yaml:"sysctls,omitempty"`
	TerminationGracePeriodSeconds *int64                 `json:"terminationGracePeriodSeconds,omitempty" yaml:"terminationGracePeriodSeconds,omitempty"`
	Uid                           *int64                 `json:"uid,omitempty" yaml:"uid,omitempty"`
	Volumes                       []Volume               `json:"volumes,omitempty" yaml:"volumes,omitempty"`
}
