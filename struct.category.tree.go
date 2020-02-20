package main

//CATEGORYTREENODE document structure
type CATEGORYTREENODE struct {
	CategoryID string            `json:"CategoryID" bson:"CategoryID" validate:"min=1,max=100,hasNoSpaces"`
	Name       string            `json:"Name" bson:"Name" validate:"min=1,max=100"`
	Images     map[string]string `json:"Images" bson:"Images" validate:"min=1,max=100"`
	Parent     string            `json:"Parent" bson:"Parent" validate:"min=1,max=100"`
	Children   []string          `json:"Children" bson:"Children" validate:"min=1,max=100"`
	Path       string            `json:"Path" bson:"Path"`
	SKUs       []string          `json:"SKUs" bson:"SKUs" validate:"min=1,max=100,hasNoSpaces"`
}
