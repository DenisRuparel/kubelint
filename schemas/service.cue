package kubelint

#ServiceType: "ClusterIP" | "NodePort" | "LoadBalancer"

service: {
  name: string & != ""
  namespace: string & != ""

  labels: {
    [string]: string
  }

  type: #ServiceType

  selector: {
    [string]: string
  }

  ports: [...{
    name?: string
    port: int & >=1 & <=65535
    targetPort: int | string
    protocol?: "TCP" | "UDP" | "SCTP"

    // 🔥 CONDITIONAL INSIDE FIELD (THIS FIXES EVERYTHING)
    nodePort?: int & >=30000 & <=32767

    if service.type == "NodePort" {
      nodePort: int & >=30000 & <=32767
    }
  }]

  sessionAffinity?: "None" | "ClientIP"
  internalTrafficPolicy?: "Cluster" | "Local"
}