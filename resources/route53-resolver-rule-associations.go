package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type Route53ResolverRuleAssociation struct {
	svc   *route53resolver.Route53Resolver
	id    *string
	name  *string
	vpcID *string
}

func init() {
	register("Route53ResolverRuleAssociation", ListRoute53ResolverRuleAssociations)
}

func ListRoute53ResolverRuleAssociations(sess *session.Session) ([]Resource, error) {
	svc := route53resolver.New(sess)
	resources := make([]Resource, 0)
	params := &route53resolver.ListResolverRuleAssociationsInput{}

	for {
		resp, err := svc.ListResolverRuleAssociations(params)

		if err != nil {
			return nil, err
		}

		for _, assoc := range resp.ResolverRuleAssociations {
			resources = append(resources, &Route53ResolverRuleAssociation{
				svc:   svc,
				name:  assoc.Name,
				id:    assoc.ResolverRuleId,
				vpcID: assoc.VPCId,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (assoc *Route53ResolverRuleAssociation) Filter() error {
	if *assoc.id == "rslvr-autodefined-rr-internet-resolver" {
		return fmt.Errorf("cannot delete default rule association for 'rslvr-autodefined-rr-internet-resolver'")
	}

	return nil
}

func (assoc *Route53ResolverRuleAssociation) Remove() error {
	dissocParams := &route53resolver.DisassociateResolverRuleInput{
		ResolverRuleId: assoc.id,
		VPCId:          assoc.vpcID,
	}

	_, err := assoc.svc.DisassociateResolverRule(dissocParams)
	if err != nil {
		return err
	}

	return nil
}

func (assoc *Route53ResolverRuleAssociation) Properties() types.Properties {
	return types.NewProperties().
		Set("Name", assoc.name).
		Set("ID", assoc.id).
		Set("VPCId", assoc.vpcID)
}

func (assoc *Route53ResolverRuleAssociation) String() string {
	if assoc.name == nil {
		return fmt.Sprintf("%s", *assoc.id)
	} else {
		return fmt.Sprintf("%s (%s)", *assoc.id, *assoc.name)
	}
}
