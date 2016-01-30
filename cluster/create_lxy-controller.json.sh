#!/usr/bin/sh

echo '{
  "kind": "ReplicationController",
  "apiVersion": "v1beta3",
  "metadata": {
    "name": "lxy-controller",
    "labels": {
      "state": "serving"
    }
  },
  "spec": {
    "replicas": 2,
    "selector": {
      "app": "lxy"
    },
    "template": {
      "metadata": {
        "labels": {
          "app": "lxy"
        }
      },
      "spec": {
        "volumes": null,
        "containers": [
          {
            "name": "lxy",
            "image": "quay.io/chris_w_beitel/lxy:$WERCKER_GIT_COMMIT",
            "ports": [
              {
                "containerPort": 5000,
                "protocol": "TCP"
              }
            ],
            "imagePullPolicy": "IfNotPresent"
          }
        ],
        "restartPolicy": "Always",
        "dnsPolicy": "ClusterFirst"
      }
    }
  }
}' > cluster/lxy-controller.json

echo '{
   "kind": "Service",
   "apiVersion": "v1alpha",
   "metadata": {
      "name": "lxy",
      "labels": {
         "name": "lxy"
      }
   },
   "spec":{
      "createExternalLoadBalancer": true,
      "ports": [
         {
           "port": 5000,
           "targetPort": "http-server",
           "protocol": "TCP"
         }
      ],
      "selector":{
         "name":"lxy"
      }
   }
}' > cluster/lxy-service.json

