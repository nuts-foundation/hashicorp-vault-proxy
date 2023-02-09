/*
 * Copyright (C) 2023 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
	err := w.vault.DeleteSecret(request.Key)
	if err != nil {
		if err == vault.ErrNotFound {
			return DeleteSecret404JSONResponse{
				Backend: backend,
				Detail:  err.Error(),
				Status:  404,
				Title:   "could not delete secret",
			}, nil
		}
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
		return LookupSecret500JSONResponse(ErrorResponse{
			Backend: backend,
			Detail:  err.Error(),
			Status:  500,
			Title:   "Could not retrieve secret",
		}), nil
	}
	return LookupSecret200JSONResponse(SecretResponse{Secret: Secret(key)}), nil
}

func (w Wrapper) ListKeys(ctx context.Context, request ListKeysRequestObject) (ListKeysResponseObject, error) {
	keys, err := w.vault.ListKeys()
	if err != nil {
		return ListKeys500JSONResponse(ErrorResponse{
			Backend: backend,
			Detail:  err.Error(),
			Status:  500,
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
	if request.Body.Secret == "" {
		return StoreSecret400JSONResponse(ErrorResponse{
			Backend: backend,
			Detail:  "Secret is required",
			Status:  400,
			Title:   "Bad request",
		}), nil
	}
	err := w.vault.StoreSecret(request.Key, []byte(request.Body.Secret))
	if err != nil {
		if err == vault.ErrKeyAlreadyExists {
			return StoreSecret409JSONResponse(ErrorResponse{
				Backend: backend,
				Detail:  err.Error(),
				Status:  409,
				Title:   "Key already exists",
			}), nil
		}
		return StoreSecret500JSONResponse(ErrorResponse{
			Backend: backend,
			Detail:  err.Error(),
			Status:  500,
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

func (w Wrapper) HealthCheck(ctx context.Context, _ HealthCheckRequestObject) (HealthCheckResponseObject, error) {
	if err := w.vault.Ping(); err != nil {
		errMessage := err.Error()
		return HealthCheck503JSONResponse{Status: Fail, Details: &errMessage}, nil
	}
	return HealthCheck200JSONResponse{Status: Pass}, nil
}
