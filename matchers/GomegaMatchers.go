package matchers

import (
	"fmt"
	"github.com/guzzlerio/rizo"
	"github.com/onsi/gomega/types"
)

func Find(predicates ...rizo.HTTPRequestPredicate) types.GomegaMatcher {
	return &rizoFindMatcher{predicates, []string{}, nil}

}

type rizoFindMatcher struct {
	predicates    []rizo.HTTPRequestPredicate
	failedMatches []string
	testServer    *rizo.RequestRecordingServer
}

func (this *rizoFindMatcher) Match(actual interface{}) (success bool, err error) {
	testServer, ok := actual.(*rizo.RequestRecordingServer)
	if !ok {
		return false, fmt.Errorf("Booom")
	}
	this.testServer = testServer
	for _, predicate := range this.predicates {
		if ok = this.testServer.Find(predicate); !ok {
			this.failedMatches = append(this.failedMatches, predicate.String())
		}
	}
	fmt.Println(len(this.failedMatches))
	return len(this.failedMatches) == 0, nil
}
func (this *rizoFindMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected predicates to match: %v %s", this.failedMatches, this.testServer.Requests)
}
func (this *rizoFindMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected predicates to not match: %v", this.failedMatches)
}
