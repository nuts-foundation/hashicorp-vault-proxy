package v1

import "context"

type Wrapper struct{}

func (w Wrapper) DeleteSecret(ctx context.Context, request DeleteSecretRequestObject) (DeleteSecretResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (w Wrapper) LookupSecret(ctx context.Context, request LookupSecretRequestObject) (LookupSecretResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (w Wrapper) StoreSecret(ctx context.Context, request StoreSecretRequestObject) (StoreSecretResponseObject, error) {
	// TODO implement me
	panic("implement me")
}
