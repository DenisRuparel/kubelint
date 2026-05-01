package kubelint

#PathType: "Prefix" | "Exact" | "ImplementationSpecific"

ingress: {
  name: string & != ""
  namespace: string & != ""

  labels?: {
    [string]: string
  }

  annotations?: {
    [string]: string
  }

  ingressClassName?: string & != ""

  // 🔒 SIMPLE TLS (manual only)
  tls?: [...{
    hosts: [...string & != ""]
    secretName: string & != ""
  }]

  // 🌐 Rules (required)
  rules: [...{
    host: string & != ""

    http: {
      paths: [...{
        path: string & != ""
        pathType: #PathType

        backend: {
          service: {
            name: string & != ""
            port: {
              number: int & >=1 & <=65535
            }
          }
        }
      }]
    }
  }]
}

// 🔥 If TLS is used → require valid tlsSecret
if ingress.tls != _|_ {
  tlsSecret: {
    name: string & != ""
    namespace: string & != ""

    data: {
      "tls.crt": string & != "" & !~"<base64-cert>"
      "tls.key": string & != "" & !~"<base64-key>"
    }
  }
}