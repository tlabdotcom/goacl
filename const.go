package goacl

// Define sample roles

var (
	SampleRoles = []Role{
		{Name: "admin", Label: "Administrator", Description: "Has full access to the system"},
		{Name: "user", Label: "Regular User", Description: "Can access user-specific features"},
	}

	// Define sample features
	SampleFeatures = []Feature{
		{Name: "inventory", Description: "Inventory Management"},
		{Name: "sales", Description: "Sales Management"},
	}

	// Define sample sub-features
	SampleSubFeatures = []SubFeature{
		{FeatureID: 1, Name: "view_inventory", Description: "View inventory items"},
		{FeatureID: 1, Name: "edit_inventory", Description: "Edit inventory items"},
		{FeatureID: 2, Name: "view_sales", Description: "View sales transactions"},
		{FeatureID: 2, Name: "create_sales", Description: "Create new sales transactions"},
	}

	// Define sample endpoints
	SampleEndpoints = []Endpoint{
		{Method: "GET", URL: "/inventory", SubFeatureID: 1},
		{Method: "POST", URL: "/inventory", SubFeatureID: 2},
		{Method: "GET", URL: "/sales", SubFeatureID: 3},
		{Method: "POST", URL: "/sales", SubFeatureID: 4},
	}

	// Define sample policies
	SamplePolicies = []Policy{
		{RoleID: 1, FeatureID: 1, SubFeatureID: 1, Status: true}, // Admin can view inventory
		{RoleID: 1, FeatureID: 1, SubFeatureID: 2, Status: true}, // Admin can edit inventory
		{RoleID: 1, FeatureID: 2, SubFeatureID: 3, Status: true}, // Admin can view sales
		{RoleID: 1, FeatureID: 2, SubFeatureID: 4, Status: true}, // Admin can create sales
		{RoleID: 2, FeatureID: 1, SubFeatureID: 1, Status: true}, // User can view inventory
		{RoleID: 2, FeatureID: 2, SubFeatureID: 3, Status: true}, // User can view sales
	}
)
