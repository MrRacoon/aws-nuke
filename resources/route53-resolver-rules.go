package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type Route53ResolverRule struct {
	svc  *route53resolver.Route53Resolver
	id   *string
	name *string
}

func init() {
	register("Route53ResolverRule", ListRoute53ResolverRules)
}

func ListRoute53ResolverRules(sess *session.Session) ([]Resource, error) {
	svc := route53resolver.New(sess)
	resources := make([]Resource, 0)
	params := &route53resolver.ListResolverRulesInput{}

	for {
		resp, err := svc.ListResolverRules(params)

		if err != nil {
			return nil, err
		}

		for _, assoc := range resp.ResolverRules {
			resources = append(resources, &Route53ResolverRule{
				svc:  svc,
				id:   assoc.Id,
				name: assoc.Name,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (rule *Route53ResolverRule) Filter() error {
	if *rule.id == "rslvr-autodefined-rr-internet-resolver" {
		return fmt.Errorf("cannot delete default rule 'rslvr-autodefined-rr-internet-resolver'")
	}

	return nil
}

func (rule *Route53ResolverRule) Remove() error {

	deleteParams := &route53resolver.DeleteResolverRuleInput{
		ResolverRuleId: rule.id,
	}

	_, err := rule.svc.DeleteResolverRule(deleteParams)
	if err != nil {
		return err
	}

	return nil
}

func (rule *Route53ResolverRule) Properties() types.Properties {
	return types.NewProperties().
		Set("Name", rule.name).
		Set("ID", rule.id)
}

func (rule *Route53ResolverRule) String() string {
	if rule.name != nil {
		return fmt.Sprintf("%s (%s)", *rule.id, *rule.name)
	} else {
		return fmt.Sprintf("%s", *rule.id)
	}
}
