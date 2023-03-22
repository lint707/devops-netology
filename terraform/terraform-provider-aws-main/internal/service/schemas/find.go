package schemas

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/schemas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func FindDiscovererByID(conn *schemas.Schemas, id string) (*schemas.DescribeDiscovererOutput, error) {
	input := &schemas.DescribeDiscovererInput{
		DiscovererId: aws.String(id),
	}

	output, err := conn.DescribeDiscoverer(input)

	if tfawserr.ErrCodeEquals(err, schemas.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, &resource.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return output, nil
}

func FindRegistryByName(conn *schemas.Schemas, name string) (*schemas.DescribeRegistryOutput, error) {
	input := &schemas.DescribeRegistryInput{
		RegistryName: aws.String(name),
	}

	output, err := conn.DescribeRegistry(input)

	if tfawserr.ErrCodeEquals(err, schemas.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, &resource.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return output, nil
}

func FindSchemaByNameAndRegistryName(conn *schemas.Schemas, name, registryName string) (*schemas.DescribeSchemaOutput, error) {
	input := &schemas.DescribeSchemaInput{
		RegistryName: aws.String(registryName),
		SchemaName:   aws.String(name),
	}

	output, err := conn.DescribeSchema(input)

	if tfawserr.ErrCodeEquals(err, schemas.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, &resource.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return output, nil
}
