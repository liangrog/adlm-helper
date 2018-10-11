package file

// Policy configuration file struct
// It's almost the same as the similar
// struct in dlm apart from it has yaml
// annotations so we can read from source
// file in yaml file
type Policy struct {
	Description      string         `yaml:"Description,omitempty"`
	ExecutionRoleArn string         `yaml:"ExecutionRoleArn"`
	State            string         `yaml:"State"`
	PolicyDetails    *PolicyDetails `yaml:"PolicyDetails"`
}

type PolicyDetails struct {
	ResourceTypes string      `yaml:"ResourceTypes"`
	TargetTags    []*Tag      `yaml:"TargetTags"`
	Schedules     []*Schedule `yaml:"Schedules"`
}

type Schedule struct {
	Name       string      `yaml:"Name"`
	CreateRule *CreateRule `yaml:"CreateRule"`
	RetainRule *RetainRule `yaml:"RetainRule"`
	TagsToAdd  []*Tag      `yaml:"Tag,omitempty"`
}

type CreateRule struct {
	Interval     int64     `yaml:"Interval"`
	IntervalUnit string    `yaml:"IntervalUnit"`
	Times        []*string `yaml:"Times"`
}

type RetainRule struct {
	Count int64 `yaml:"Count"`
}

type Tag struct {
	Key   string `yaml:"Key"`
	Value string `yaml:"Value"`
}
