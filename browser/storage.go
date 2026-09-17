package browser

import (
	"context"
	"strconv"

	"github.com/yvv4git/go-bidi/protocol"
)

// GetCookies lists the cookies matching the optional filter.
func (s *Session) GetCookies(
	ctx context.Context,
	params protocol.StorageGetCookiesParams,
) ([]protocol.StorageCookie, error) {
	var result protocol.StorageGetCookiesResult
	if err := s.caller.Call(ctx, protocol.StorageGetCookies, params, &result); err != nil {
		return nil, err
	}

	return result.Cookies, nil
}

// SetCookie stores a cookie.
func (s *Session) SetCookie(ctx context.Context, params protocol.StorageSetCookieParams) error {
	return s.caller.Call(ctx, protocol.StorageSetCookie, params, nil)
}

// DeleteCookies deletes the cookies matching the optional filter.
func (s *Session) DeleteCookies(
	ctx context.Context,
	params protocol.StorageDeleteCookiesParams,
) error {
	return s.caller.Call(ctx, protocol.StorageDeleteCookies, params, nil)
}

// GetLocalStorage returns the value of a localStorage key, or "" when the
// key is missing.
func (p *Page) GetLocalStorage(ctx context.Context, key string) (string, error) {
	return p.storageGet(ctx, "localStorage", key)
}

// SetLocalStorage stores a key in localStorage.
func (p *Page) SetLocalStorage(ctx context.Context, key, value string) error {
	return p.storageSet(ctx, "localStorage", key, value)
}

// RemoveLocalStorage deletes a key from localStorage.
func (p *Page) RemoveLocalStorage(ctx context.Context, key string) error {
	return p.storageRemove(ctx, "localStorage", key)
}

// GetSessionStorage returns the value of a sessionStorage key, or "" when
// the key is missing.
func (p *Page) GetSessionStorage(ctx context.Context, key string) (string, error) {
	return p.storageGet(ctx, "sessionStorage", key)
}

// SetSessionStorage stores a key in sessionStorage.
func (p *Page) SetSessionStorage(ctx context.Context, key, value string) error {
	return p.storageSet(ctx, "sessionStorage", key, value)
}

// RemoveSessionStorage deletes a key from sessionStorage.
func (p *Page) RemoveSessionStorage(ctx context.Context, key string) error {
	return p.storageRemove(ctx, "sessionStorage", key)
}

func (p *Page) storageGet(ctx context.Context, store, key string) (string, error) {
	result, err := p.Evaluate(ctx, store+".getItem("+strconv.Quote(key)+")")
	if err != nil {
		return "", err
	}

	value, err := result.Result.Interface()
	if err != nil {
		return "", err
	}

	if value == nil {
		return "", nil
	}

	text, ok := value.(string)
	if !ok {
		return "", nil
	}

	return text, nil
}

func (p *Page) storageSet(ctx context.Context, store, key, value string) error {
	_, err := p.Evaluate(ctx,
		store+".setItem("+strconv.Quote(key)+", "+strconv.Quote(value)+")")

	return err
}

func (p *Page) storageRemove(ctx context.Context, store, key string) error {
	_, err := p.Evaluate(ctx, store+".removeItem("+strconv.Quote(key)+")")

	return err
}
