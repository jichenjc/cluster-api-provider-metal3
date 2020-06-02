/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package baremetal

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/cluster-api/util/patch"
)

// Filter filters a list for a string.
func Filter(list []string, strToFilter string) (newList []string) {
	for _, item := range list {
		if item != strToFilter {
			newList = append(newList, item)
		}
	}
	return
}

// Contains returns true if a list contains a string.
func Contains(list []string, strToSearch string) bool {
	for _, item := range list {
		if item == strToSearch {
			return true
		}
	}
	return false
}

// NotFoundError represents that an object was not found
type NotFoundError struct {
}

// Error implements the error interface
func (e *NotFoundError) Error() string {
	return "Object not found"
}

func patchIfFound(ctx context.Context, helper *patch.Helper, host runtime.Object) error {
	err := helper.Patch(ctx, host)
	if err != nil {
		notFound := true
		if aggr, ok := err.(kerrors.Aggregate); ok {
			for _, kerr := range aggr.Errors() {
				if !apierrors.IsNotFound(kerr) {
					notFound = false
				}
				if apierrors.IsConflict(kerr) {
					return &RequeueAfterError{}
				}
			}
		} else {
			notFound = false
		}
		if notFound {
			return nil
		}
	}
	return err
}
