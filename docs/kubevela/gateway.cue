package main

gateway: {
	description: "Enable public web traffic for the component, the ingress API matches K8s v1.20+."
	type:        "trait"
	attributes: {
		podDisruptive: false
		appliesToWorkloads: ["*"]
	}
}
template: {
	parameter: {
		domain: string
		http: [string]: int
		class:       *"traefik" | string
		secretName?: string
		secretNamespace?: string
		listenerPort: *80 | int
		annotations?: [string]: string
		labels?: [string]: string
	}
	// Create Service
	outputs: {
		service: {
			apiVersion: "v1"
			kind:       "Service"
			metadata: {
				name:      context.name
				namespace: context.namespace
				if parameter.labels != _|_ {
					labels: parameter.labels
				}
				if parameter.annotations != _|_ {
					annotations: parameter.annotations
				}
			}
			spec: {
				selector: {
					"app.oam.dev/component": context.name
				}
				ports: [
					for k, v in parameter.http {
						name:       "http"
						port:       v
						targetPort: v
						protocol:   "TCP"
					},
				]
				type: "ClusterIP"
			}
		}
		// Create Ingress
		ingress: {
			apiVersion: "networking.k8s.io/v1"
			kind:       "Ingress"
			metadata: {
				name:      context.name
				namespace: context.namespace
				if parameter.labels != _|_ {
					labels: parameter.labels
				}
				if parameter.annotations != _|_ {
					annotations: parameter.annotations
				}
			}
			spec: {
				ingressClassName: parameter.class
				if parameter.secretName != _|_ {
					tls: [{
						hosts: [parameter.domain]
						secretName: parameter.secretName
					}]
				}
				rules: [{
					host: parameter.domain
					http: {
						paths: [
							for k, v in parameter.http {
								path:     k
								pathType: "Prefix"
								backend: {
									service: {
										name: context.name
										port: number: v
									}
								}
							},
						]
					}
				}]
			}
		}
	}
}