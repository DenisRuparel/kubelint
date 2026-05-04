package kubelint

#SecretType: "Opaque" | "kubernetes.io/dockerconfigjson" | "kubernetes.io/tls"

secret: {
  name: string & != ""
  namespace: string & != ""

  labels?: {
    [string]: string
  }

  // Default to Opaque if not provided
  type: *"Opaque" | #SecretType

  // 🔥 Use stringData for human-readable input
  stringData: {
    [string]: string & != ""
  }
}