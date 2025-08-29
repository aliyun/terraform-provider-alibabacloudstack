package schema

import (
	"github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
)


type FullProviderSchemaOutput struct {
    FormatVersion   string                       `json:"format_version"`
    ProviderSchemas map[string]*tfjson.ProviderSchema `json:"provider_schemas"`
}

func ConvertAndWrapProviderSchema(providerName string, provider *schema.Provider) *FullProviderSchemaOutput {
    return &FullProviderSchemaOutput{
        FormatVersion: "1.0",
        ProviderSchemas: map[string]*tfjson.ProviderSchema{
            providerName: ConvertProviderSchema(provider),
        },
    }
}

func ConvertProviderSchema(provider *schema.Provider) *tfjson.ProviderSchema {
	providerSchema := &tfjson.ProviderSchema{
		ConfigSchema:      convertProviderConfigSchema(provider.Schema),
		ResourceSchemas:   make(map[string]*tfjson.Schema),
		DataSourceSchemas: make(map[string]*tfjson.Schema),
	}

	for name, res := range provider.ResourcesMap {
		providerSchema.ResourceSchemas[name] = convertResourceSchema(res)
	}

	for name, ds := range provider.DataSourcesMap {
		providerSchema.DataSourceSchemas[name] = convertResourceSchema(ds)
	}

	return providerSchema
}

func convertProviderConfigSchema(schemaMap map[string]*schema.Schema) *tfjson.Schema {
	if schemaMap == nil {
		return nil
	}

	return &tfjson.Schema{
		Block: convertSchemaBlock(schemaMap),
	}
}

func convertResourceSchema(res *schema.Resource) *tfjson.Schema {
	return &tfjson.Schema{
		Block: convertSchemaBlock(res.Schema),
	}
}

func convertSchemaBlock(schemaMap map[string]*schema.Schema) *tfjson.SchemaBlock {
	block := &tfjson.SchemaBlock{
		Attributes:   make(map[string]*tfjson.SchemaAttribute),
		NestedBlocks: make(map[string]*tfjson.SchemaBlockType),
	}

	for name, s := range schemaMap {
		if s.Elem != nil {
			if _, isResource := s.Elem.(*schema.Resource); isResource {
				blockType := convertBlockType(s)
				block.NestedBlocks[name] = blockType
				continue
			}
		}

		attribute := &tfjson.SchemaAttribute{
			AttributeType: convertType(s.Type),
			Description:   s.Description,
			Required:      s.Required,
			Optional:      s.Optional,
			Computed:      s.Computed,
			Sensitive:     s.Sensitive,
		}
		block.Attributes[name] = attribute
	}

	return block
}

func convertBlockType(s *schema.Schema) *tfjson.SchemaBlockType {
	nestedBlock := &tfjson.SchemaBlockType{
		NestingMode: convertNestingMode(s.Type),
		Block:       convertSchemaBlock(s.Elem.(*schema.Resource).Schema),
	}

	if s.MaxItems > 0 {
		nestedBlock.MinItems = uint64(s.MinItems)
		nestedBlock.MaxItems = uint64(s.MaxItems)
	}

	return nestedBlock
}

func convertNestingMode(t schema.ValueType) tfjson.SchemaNestingMode {
	switch t {
	case schema.TypeList:
		return tfjson.SchemaNestingModeList
	case schema.TypeSet:
		return tfjson.SchemaNestingModeSet
	case schema.TypeMap:
		return tfjson.SchemaNestingModeMap
	default:
		return tfjson.SchemaNestingModeSingle
	}
}

func convertType(t schema.ValueType) cty.Type {
	switch t {
	case schema.TypeString:
		return cty.String
	case schema.TypeBool:
		return cty.Bool
	case schema.TypeInt, schema.TypeFloat:
		return cty.Number
	case schema.TypeList, schema.TypeSet:
		elemType := cty.DynamicPseudoType
		return cty.List(elemType)
	case schema.TypeMap:
		return cty.Map(cty.DynamicPseudoType)
	default:
		return cty.NilType
	}
}
