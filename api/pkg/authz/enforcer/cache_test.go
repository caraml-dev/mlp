package enforcer

import (
	"testing"

	"github.com/ory/keto-client-go/models"
	"github.com/stretchr/testify/assert"
)

func TestBuildCacheKey(t *testing.T) {
	tests := map[string]struct {
		cache    *InMemoryCache
		input    models.OryAccessControlPolicyAllowedInput
		expected string
	}{
		"first item": {
			cache: &InMemoryCache{
				actionMap:   map[string]string{},
				resourceMap: map[string]string{},
				subjectMap:  map[string]string{},
			},
			input: models.OryAccessControlPolicyAllowedInput{
				Action:   "abc",
				Resource: "def",
				Subject:  "xyz",
			},
			expected: "1:1:1",
		},
		"all items exist": {
			cache: &InMemoryCache{
				actionMap: map[string]string{
					"abc": "2",
				},
				resourceMap: map[string]string{
					"def": "4",
				},
				subjectMap: map[string]string{
					"xyz": "10",
				},
			},
			input: models.OryAccessControlPolicyAllowedInput{
				Action:   "abc",
				Resource: "def",
				Subject:  "xyz",
			},
			expected: "2:4:10",
		},
		"some items exist": {
			cache: &InMemoryCache{
				actionMap: map[string]string{
					"abc": "0",
				},
				resourceMap: map[string]string{
					"def": "4",
				},
				subjectMap: map[string]string{
					"xyz": "0",
				},
			},
			input: models.OryAccessControlPolicyAllowedInput{
				Action:   "abc",
				Resource: "abc",
				Subject:  "xyz",
			},
			expected: "0:2:0", // Resource ID (missing) will be generated using the map key count
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.cache.buildCacheKey(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestLookUp(t *testing.T) {
	trueValue, falseValue := true, false
	type seedData struct {
		item      models.OryAccessControlPolicyAllowedInput
		isAllowed *bool
	}

	tests := map[string]struct {
		seedData      []seedData
		input         models.OryAccessControlPolicyAllowedInput
		expectedVal   *bool
		expectedFound bool
	}{
		"item exists | true value": {
			seedData: []seedData{
				{
					item: models.OryAccessControlPolicyAllowedInput{
						Action:   "abc",
						Resource: "def",
						Subject:  "xyz",
					},
					isAllowed: &trueValue,
				},
				{
					item: models.OryAccessControlPolicyAllowedInput{
						Action:   "def",
						Resource: "ghi",
						Subject:  "xyz",
					},
					isAllowed: &falseValue,
				},
			},
			input: models.OryAccessControlPolicyAllowedInput{
				Action:   "abc",
				Resource: "def",
				Subject:  "xyz",
			},
			expectedVal:   &trueValue,
			expectedFound: true,
		},
		"item exists | false value": {
			seedData: []seedData{
				{
					item: models.OryAccessControlPolicyAllowedInput{
						Action:   "abc",
						Resource: "def",
						Subject:  "xyz",
					},
					isAllowed: &trueValue,
				},
				{
					item: models.OryAccessControlPolicyAllowedInput{
						Action:   "def",
						Resource: "ghi",
						Subject:  "xyz",
					},
					isAllowed: &falseValue,
				},
			},
			input: models.OryAccessControlPolicyAllowedInput{
				Action:   "def",
				Resource: "ghi",
				Subject:  "xyz",
			},
			expectedVal:   &falseValue,
			expectedFound: true,
		},
		"item does not exist": {
			seedData: []seedData{
				{
					item: models.OryAccessControlPolicyAllowedInput{
						Action:   "abc",
						Resource: "def",
						Subject:  "xyz",
					},
					isAllowed: &trueValue,
				},
			},
			input: models.OryAccessControlPolicyAllowedInput{
				Action:   "def",
				Resource: "ghi",
				Subject:  "xyz",
			},
			expectedFound: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cache := newInMemoryCache(60, 60)
			for _, item := range tt.seedData {
				cache.StorePermission(item.item, item.isAllowed)
			}
			cachedVal, found := cache.LookUpPermission(tt.input)
			assert.Equal(t, tt.expectedVal, cachedVal)
			assert.Equal(t, tt.expectedFound, found)
		})
	}
}
