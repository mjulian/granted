package cfaws

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	ssotypes "github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/pkg/browser"
)

func (c *CFSharedConfig) SSOLogin(ctx context.Context) (aws.Credentials, error) {
	if c.ProfileType != ProfileTypeSSO {
		return aws.Credentials{}, errors.New("cannot ssologin to non sso profile")
	}
	// IO = Determine Required Config Values to Establish SSO Session ✅
	// IO = Establish an SSO Session with config vars ✅
	// IO = Retreive any relevant Credentials from the SSO Session
	// IO = Export the credentials to the environment

	rootProfile := c
	requiresAssuming := false
	if len(c.Parents) > 0 {
		rootProfile = c.Parents[0]
		requiresAssuming = true
	}
	cfg, err := rootProfile.AwsConfig(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	if err != nil {
		return aws.Credentials{}, err
	}
	ssooidcClient := ssooidc.NewFromConfig(cfg)
	if err != nil {
		return aws.Credentials{}, err
	}
	register, err := ssooidcClient.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName: aws.String("granted-cli-client"),
		ClientType: aws.String("public"),
		Scopes:     []string{"sso-portal:*"},
	})
	if err != nil {
		return aws.Credentials{}, err
	}

	// authorize your device using the client registration response
	deviceAuth, err := ssooidcClient.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     register.ClientId,
		ClientSecret: register.ClientSecret,
		StartUrl:     aws.String(rootProfile.RawConfig.SSOStartURL),
	})
	if err != nil {
		return aws.Credentials{}, err
	}
	// trigger OIDC login. open browser to login. close tab once login is done. press enter to continue
	url := aws.ToString(deviceAuth.VerificationUriComplete)
	fmt.Printf("If browser is not opened automatically, please open link:\n%v\n", url)
	err = browser.OpenURL(url)
	if err != nil {
		return aws.Credentials{}, err
	}
	fmt.Println("Press ENTER key once login is done")
	_ = bufio.NewScanner(os.Stdin).Scan()
	// generate sso token
	token, err := ssooidcClient.CreateToken(ctx, &ssooidc.CreateTokenInput{
		ClientId:     register.ClientId,
		ClientSecret: register.ClientSecret,
		DeviceCode:   deviceAuth.DeviceCode,
		GrantType:    aws.String("urn:ietf:params:oauth:grant-type:device_code"),
	})
	if err != nil {
		return aws.Credentials{}, err
	}
	// create sso client
	ssoClient := sso.NewFromConfig(cfg)
	res, err := ssoClient.GetRoleCredentials(ctx, &sso.GetRoleCredentialsInput{AccessToken: token.AccessToken, AccountId: &rootProfile.RawConfig.SSOAccountID, RoleName: &rootProfile.RawConfig.SSORoleName})
	if err != nil {
		return aws.Credentials{}, err
	}

	rootCreds := TypeRoleCredsToAwsCreds(*res.RoleCredentials)
	credProvider := &CredProv{rootCreds}
	if requiresAssuming {
		toAssume := append([]*CFSharedConfig{}, c.Parents[1:]...)
		toAssume = append(toAssume, c)
		for _, p := range toAssume {
			stsClient := sts.New(sts.Options{Credentials: aws.NewCredentialsCache(credProvider), Region: p.RawConfig.Region})
			stsRes, err := stsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
				RoleArn:         &p.RawConfig.RoleARN,
				RoleSessionName: &p.Name,
			})
			if err != nil {
				return aws.Credentials{}, err
			}
			fmt.Fprintf(os.Stderr, "\n\033[32mAssumed role for : %s\033[0m\n", p.Name)
			credProvider = &CredProv{TypeCredsToAwsCreds(*stsRes.Credentials)}

		}
	}
	return credProvider.Credentials, nil

	// // This may be unnecessary, but it reveals the full list of accounts per ssoClient
	// fmt.Println("Fetching list of all accounts for user")
	// accountPaginator := sso.NewListAccountsPaginator(ssoClient, &sso.ListAccountsInput{
	// 	AccessToken: token.AccessToken,
	// })
	// for accountPaginator.HasMorePages() {
	// 	x, err := accountPaginator.NextPage(ctx)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for _, y := range x.AccountList {
	// 		fmt.Println("-------------------------------------------------------")
	// 		fmt.Printf("\nAccount ID: %v\nName: %v\nEmail: %v\n", aws.ToString(y.AccountId), aws.ToString(y.AccountName), aws.ToString(y.EmailAddress))

	// 		// list roles for a given account [ONLY provided for better example coverage]
	// 		fmt.Printf("\n\nFetching roles of account %v for user\n", aws.ToString(y.AccountId))
	// 		rolePaginator := sso.NewListAccountRolesPaginator(ssoClient, &sso.ListAccountRolesInput{
	// 			AccessToken: token.AccessToken,
	// 			AccountId:   y.AccountId,
	// 		})
	// 		for rolePaginator.HasMorePages() {
	// 			z, err := rolePaginator.NextPage(ctx)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			for _, p := range z.RoleList {
	// 				fmt.Printf("Account ID: %v Role Name: %v\n", aws.ToString(p.AccountId), aws.ToString(p.RoleName))
	// 			}
	// 		}

	// 	}
	// }
	// fmt.Println("-------------------------------------------------------")

}

func TypeCredsToAwsCreds(c types.Credentials) aws.Credentials {
	return aws.Credentials{AccessKeyID: *c.AccessKeyId, SecretAccessKey: *c.SecretAccessKey, SessionToken: *c.SessionToken, CanExpire: true, Expires: *c.Expiration}
}
func TypeRoleCredsToAwsCreds(c ssotypes.RoleCredentials) aws.Credentials {
	return aws.Credentials{AccessKeyID: *c.AccessKeyId, SecretAccessKey: *c.SecretAccessKey, SessionToken: *c.SessionToken, CanExpire: true, Expires: time.UnixMilli(c.Expiration)}
}

type CredProv struct{ aws.Credentials }

func (c *CredProv) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return c.Credentials, nil
}
