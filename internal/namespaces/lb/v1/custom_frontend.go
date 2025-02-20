package lb

import (
	"context"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

func lbFrontendMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp lb.Frontend
	frontend := tmp(i.(lb.Frontend))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "LB",
		},
		{
			FieldName: "Backend",
		},
	}

	if len(frontend.LB.Tags) != 0 && frontend.LB.Tags[0] == kapsuleTag {
		frontendResp, err := human.Marshal(frontend, opt)
		if err != nil {
			return "", err
		}
		return strings.Join([]string{
			frontendResp,
			warningKapsuleTaggedMessageView(),
		}, "\n\n"), nil
	}

	str, err := human.Marshal(frontend, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func frontendGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptFrontend()
	return c
}

func frontendCreateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptFrontend()
	return c
}

func frontendUpdateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptFrontend()
	return c
}

func frontendDeleteBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptFrontend()
	return c
}

func interceptFrontend() core.CommandInterceptor {
	return func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		if _, ok := res.(*core.SuccessResult); ok {
			getFrontend, err := api.GetFrontend(&lb.ZonedAPIGetFrontendRequest{
				Zone:       argsI.(*lb.ZonedAPIDeleteFrontendRequest).Zone,
				FrontendID: argsI.(*lb.ZonedAPIDeleteFrontendRequest).FrontendID,
			})
			if err != nil {
				return nil, err
			}
			if len(getFrontend.LB.Tags) != 0 && getFrontend.LB.Tags[0] == kapsuleTag {
				return warningKapsuleTaggedMessageView(), nil
			}
		}

		return res, nil
	}
}
