package kubelint

deployment: {

  name: string & != ""
  namespace: string & != ""

  labels: {
    [string]: string
  }

  replicas: int & >0

  revisionHistoryLimit?: int & >=0
  progressDeadlineSeconds?: int & >0

  strategy?: {
    type: "RollingUpdate" | "Recreate"

    rollingUpdate?: {
      maxUnavailable?: int | string
      maxSurge?: int | string
    }
  }

  serviceAccountName?: string
  automountServiceAccountToken?: bool

  terminationGracePeriodSeconds?: int & >=0

  securityContext?: {
    runAsNonRoot?: bool
    runAsUser?: int & >=0
    runAsGroup?: int & >=0
    fsGroup?: int & >=0
  }

  containers: [...{
    name: string & != ""

    image: {
      repository: string & != ""
      tag: string & != ""
    }

    imagePullPolicy?: "Always" | "IfNotPresent" | "Never"

    ports?: [...{
      containerPort: int & >=1 & <=65535
      name?: string
    }]

    resources?: {
      requests?: {
        cpu?: string
        memory?: string
      }
      limits?: {
        cpu?: string
        memory?: string
      }
    }

    securityContext?: {
      privileged?: bool
      allowPrivilegeEscalation?: bool
      readOnlyRootFilesystem?: bool
    }

    envFrom?: [...{
      configMapRef?: {
        name: string
      }
      secretRef?: {
        name: string
      }
    }]

    livenessProbe?: probe
    readinessProbe?: probe
    startupProbe?: probe

    volumeMounts?: [...{
      name: string
      mountPath: string
    }]
  }]

  volumes?: [...{
    name: string

    emptyDir?: {
      medium?: "Memory"
    }
  }]

  imagePullSecrets?: [...{
    name: string
  }]

  nodeSelector?: {
    [string]: string
  }

  tolerations?: [...{
    key?: string
    operator?: "Exists" | "Equal"
    value?: string
    effect?: "NoSchedule" | "PreferNoSchedule" | "NoExecute"
  }]

  dnsPolicy?: "ClusterFirst" | "Default"
  restartPolicy?: "Always" | "OnFailure" | "Never"
}

probe: {
  httpGet?: {
    path: string
    port: string | int
  }

  initialDelaySeconds?: int & >=0
  periodSeconds?: int & >0
  timeoutSeconds?: int & >0
  failureThreshold?: int & >0
}