package helper_functions

import (
  "k8s.io/apimachinery/pkg/util/wait"
  "k8s.io/client-go/util/retry"
  "time"
)

var defaultRetry = wait.Backoff{
  Steps:    9,
  Duration: 20 * time.Millisecond,
  Factor:   3.0,
  Jitter:   0.1,
}

func TryWait(f func() error) error {
  return retry.RetryOnConflict(defaultRetry, f)
}
