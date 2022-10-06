package project

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	artifactsv1 "operators.kloudlite.io/apis/artifacts/v1"
	"operators.kloudlite.io/lib/constants"
	"operators.kloudlite.io/lib/harbor"
	"operators.kloudlite.io/lib/logging"
	rApi "operators.kloudlite.io/lib/operator"
	stepResult "operators.kloudlite.io/lib/operator/step-result"
	"operators.kloudlite.io/operators/harbor/internal/env"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	HarborCli *harbor.Client
	logger    logging.Logger
	Name      string
	Env       *env.Env
}

func (r *Reconciler) GetName() string {
	return r.Name
}

const (
	DefaultsPatched    string = "defaults-patched"
	HarborProjectReady string = "harbor-project-ready"
	HarborWebhookReady string = "harbor-project-webhook-ready"
)

// +kubebuilder:rbac:groups=artifacts.kloudlite.io,resources=harborprojects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=artifacts.kloudlite.io,resources=harborprojects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=artifacts.kloudlite.io,resources=harborprojects/finalizers,verbs=update

func (r *Reconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	req, err := rApi.NewRequest(context.WithValue(ctx, "logger", r.logger), r.Client, request.NamespacedName, &artifactsv1.HarborProject{})
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if req.Object.GetDeletionTimestamp() != nil {
		if x := r.finalize(req); !x.ShouldProceed() {
			return x.ReconcilerResponse()
		}
		return ctrl.Result{}, nil
	}

	req.Logger.Infof("NEW RECONCILATION")

	if step := req.ClearStatusIfAnnotated(); !step.ShouldProceed() {
		return step.ReconcilerResponse()
	}

	if step := req.RestartIfAnnotated(); !step.ShouldProceed() {
		return step.ReconcilerResponse()
	}

	// TODO: initialize all checks here
	if step := req.EnsureChecks(DefaultsPatched, HarborProjectReady, HarborWebhookReady); !step.ShouldProceed() {
		return step.ReconcilerResponse()
	}

	if step := req.EnsureLabelsAndAnnotations(); !step.ShouldProceed() {
		return step.ReconcilerResponse()
	}

	if step := req.EnsureFinalizers(constants.ForegroundFinalizer, constants.CommonFinalizer); !step.ShouldProceed() {
		return step.ReconcilerResponse()
	}

	if step := r.reconDefaults(req); !step.ShouldProceed() {
		return step.ReconcilerResponse()
	}

	if step := r.reconHarborProject(req); !step.ShouldProceed() {
		return step.ReconcilerResponse()
	}

	if step := r.reconHarborWebhook(req); !step.ShouldProceed() {
		return step.ReconcilerResponse()
	}

	req.Object.Status.IsReady = true
	// req.Object.Status.LastReconcileTime = metav1.Time{Time: time.Now()}
	req.Logger.Infof("RECONCILATION COMPLETE")
	return ctrl.Result{RequeueAfter: r.Env.ReconcilePeriod * time.Second}, r.Status().Update(ctx, req.Object)
}

func (r *Reconciler) finalize(req *rApi.Request[*artifactsv1.HarborProject]) stepResult.Result {
	return req.Finalize()
}

func (r *Reconciler) reconDefaults(req *rApi.Request[*artifactsv1.HarborProject]) stepResult.Result {
	ctx, obj, checks := req.Context(), req.Object, req.Object.Status.Checks

	check := rApi.Check{Generation: obj.Generation}

	if obj.Spec.Project == nil || obj.Spec.Webhook == nil {
		obj.Spec.Project = &harbor.Project{
			Name: obj.Name,
		}

		obj.Spec.Webhook = &harbor.Webhook{
			Name: "kloudlite-webhook",
		}

		if err := r.Update(ctx, obj); err != nil {
			return req.CheckFailed(DefaultsPatched, check, err.Error())
		}

		checks[DefaultsPatched] = check
		return req.UpdateStatus()
	}

	check.Status = true
	if check != checks[DefaultsPatched] {
		checks[DefaultsPatched] = check
		return req.UpdateStatus()
	}

	return req.Next()
}

func (r *Reconciler) reconHarborProject(req *rApi.Request[*artifactsv1.HarborProject]) stepResult.Result {
	ctx, obj, checks := req.Context(), req.Object, req.Object.Status.Checks

	check := rApi.Check{Generation: obj.Generation}

	exists, err := r.HarborCli.CheckIfProjectExists(ctx, obj.Spec.Project.Name)
	if err != nil {
		return req.CheckFailed(HarborProjectReady, check, err.Error())
	}

	if !exists {
		project, err := r.HarborCli.CreateProject(ctx, obj.Spec.Project.Name)
		if err != nil {
			return req.CheckFailed(HarborProjectReady, check, err.Error())
		}
		obj.Spec.Project = project
		if err := r.Update(ctx, obj); err != nil {
			return nil
		}
		checks[HarborProjectReady] = check
		return req.UpdateStatus()
	}

	if obj.Spec.Project.Location == "" {
		req.Logger.Infof("project location is empty, going to query harbor for it")
		project, err := r.HarborCli.GetProject(ctx, obj.Spec.Project.Name)
		if err != nil {
			return req.CheckFailed(HarborProjectReady, check, err.Error())
		}
		obj.Spec.Project = project
		if err := r.Update(ctx, obj); err != nil {
			return req.CheckFailed(HarborProjectReady, check, err.Error())
		}
		checks[HarborProjectReady] = check
		return req.UpdateStatus()
	}

	check.Status = true
	if check != checks[HarborProjectReady] {
		checks[HarborProjectReady] = check
		return req.UpdateStatus()
	}

	return req.Next()
}

func (r *Reconciler) reconHarborWebhook(req *rApi.Request[*artifactsv1.HarborProject]) stepResult.Result {
	ctx, obj, checks := req.Context(), req.Object, req.Object.Status.Checks
	check := rApi.Check{Generation: obj.Generation}

	exists, err := r.HarborCli.CheckWebhookExists(ctx, obj.Spec.Webhook)
	if err != nil {
		return req.CheckFailed(HarborWebhookReady, check, err.Error())
	}

	if !exists {
		webhook, err := r.HarborCli.CreateWebhook(ctx, obj.Spec.Project.Name, obj.Spec.Webhook.Name)
		if err != nil {
			return req.CheckFailed(HarborWebhookReady, check, err.Error())
		}
		obj.Spec.Webhook = webhook
		if err := r.Update(ctx, obj); err != nil {
			return req.CheckFailed(HarborProjectReady, check, err.Error())
		}
		checks[HarborWebhookReady] = check
		return req.UpdateStatus()
	}

	check.Status = true
	if check != checks[HarborWebhookReady] {
		checks[HarborWebhookReady] = check
		return req.UpdateStatus()
	}

	return req.Next()
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager, logger logging.Logger) error {
	r.Client = mgr.GetClient()
	r.Scheme = mgr.GetScheme()
	r.logger = logger.WithName(r.Name)

	builder := ctrl.NewControllerManagedBy(mgr).For(&artifactsv1.HarborProject{})
	return builder.Complete(r)
}
