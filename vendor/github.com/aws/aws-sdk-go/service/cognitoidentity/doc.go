// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package cognitoidentity provides the client and types for making API
// requests to Amazon Cognito Identity.
//
// Amazon Cognito is a web service that delivers scoped temporary credentials
// to mobile devices and other untrusted environments. Amazon Cognito uniquely
// identifies a device and supplies the user with a consistent identity over
// the lifetime of an application.
//
// Using Amazon Cognito, you can enable authentication with one or more third-party
// identity providers (Facebook, Google, or Login with Amazon), and you can
// also choose to support unauthenticated access from your app. Cognito delivers
// a unique identifier for each user and acts as an OpenID token provider trusted
// by AWS Security Token Service (STS) to access temporary, limited-privilege
// AWS credentials.
//
// To provide end-user credentials, first make an unsigned call to GetId. If
// the end user is authenticated with one of the supported identity providers,
// set the Logins map with the identity provider token. GetId returns a unique
// identifier for the user.
//
// Next, make an unsigned call to GetCredentialsForIdentity. This call expects
// the same Logins map as the GetId call, as well as the IdentityID originally
// returned by GetId. Assuming your identity pool has been configured via the
// SetIdentityPoolRoles operation, GetCredentialsForIdentity will return AWS
// credentials for your use. If your pool has not been configured with SetIdentityPoolRoles,
// or if you want to follow legacy flow, make an unsigned call to GetOpenIdToken,
// which returns the OpenID token necessary to call STS and retrieve AWS credentials.
// This call expects the same Logins map as the GetId call, as well as the IdentityID
// originally returned by GetId. The token returned by GetOpenIdToken can be
// passed to the STS operation AssumeRoleWithWebIdentity (http://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRoleWithWebIdentity.html)
// to retrieve AWS credentials.
//
// If you want to use Amazon Cognito in an Android, iOS, or Unity application,
// you will probably want to make API calls via the AWS Mobile SDK. To learn
// more, see the AWS Mobile SDK Developer Guide (http://docs.aws.amazon.com/mobile/index.html).
//
// See https://docs.aws.amazon.com/goto/WebAPI/cognito-identity-2014-06-30 for more information on this service.
//
// See cognitoidentity package documentation for more information.
// https://docs.aws.amazon.com/sdk-for-go/api/service/cognitoidentity/
//
// Using the Client
//
// To contact Amazon Cognito Identity with the SDK use the New function to create
// a new service client. With that client you can make API requests to the service.
// These clients are safe to use concurrently.
//
// See the SDK's documentation for more information on how to use the SDK.
// https://docs.aws.amazon.com/sdk-for-go/api/
//
// See aws.Config documentation for more information on configuring SDK clients.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
//
// See the Amazon Cognito Identity client CognitoIdentity for more
// information on creating client for this service.
// https://docs.aws.amazon.com/sdk-for-go/api/service/cognitoidentity/#New
package cognitoidentity
