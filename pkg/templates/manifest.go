package templates

// Manifest is the template found in examples
var Manifest = NewTemplate("manifest", `
---
apiVersion: v1
kind: Namespace
metadata:
  name: heptio-sonobuoy
---
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    component: sonobuoy
  name: sonobuoy-serviceaccount
  namespace: heptio-sonobuoy
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  labels:
    component: sonobuoy
  name: sonobuoy-serviceaccount
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: sonobuoy-serviceaccount
subjects:
- kind: ServiceAccount
  name: sonobuoy-serviceaccount
  namespace: heptio-sonobuoy
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  labels:
    component: sonobuoy
  name: sonobuoy-serviceaccount
  namespace: heptio-sonobuoy
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
---
apiVersion: v1
data:
  config.json: |
    {
        "Description": "EXAMPLE",
        "Filters": {
            "LabelSelector": "",
            "Namespaces": ".*"
        },
        "PluginNamespace": "heptio-sonobuoy",
        "Plugins": {{.PluginSelector}},
        "Resources": [
            "CertificateSigningRequests",
            "ClusterRoleBindings",
            "ClusterRoles",
            "ComponentStatuses",
            "CustomResourceDefinitions",
            "Nodes",
            "PersistentVolumes",
            "PodSecurityPolicies",
            "ServerVersion",
            "StorageClasses",
            "ConfigMaps",
            "DaemonSets",
            "Deployments",
            "Endpoints",
            "Events",
            "HorizontalPodAutoscalers",
            "Ingresses",
            "Jobs",
            "LimitRanges",
            "PersistentVolumeClaims",
            "Pods",
            "PodDisruptionBudgets",
            "PodTemplates",
            "ReplicaSets",
            "ReplicationControllers",
            "ResourceQuotas",
            "RoleBindings",
            "Roles",
            "ServerGroups",
            "ServiceAccounts",
            "Services",
            "StatefulSets"
        ],
        "ResultsDir": "/tmp/sonobuoy",
        "Server": {
            "advertiseaddress": "sonobuoy-master:8080",
            "bindaddress": "0.0.0.0",
            "bindport": 8080,
            "timeoutseconds": 5400
        },
        "Version": "0.11.0",{{/* TODO(EKF): Should be buildinfo.Version, when that's set*/}}
        "WorkerImage": "{{.SonobuoyImage}}"
    }
kind: ConfigMap
metadata:
  labels:
    component: sonobuoy
  name: sonobuoy-config-cm
  namespace: heptio-sonobuoy
---
apiVersion: v1
data:
  e2e.yaml: |
    sonobuoy-config:
      driver: Job
      plugin-name: e2e
      result-type: e2e
    spec:
      env:
      - name: E2E_FOCUS
        value: {{.E2EFocus}}
      image: gcr.io/heptio-images/kube-conformance:latest
      imagePullPolicy: Always
      name: e2e
      volumeMounts:
      - mountPath: /tmp/results
        name: results
        readOnly: false
  systemd-logs.yaml: |
    sonobuoy-config:
      driver: DaemonSet
      plugin-name: systemd-logs
      result-type: systemd_logs
    spec:
      command:
      - sh
      - -c
      - /get_systemd_logs.sh && sleep 3600
      env:
      - name: NODE_NAME
        valueFrom:
          fieldRef:
            fieldPath: spec.nodeName
      - name: RESULTS_DIR
        value: /tmp/results
      - name: CHROOT_DIR
        value: /node
      image: gcr.io/heptio-images/sonobuoy-plugin-systemd-logs:latest
      imagePullPolicy: Always
      name: sonobuoy-systemd-logs-config
      securityContext:
        privileged: true
      volumeMounts:
      - mountPath: /tmp/results
        name: results
        readOnly: false
      - mountPath: /node
        name: root
        readOnly: false
kind: ConfigMap
metadata:
  labels:
    component: sonobuoy
  name: sonobuoy-plugins-cm
  namespace: heptio-sonobuoy
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    component: sonobuoy
    run: sonobuoy-master
    tier: analysis
  name: sonobuoy
  namespace: heptio-sonobuoy
spec:
  containers:
  - command:
    - /bin/bash
    - -c
    - /sonobuoy master --no-exit=true -v 3 --logtostderr
    env:
    - name: SONOBUOY_ADVERTISE_IP
      valueFrom:
        fieldRef:
          fieldPath: status.podIP
    image: {{.SonobuoyImage}}
    imagePullPolicy: Always
    name: kube-sonobuoy
    volumeMounts:
    - mountPath: /etc/sonobuoy
      name: sonobuoy-config-volume
    - mountPath: /plugins.d
      name: sonobuoy-plugins-volume
    - mountPath: /tmp/sonobuoy
      name: output-volume
  restartPolicy: Never
  serviceAccountName: sonobuoy-serviceaccount
  volumes:
  - configMap:
      name: sonobuoy-config-cm
    name: sonobuoy-config-volume
  - configMap:
      name: sonobuoy-plugins-cm
    name: sonobuoy-plugins-volume
  - emptyDir: {}
    name: output-volume
---
apiVersion: v1
kind: Service
metadata:
  labels:
    component: sonobuoy
    run: sonobuoy-master
  name: sonobuoy-master
  namespace: heptio-sonobuoy
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    run: sonobuoy-master
  type: ClusterIP
`)
