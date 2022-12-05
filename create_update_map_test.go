package gormupdatemap

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGormpatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gormupdatemap Suite")
}

var _ = Describe("CreateUpdateMap", Ordered, func() {
	It("Creates a valid map", func() {
		type t struct {
			Field1 *bool `json:"field_one"`
			Field2 *bool `json:"field_two"`
			Field3 *bool `json:"field_three"          admin_only:"true"`
			Field4 bool  `json:"field_four"`
			Field5 *bool `json:"field_five,omitempty"`
		}

		By("Ignores nil fields")
		v := t{
			Field1: nil,
			Field2: nil,
		}

		patchMap, validationErr := CreateUpdateMap(v, false)
		Expect(validationErr).To(BeNil())
		Expect(patchMap).To(Not(HaveKey("field_one")))
		Expect(patchMap).To(Not(HaveKey("field_two")))

		By("Sets non-nil fields")
		v = t{
			Field1: boolPtr(false),
			Field2: nil,
		}

		patchMap, validationErr = CreateUpdateMap(v, false)
		Expect(patchMap).To(HaveKey("field_one"))
		Expect(patchMap).To(Not(HaveKey("field_two")))

		By("Ignoring non-pointer fields")
		v = t{
			Field1: boolPtr(true),
			Field4: true,
		}
		patchMap, validationErr = CreateUpdateMap(v, false)
		Expect(patchMap).To(HaveKey("field_one"))
		Expect(patchMap).To(Not(HaveKey("field_four")))

		By("Checks for admin permissions")
		By("Without admin permissions")
		v = t{
			Field1: boolPtr(false),
			Field2: boolPtr(true),
			Field3: boolPtr(false),
		}
		patchMap, validationErr = CreateUpdateMap(v, false)
		Expect(validationErr).To(Not(BeNil()))
		Expect(*validationErr).To(Equal("Only admins can update field 'field_three'."))

		By("With admin permissions")
		patchMap, validationErr = CreateUpdateMap(v, true)
		Expect(patchMap).To(HaveKey("field_one"))
		Expect(patchMap).To(HaveKey("field_two"))
		Expect(patchMap).To(HaveKey("field_three"))

		By("Grabbing the first json tag value")
		v = t{
			Field2: boolPtr(true),
			Field5: boolPtr(true),
		}
		patchMap, validationErr = CreateUpdateMap(v, false)
		Expect(patchMap).To(HaveKey("field_two"))
		Expect(patchMap).To(HaveKey("field_five"))
	})

	It("Follows GORM column naming conventions", func() {
		type updateBusiness struct {
			Limit24 *float64 `json:"limit_24"`
		}

		By("Removing underscores before digits")
		patchMap, validationErr := CreateUpdateMap(updateBusiness{Limit24: float64Ptr(10.10)}, false)
		Expect(validationErr).To(BeNil())
		Expect(patchMap).To(HaveKey("limit24"))
	})
})
