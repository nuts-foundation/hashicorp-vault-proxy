package v1

import (
	"context"
	"github.com/nuts-foundation/hashicorp-vault-proxy/vault"
)

type Wrapper struct {
	vault vault.Storage
}

const backend = "vault"

func NewWrapper(vault vault.Storage) Wrapper {
	return Wrapper{vault: vault}
}

func (w Wrapper) DeleteSecret(ctx context.Context, request DeleteSecretRequestObject) (DeleteSecretResponseObject, error) {
	err := w.vault.DeleteSecret(string(request.Key))
	if err != nil {
		return nil, err
	}
	return DeleteSecret204Response{}, nil
}

func (w Wrapper) LookupSecret(ctx context.Context, request LookupSecretRequestObject) (LookupSecretResponseObject, error) {
	key, err := w.vault.GetSecret(string(request.Key))
	if err != nil {
		if err == vault.ErrNotFound {
			return LookupSecret404JSONResponse(ErrorResponse{
				Backend: backend,
				Detail:  err.Error(),
				Status:  404,
				Title:   "Secret not found",
			}), nil
		}
		return LookupSecret400JSONResponse(ErrorResponse{
			Backend: backend,
			Detail:  err.Error(),
			Status:  400,
			Title:   "Could not retrieve secret",
		}), nil
	}
	return LookupSecret200JSONResponse(SecretResponse{Secret: Secret(key)}), nil
}

func (w Wrapper) ListKeys(ctx context.Context, request ListKeysRequestObject) (ListKeysResponseObject, error) {
	keys, err := w.vault.ListKeys()
	if err != nil {
		return ListKeys400JSONResponse(ErrorResponse{
			Backend: backend,
			Detail:  err.Error(),
			Status:  400,
			Title:   "Could not list keys",
		}), nil
	}
	keyList := make([]Key, len(keys))
	for i, key := range keys {
		keyList[i] = Key(key)
	}
	return ListKeys200JSONResponse(keyList), nil
}

func (w Wrapper) StoreSecret(ctx context.Context, request StoreSecretRequestObject) (StoreSecretResponseObject, error) {
	err := w.vault.StoreSecret(string(request.Key), []byte(request.Body.Secret))
	if err != nil {
		if err == vault.ErrKeyAlreadyExists {
			return StoreSecret409JSONResponse(ErrorResponse{
				Backend: backend,
				Detail:  err.Error(),
				Status:  409,
				Title:   "Key already exists",
			}), nil
		}
		return StoreSecret400JSONResponse(ErrorResponse{
			Backend: backend,
			Detail:  err.Error(),
			Status:  400,
			Title:   "Could not store secret",
		}), nil
	}
	result, err := w.vault.GetSecret(string(request.Key))
	if err != nil {
		return StoreSecret400JSONResponse(ErrorResponse{
			Backend: backend,
			Detail:  err.Error(),
			Status:  400,
			Title:   "Could not retrieve stored secret",
		}), nil
	}
	return StoreSecret200JSONResponse{Secret: Secret(result)}, nil
}

func (w Wrapper) HealthCheck(ctx context.Context, request HealthCheckRequestObject) (HealthCheckResponseObject, error) {
	if err := w.vault.Ping(); err != nil {
		return HealthCheck503JSONResponse{Status: Fail}, nil
	}
	return HealthCheck200JSONResponse{Status: Pass}, nil
}
