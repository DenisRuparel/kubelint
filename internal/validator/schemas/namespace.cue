package kubelint

namespace: {
  name: string & != ""

  labels?: {
    [string]: string
  }

  annotations?: {
    [string]: string
  }
}