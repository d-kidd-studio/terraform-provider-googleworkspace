// SPDX-License-Identifier: MPL-2.0

package googleworkspace

import (
	"testing"
	"time"
)

func TestConsistencyCheckReachedConsistency(t *testing.T) {
	numInserts := 3

	cc := consistencyCheck{
		timeout:        time.Duration(time.Minute * 5),
		currConsistent: 1,
		etagChanges:    1,
		lastEtag:       "12345",
	}

	if cc.reachedConsistency(numInserts) {
		t.Errorf("Failed: reached consistency (numInserts: %d, currConsistent: %d, etagChanges: %d, timeout: %d)", numInserts, cc.currConsistent, cc.etagChanges, int(cc.timeout.Minutes()))
	}

	cc.etagChanges = 2
	cc.currConsistent = 18

	if !cc.reachedConsistency(numInserts) {
		t.Errorf("Failed: did not reach consistency (numInserts: %d, currConsistent: %d, etagChanges: %d, timeout: %d)", numInserts, cc.currConsistent, cc.etagChanges, int(cc.timeout.Minutes()))
	}

	cc.etagChanges = 3
	cc.currConsistent = 1

	if cc.reachedConsistency(numInserts) {
		t.Errorf("Failed: reached consistency (numInserts: %d, currConsistent: %d, etagChanges: %d, timeout: %d)", numInserts, cc.currConsistent, cc.etagChanges, int(cc.timeout.Minutes()))
	}

	cc.currConsistent = numInserts

	if !cc.reachedConsistency(numInserts) {
		t.Errorf("Failed: did not reach consistency (numInserts: %d, currConsistent: %d, etagChanges: %d, timeout: %d)", numInserts, cc.currConsistent, cc.etagChanges, int(cc.timeout.Minutes()))
	}
}

func TestConsistencyHandleNewEtag(t *testing.T) {
	cc := consistencyCheck{
		resourceType: "test",
	}

	cc.handleNewEtag("12345")
	if cc.currConsistent != 0 {
		t.Errorf("Failed ['12345']: new etag shows currConsistent > 0 (%d)", cc.currConsistent)
	}

	cc.handleNewEtag("abcde")
	if cc.lastEtag != "abcde" {
		t.Errorf("Failed ['abcde']: ends with incorrect lastEtag (%s)", cc.lastEtag)
	}

	cc.handleNewEtag("54321")
	if cc.etagChanges != 3 {
		t.Errorf("Failed ['abcde']: shows more/less etag changes (expected: %d, got: %d)", 3, cc.etagChanges)
	}
}
