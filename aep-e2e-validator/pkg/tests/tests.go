package tests

func NewTests() []Test {
	return []Test{
		TestAEP132ListResourcesLimit1,
		TestAEP132ListResourcesPageToken,
		TestAEP133Create,
		TestAEP133DuplicateCreationCheck,
		TestAEP134UpdateResource,
		TestAEP135DeleteResource,
		TestAEP135DeleteNonExistentResource,
	}
}
