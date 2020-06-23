package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type Route53ResolverEndpoint struct {
	svc  *route53resolver.Route53Resolver
	id   *string
	name *string
}

func init() {
	register("Route53ResolverEndpoint", ListRoute53ResolverEndpoints)
}

func ListRoute53ResolverEndpoints(sess *session.Session) ([]Resource, error) {
	svc := route53resolver.New(sess)
	resources := make([]Resource, 0)
	params := &route53resolver.ListResolverEndpointsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.ListResolverEndpoints(params)
		if err != nil {
			return nil, err
		}

		for _, endp := range resp.ResolverEndpoints {
			resources = append(resources, &Route53ResolverEndpoint{
				svc:  svc,
				id:   endp.Id,
				name: endp.Name,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (endp *Route53ResolverEndpoint) Remove() error {
	params := &route53resolver.DeleteResolverEndpointInput{
		ResolverEndpointId: endp.id,
	}

	_, err := endp.svc.DeleteResolverEndpoint(params)
	if err != nil {
		return err
	}

	return nil
}

func (endp *Route53ResolverEndpoint) Properties() types.Properties {
	return types.NewProperties().
		Set("Name", endp.name).
		Set("ID", endp.id)
}

func (endp *Route53ResolverEndpoint) String() string {
	return fmt.Sprintf("%s (%s)", *endp.id, *endp.name)
}
