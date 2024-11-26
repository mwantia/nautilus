package plugin

const (
	None     PipelineProcessorCapabilityType = 0
	Metadata PipelineProcessorCapabilityType = 1
	Storage  PipelineProcessorCapabilityType = 2
)

type PipelineProcessorCapabilityType int32

type PipelineProcessorCapability struct {
	Types []PipelineProcessorCapabilityType `json:"types"`
}

var PipelineProcessorCapabilityType_Name = map[int32]string{
	0: "None",
	1: "Metadata",
	2: "Storage",
}

var PipelineProcessorCapabilityType_Value = map[string]int32{
	"None":     0,
	"Metadata": 1,
	"Storage":  2,
}

func (t PipelineProcessorCapabilityType) String() string {
	value := int32(t)

	name, exists := PipelineProcessorCapabilityType_Name[value]
	if !exists {
		return ""
	}

	return name
}
