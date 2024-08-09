package sas

import (
	"testing"

	"github.com/imthatgin/sas/pkg/endpoints"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type testPolicyModel struct {
}

func TestDefaultPolicy(t *testing.T) {
	policy := NewPolicy[testPolicyModel](endpoints.AllEndpoints)
	e := echo.New()
	ctx := e.NewContext(nil, nil)

	assert.Equal(t, false, policy.canCreate(ctx))
	assert.Equal(t, false, policy.canListAll(ctx))
	assert.Equal(t, false, policy.canListById(ctx, testPolicyModel{}))
	assert.Equal(t, false, policy.canDeleteById(ctx, testPolicyModel{}))
}

func TestPolicy_CanDeleteById(t *testing.T) {
	policy := NewPolicy[testPolicyModel](endpoints.AllEndpoints)
	e := echo.New()
	ctx := e.NewContext(nil, nil)

	assert.Equal(t, false, policy.canDeleteById(ctx, testPolicyModel{}))

	policy.CanDeleteById(func(c echo.Context, entity testPolicyModel) bool {
		return true
	})

	assert.Equal(t, true, policy.canDeleteById(ctx, testPolicyModel{}))
}

func TestPolicy_CanListAll(t *testing.T) {
	policy := NewPolicy[testPolicyModel](endpoints.AllEndpoints)
	e := echo.New()
	ctx := e.NewContext(nil, nil)

	assert.Equal(t, false, policy.canListAll(ctx))

	policy.CanListAll(func(c echo.Context) bool {
		return true
	})

	assert.Equal(t, true, policy.canListAll(ctx))
}

func TestPolicy_CanListById(t *testing.T) {
	policy := NewPolicy[testPolicyModel](endpoints.AllEndpoints)
	e := echo.New()
	ctx := e.NewContext(nil, nil)

	assert.Equal(t, false, policy.canListById(ctx, testPolicyModel{}))

	policy.CanListById(func(c echo.Context, entity testPolicyModel) bool {
		return true
	})

	assert.Equal(t, true, policy.canListById(ctx, testPolicyModel{}))
}

func TestPolicy_CanWriteById(t *testing.T) {
	policy := NewPolicy[testPolicyModel](endpoints.AllEndpoints)
	e := echo.New()
	ctx := e.NewContext(nil, nil)

	assert.Equal(t, false, policy.canWriteById(ctx, testPolicyModel{}))

	policy.CanWriteById(func(c echo.Context, entity testPolicyModel) bool {
		return true
	})

	assert.Equal(t, true, policy.canWriteById(ctx, testPolicyModel{}))
}

func TestPolicy_CanCreate(t *testing.T) {
	policy := NewPolicy[testPolicyModel](endpoints.AllEndpoints)
	e := echo.New()
	ctx := e.NewContext(nil, nil)

	assert.Equal(t, false, policy.canCreate(ctx))

	policy.CanCreate(func(c echo.Context) bool {
		return true
	})

	assert.Equal(t, true, policy.canCreate(ctx))
}
