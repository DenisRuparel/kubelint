package kubelint

configMap: {
  name: string & != ""
  namespace: string & != ""

  labels?: {
    [string]: string
  }

  // 🔥 Core: arbitrary key-value pairs (all strings)
  data: {
    [string]: string
  }
}