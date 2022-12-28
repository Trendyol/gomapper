package tests

// type Root1 struct {
// 	Flavor *Flavor1
// }

// type Flavor1 struct {
// 	Type     string            `json:"type"`
// 	Override *map[string]Role1 `json:"override,omitempty"`
// }

// type Role1 struct {
// 	DiskSize  *int   `json:"diskSize,omitempty"`
// 	VmFlavor  string `json:"vmFlavor,omitempty"`
// 	NodeCount *int   `json:"nodeCount,omitempty"`
// }

// type Root2 struct {
// 	Flavor *Flavor2
// }

// type Flavor2 struct {
// 	Type     string            `json:"type"`
// 	Override *map[string]Role2 `json:"override,omitempty"`
// }

// type Role2 struct {
// 	DiskSize  *int   `json:"diskSize,omitempty"`
// 	VmFlavor  string `json:"vmFlavor,omitempty"`
// 	NodeCount *int   `json:"nodeCount,omitempty"`
// }

// func Test_Dest_Must(t *testing.T) {
// 	roleMap := make(map[string]Role1)

// 	diskSize := 50
// 	nodeCount := 5

// 	roleMap["data"] = Role1{
// 		DiskSize:  &diskSize,
// 		VmFlavor:  "ty_small",
// 		NodeCount: &nodeCount,
// 	}

// 	source := &Flavor1{
// 		Type:     "small",
// 		Override: &roleMap,
// 	}

// 	var dest Flavor2
// 	gomapper.Map(source, &dest)

// 	assert.NotNil(t, dest.Override)
// }
