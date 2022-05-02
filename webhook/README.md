# Webhooks

This package contains methods and types used to customize validating/mutating webhooks for resources

## Motivation

There are a few factors that motivate customizing new webhooks:

 - provide a `context.Context` for each function: Just like `knative/pkg` and others, using a context inside each function can provide multiple helper methods to check specific request type, fetch the old object, etc.
 - `logger` fetched from context: Instead of having a global variable hosting an imutable logger, this project already supports having multiple loggers with dynamic log level controlled by a configmap inside the cluster. This project uses uber's zap logger in its Suggared implementation, with specific `Debug`, `Info`, `Warning`, etc. methods. The default `controller-runtime` provided logger has very unclear log level (0-9) which is counter-intuitive leading to not-standard and unexpected log output
 - Adding validation/mutation methods outside of the `apis` scope: For example validating if another object exists, adding feature-flag controlled validations, etc.


## Usage

A few changes are needed in the current objects in order to use this customized webhooks:

 1. Add `context.Context` to `ValidateCreate`, `ValidateUpdate`, `ValidateDelete` and `Default` methods as the first parameter. `Default` and/or `Validation` can be changed as necessary. If required, one or both can keep using regular webhooks.
 1. Import `"github.com/katanomi/pkg/webhook/admission"` on your `<type>_webhook.go` files
 1. Change the interface check from `var _ webhook.Defaulter = &Type{}` to  `var _ admission.Defaulter = &Type{}` the last one coming from this package. Do the same for `Validator`
 1. Use `"knative.dev/pkg/logging"` to fetch a logger from the context and `"knative.dev/pkg/apis"` methods like `apis.IsWithinCreate`
 1. Change one of the below methods to setup the webhook within the app


> The [sample-controller](../examples/sample-controller) has full usable example implemented

### As default setup for object

In this scenario, the `apis/<version>` package will declare all setup methods necessary and nothing else needs to be changed

 - Implement `sharedmain.WebhookRegisterSetup` interface from [`sharedmain`](../sharedmain) package:

*Here is a simple example from `deliveries`*:

```golang
var _ sharedmain.WebhookRegisterSetup = &Delivery{}

// GetLoggerName returns the logger name
func (r *Delivery) GetLoggerName() string {
	return "delivery"
}

// SetupRegisterWithManager implements method to register a webhook
func (r *Delivery) SetupRegisterWithManager(ctx context.Context, mgr ctrl.Manager) {
    // gets logger for this webhook with the specific level
	log := logging.FromContext(ctx)
    // defaulting and validating should be added here, if not added the standard one will be used instead
	err := admission.RegisterDefaultWebhookFor(ctx, mgr, r, admission.WithCreatedBy())
	if err != nil {
		log.Fatalw("register webhook failed", "err", err)
	}

	// adds validate webhook for object
	err := admission.RegisterValidateWebhookFor(ctx, mgr, r, []admission.ValidateCreateFunc{}, []admission.ValidateUpdateFunc{}, []admission.ValidateDeleteFunc{})
	if err != nil {
		log.Fatalw("register webhook failed", "err", err)
	}
}
```

Thats it. Once the object is passed to the `sharedmain.App("").Webhooks()` method, it will automatically will use the new methods, as seen below from `deliveries` controller:

```golang
  sharedmain.App("deliveries-controller").
		Scheme(scheme).
		Log().
		Profiling().
		Controllers(
			&deliveriescontrollers.StageReconciler{},
			&deliveriescontrollers.StageRunReconciler{},
			&deliveriescontrollers.DeliveryReconciler{},
			&deliveriescontrollers.DeliveryRunReconciler{},
			&deliveriescontrollers.ClusterStageReconciler{},
		).
		Webhooks(
			&deliveriesv1alpha1.Stage{},
			&deliveriesv1alpha1.StageRun{},
			&deliveriesv1alpha1.Delivery{},
			&deliveriesv1alpha1.DeliveryRun{},
			&deliveriesv1alpha1.ClusterStage{},
		).
		Run()
```


### As custom webhook with external methods

Using this method will give the ability to extend validation methods with new extra methods while keeping the default method pristine.

 - Implements all the custom `Validation` and `Default` methods, see [`validation.go`](admission/validation.go) and [`transform.go`](admission/transform.go)
 - Init the webhook object inside the `sharedmain.App("").Webhooks()` method:

*Lets say we want to add new defaulting and validating methods to the `Delivery` object*:

```golang
  sharedmain.App("deliveries-controller").
		Scheme(scheme).
		Log().
		Profiling().
		Controllers(
			&deliveriescontrollers.StageReconciler{},
			&deliveriescontrollers.StageRunReconciler{},
			&deliveriescontrollers.DeliveryReconciler{},
			&deliveriescontrollers.DeliveryRunReconciler{},
			&deliveriescontrollers.ClusterStageReconciler{},
		).
		Webhooks(
			&deliveriesv1alpha1.Stage{},
			&deliveriesv1alpha1.StageRun{},
			// custom default webhook
			admission.NewDefaulterWebhook(&deliveriesv1alpha1.Delivery{}).WithTransformer(
				somepackage.SomeTransformerMethod,
				somepackage.AnotherMethod,
			),
			// custom validation webhook
			admission.NewValidatorWebhook(&deliveriesv1alpha1.Delivery{}).
				WithValidateCreate(somepackage.SomeValidateCreateFunc).
				WithValidateUpdate(somepackage.SomeValidateUpdateFunc).
				WithValidateDelete(somepackage.SomeValidateDeleteFunc),
			&deliveriesv1alpha1.Delivery{},
			&deliveriesv1alpha1.DeliveryRun{},
			&deliveriesv1alpha1.ClusterStage{},
		).
		Run()
```

## Implementing validation and default methods

A brief introduction on the methods, see  [`validation.go`](admission/validation.go) and [`transform.go`](admission/transform.go) for specifics

### Default (transform) methods

```golang
// TransformFuncused to make common defaulting logic amongst multiple resource
// using a context, an object and a request
type TransformFunc func(context.Context, runtime.Object, admission.Request)
```

An implementation example for automatically adding the current update time/creation time would be:

```golang
func WithCreateUpdateTimes() TransformFunc {
	return func(ctx context.Context, obj runtime.Object, req admission.Request) {
		metaobj, ok := obj.(metav1.Object)
		if !ok {
			return
		}
		log := logging.FromContext(ctx)
		annotations := metaobj.GetAnnotations()
		if annotations == nil {
			annotations = map[string]string{}
		}

		now := time.Now().Format(time.RFC3339)
		if apis.IsInCreate(ctx) {
			annotations["createdAt"] = now
		} else if apis.IsInUpdate(ctx) {
			annotations["updatedAt"] = now
		}
		metaobj.SetAnnotations(annotations)
	}
}
```

### Validation methods

```golang

// ValidateCreateFunc function to add validation functions when operation is create
// using a context, an object and a request
type ValidateCreateFunc func(ctx context.Context, obj runtime.Object, req admission.Request) error

// ValidateUpdateFunc function to add validation functions when operation is update
// using a context, the current object, the old object and a request
type ValidateUpdateFunc func(ctx context.Context, obj runtime.Object, old runtime.Object, req admission.Request) error

// ValidateDeleteFunc function to add validation functions when operation is delete
// using a context, an object and a request
type ValidateDeleteFunc func(ctx context.Context, obj runtime.Object, req admission.Request) error
```

An implementation example for validating if the `createdAt` annotation is preset:

```golang
func HasCreatedAtAnnotation() ValidateCreateFunc {
	return func(ctx context.Context, obj runtime.Object, req admission.Request) error {
		metaobj, ok := obj.(metav1.Object)
		if !ok {
			return
		}
		log := logging.FromContext(ctx)
		annotations := metaobj.GetAnnotations()
		if annotations == nil {
			annotations = map[string]string{}
		}
		if annotations["createdAt"] == "" {
			return fmt.Errorf("some validation error")
		}
		return nil
	}
}
```
