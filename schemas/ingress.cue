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

  // 🔒 TLS (optional)
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