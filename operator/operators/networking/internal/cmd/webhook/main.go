package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"

	"github.com/codingconcepts/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kloudlite/operator/common"
	"github.com/kloudlite/operator/pkg/constants"
	"github.com/kloudlite/operator/pkg/logging"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

type Resource string

const (
	ResourcePod     Resource = "pod"
	ResourceService Resource = "service"
)

const (
	podBindingIP        string = "kloudlite.io/podbinding.ip"
	podReservationToken string = "kloudlite.io/podbinding.reservation-token"

	svcBindingIPLabel            string = "kloudlite.io/servicebinding.ip"
	svcReservationTokenLabel     string = "kloudlite.io/servicebinding.reservation-token"
	kloudliteWebhookTriggerLabel string = "kloudlite.io/webhook.trigger"
)

const (
	debugWebhookAnnotation string = "kloudlite.io/networking.webhook.debug"
)

type Env struct {
	GatewayAdminApiAddr string `env:"GATEWAY_ADMIN_API_ADDR" required:"true"`
}

type Flags struct {
	WgImage           string
	WgImagePullPolicy string
}

type HandlerContext struct {
	context.Context
	Env
	Flags
	Resource
	*slog.Logger
}

func main() {
	var ev Env
	if err := env.Set(&ev); err != nil {
		panic(err)
	}

	var addr string
	flag.StringVar(&addr, "addr", "", "--addr <host:port>")

	var logLevel string
	flag.StringVar(&logLevel, "log-level", "info", "--log-level <debug|warn|info|error>")

	var flags Flags

	flag.StringVar(&flags.WgImage, "wg-image", "ghcr.io/kloudlite/hub/wireguard:latest", "--wg-image <image>")

	flag.StringVar(&flags.WgImagePullPolicy, "wg-image-pull-policy", "IfNotPresent", "--wg-image-pull-policy <image-pull-policy>")

	flag.Parse()

	logger := logging.NewSlogLogger(logging.SlogOptions{
		Prefix:        "[webhook]",
		ShowCaller:    true,
		ShowDebugLogs: logLevel == "debug",
	})

	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	httpLogger := logging.NewHttpLogger(logging.HttpLoggerOptions{})
	r.Use(httpLogger.Use)

	r.HandleFunc("/mutate/pod", func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		handleMutate(HandlerContext{Context: r.Context(), Env: ev, Flags: flags, Resource: ResourcePod, Logger: logger.With("request-id", requestID)}, w, r)
	})

	r.HandleFunc("/mutate/service", func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		handleMutate(HandlerContext{Context: r.Context(), Env: ev, Flags: flags, Resource: ResourceService, Logger: logger.With("request-id", requestID)}, w, r)
	})

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	logger.Info("starting http server", "addr", addr)

	common.PrintReadyBanner()

	err := server.ListenAndServeTLS("/tmp/tls/tls.crt", "/tmp/tls/tls.key")
	if err != nil {
		panic(err)
	}
}

func handleMutate(ctx HandlerContext, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	review := admissionv1.AdmissionReview{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err = deserializer.Decode(body, nil, &review); err != nil {
		http.Error(w, "could not decode admission review", http.StatusBadRequest)
		return
	}

	var response admissionv1.AdmissionReview

	switch ctx.Resource {
	case ResourcePod:
		{
			response = processPodAdmission(ctx, review)
		}
	case ResourceService:
		{
			response = processServiceAdmission(ctx, review)
		}
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

func processPodAdmission(ctx HandlerContext, review admissionv1.AdmissionReview) admissionv1.AdmissionReview {
	ctx.InfoContext(ctx, "pod admission", "ref", review.Request.UID, "op", review.Request.Operation)

	switch review.Request.Operation {
	case admissionv1.Create:
		{
			ctx.Info("[INCOMING] pod", "op", review.Request.Operation, "uid", review.Request.UID, "name", review.Request.Name, "namespace", review.Request.Namespace)
			pod := corev1.Pod{}
			err := json.Unmarshal(review.Request.Object.Raw, &pod)
			if err != nil {
				return errResponse(ctx, err, review.Request.UID)
			}

			if pod.GetLabels()[constants.KloudliteGatewayEnabledLabel] == "false" {
				return mutateAndAllow(review, nil)
			}

			isDebugMode := pod.GetAnnotations()[debugWebhookAnnotation] == "true"

			wgContainer := corev1.Container{
				Name:            "kloudlite-wg",
				Image:           ctx.WgImage,
				ImagePullPolicy: corev1.PullPolicy(ctx.WgImagePullPolicy),
				Command: []string{
					"bash",
					"-c",
					fmt.Sprintf(`
cat > /tmp/script.sh <<EOF
while true; do
  set -x
  curl -X PUT --fail --silent %q > /etc/wireguard/kloudlite-wg.conf
  set +x
  ec=\$?
  echo "exit code: \$ec"
  if [ \$ec -eq 0 ]; then
    wg-quick down kloudlite-wg || echo "[starting] wireguard kloudlite-wg intrerface"
    wg-quick up kloudlite-wg
    %s
    echo "[SUCCESS] wireguard is up"
    echo "search $POD_NAMESPACE.svc.cluster.local svc.cluster.local cluster.local" >> /etc/resolv.conf
    echo "options ndots:5" >> /etc/resolv.conf
    exit 0
  fi
  echo "[RETRY] wireguard configuration could not be fetched from gateway ip-manager, retrying in 1 seconds"
  sleep 1
done
EOF

bash /tmp/script.sh
`, fmt.Sprintf("%s/pod/$POD_NAMESPACE/$POD_NAME/$POD_IP", ctx.GatewayAdminApiAddr),
						func() string {
							if isDebugMode {
								return "tail -f /dev/null"
							}
							return "echo waiting 2 seconds to allow wireguard to be ready; sleep 2"
						}(),
					),
				},
				Env: []corev1.EnvVar{
					{
						Name: "POD_IP",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								FieldPath: "status.podIP",
							},
						},
					},
					{
						Name: "POD_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								FieldPath: "metadata.name",
							},
						},
					},
					{
						Name: "POD_NAMESPACE",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								FieldPath: "metadata.namespace",
							},
						},
					},
				},
				SecurityContext: &corev1.SecurityContext{
					Capabilities: &corev1.Capabilities{
						Drop: []corev1.Capability{
							"ALL",
						},
						Add: []corev1.Capability{
							"NET_ADMIN",
						},
					},
				},
			}

			if isDebugMode {
				pod.Spec.Containers = append(pod.Spec.Containers, wgContainer)
			} else {
				pod.Spec.InitContainers = append(pod.Spec.InitContainers, wgContainer)
			}

			lb := pod.GetLabels()
			if lb == nil {
				lb = make(map[string]string, 1)
			}
			lb[constants.KloudliteGatewayEnabledLabel] = "true"

			patchBytes, err := json.Marshal([]map[string]any{
				{
					"op":    "add",
					"path":  "/metadata/labels",
					"value": lb,
				},
				{
					"op":    "add",
					"path":  "/spec",
					"value": pod.Spec,
				},
			})
			if err != nil {
				return errResponse(ctx, err, review.Request.UID)
			}

			return mutateAndAllow(review, patchBytes)
		}
	case admissionv1.Delete:
		{
			pod := corev1.Pod{}
			err := json.Unmarshal(review.Request.OldObject.Raw, &pod)
			if err != nil {
				return errResponse(ctx, err, review.Request.UID)
			}

			if pod.GetDeletionTimestamp() == nil {
				return mutateAndAllow(review, nil)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/pod/%s/%s", ctx.Env.GatewayAdminApiAddr, pod.GetNamespace(), pod.GetName()), nil)
			if err != nil {
				return errResponse(ctx, err, review.Request.UID)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return errResponse(ctx, err, review.Request.UID)
			}

			if resp.StatusCode != 200 {
				return errResponse(ctx, fmt.Errorf("unexpected status code: %d", resp.StatusCode), review.Request.UID)
			}

			return mutateAndAllow(review, nil)
		}
	default:
		{
			return mutateAndAllow(review, nil)
		}
	}
}

func processServiceAdmission(ctx HandlerContext, review admissionv1.AdmissionReview) admissionv1.AdmissionReview {
	switch review.Request.Operation {
	case admissionv1.Create, admissionv1.Update:
		{
			ctx.Info("[INCOMING] service", "op", review.Request.Operation, "name", review.Request.Name, "namespace", review.Request.Namespace)
			svc := corev1.Service{}
			err := json.Unmarshal(review.Request.Object.Raw, &svc)
			if err != nil {
				return errResponse(ctx, err, review.Request.UID)
			}

			lb := svc.GetLabels()
			if lb == nil {
				lb = make(map[string]string, 2)
			}
			lb[constants.KloudliteGatewayEnabledLabel] = "true"
			delete(lb, kloudliteWebhookTriggerLabel)
			svc.SetLabels(lb)

			patchBytes, err := json.Marshal([]map[string]any{
				{
					"op":    "add",
					"path":  "/metadata/labels",
					"value": svc.GetLabels(),
				},
			})
			if err != nil {
				return errResponse(ctx, err, review.Request.UID)
			}

			return mutateAndAllow(review, patchBytes)
		}
	case admissionv1.Delete:
		{
			svc := corev1.Service{}
			err := json.Unmarshal(review.Request.OldObject.Raw, &svc)
			if err != nil {
				return errResponse(ctx, err, review.Request.UID)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/service/%s/%s", ctx.Env.GatewayAdminApiAddr, svc.GetNamespace(), svc.GetName()), nil)
			if err != nil {
				return errResponse(ctx, err, review.Request.UID)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return errResponse(ctx, err, review.Request.UID)
			}

			if resp.StatusCode != 200 {
				return errResponse(ctx, fmt.Errorf("unexpected status code: %d", resp.StatusCode), review.Request.UID)
			}
			return mutateAndAllow(review, nil)
		}
	default:
		{
			return mutateAndAllow(review, nil)
		}
	}
}

func errResponse(ctx HandlerContext, err error, uid types.UID) admissionv1.AdmissionReview {
	ctx.Error("encountered error", "err", err)
	return admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     uid,
			Allowed: false,
			Result: &metav1.Status{
				Message: err.Error(),
			},
		},
	}
}

func mutateAndAllow(review admissionv1.AdmissionReview, patch []byte) admissionv1.AdmissionReview {
	patchType := admissionv1.PatchTypeJSONPatch

	resp := admissionv1.AdmissionResponse{
		UID:     review.Request.UID,
		Allowed: true,
	}

	if patch != nil {
		resp.Patch = patch
		resp.PatchType = &patchType
	}

	return admissionv1.AdmissionReview{
		TypeMeta: review.TypeMeta,
		// Request:  review.Request,
		Response: &resp,
	}
}
