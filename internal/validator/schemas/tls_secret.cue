package kubelint

tlsSecret?: {
  name: string & != ""
  namespace: string & != ""

  labels?: {
    [string]: string
  }

  // 🔒 Fixed type
  type: "kubernetes.io/tls"

  // 🔥 Required TLS fields (base64 encoded)
  data: {
    "tls.crt": string & != "" & !="<base64-cert>"
    "tls.key": string & != "" & !="<base64-key>"
  }
}