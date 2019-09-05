/*
Copyright 2015 The Kubernetes Authors.

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

package podgc

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/sets"
)

func newTestBasicUniqueWorkQueue() (*basicUniqueWorkQueue, *clock.FakeClock) {
	fakeClock := clock.NewFakeClock(time.Now())
	wq := &basicUniqueWorkQueue{
		clock: fakeClock,
		queue: make(map[string]time.Time),
	}
	return wq, fakeClock
}

func compareResults(t *testing.T, expected, actual []string) {
	expectedSet := sets.NewString()
	for _, u := range expected {
		expectedSet.Insert(u)
	}
	actualSet := sets.NewString()
	for _, u := range actual {
		actualSet.Insert(u)
	}
	if !expectedSet.Equal(actualSet) {
		t.Errorf("Expected %#v, got %#v", expectedSet.List(), actualSet.List())
	}
}

func TestGetWork(t *testing.T) {
	q, clock := newTestBasicUniqueWorkQueue()
	q.EnqueueIfNew("foo1", -1*time.Minute)
	q.EnqueueIfNew("foo2", -1*time.Minute)
	q.EnqueueIfNew("foo3", 5*time.Minute)
	q.EnqueueIfNew("foo4", 5*time.Minute)
	q.EnqueueIfNew("foo4", 1*time.Minute)
	q.EnqueueIfNew("foo5", 1*time.Minute)
	expected := []string{"foo1", "foo2"}
	compareResults(t, expected, q.GetWork())
	compareResults(t, []string{}, q.GetWork())
	// Dial the time to 2 minutes ahead.
	clock.Step(2 * time.Minute)
	expected = []string{"foo5"}
	compareResults(t, expected, q.GetWork())
	compareResults(t, []string{}, q.GetWork())
	// Dial the time to 1 hour ahead.
	clock.Step(time.Hour)
	expected = []string{"foo3", "foo4"}
	compareResults(t, expected, q.GetWork())
	compareResults(t, []string{}, q.GetWork())
}
