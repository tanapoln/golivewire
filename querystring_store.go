package golivewire

import (
	"context"
	"net/url"
)

type qsStoreKey struct{}

type qsStore struct {
	url.Values
}

func (q qsStore) Merge(vals url.Values) {
	for key, val := range vals {
		for _, s := range val {
			q.Add(key, s)
		}
	}
}

func qsStoreFromCtx(ctx context.Context) *qsStore {
	v := ctx.Value(qsStoreKey{})
	if qs, ok := v.(*qsStore); ok {
		return qs
	} else {
		return nil
	}
}

func withQsStore(ctx context.Context) context.Context {
	if qsStoreFromCtx(ctx) != nil {
		return ctx
	}

	qs := &qsStore{
		Values: url.Values{},
	}

	httpReq := httpRequestFromContext(ctx)
	if httpReq != nil {
		q := httpReq.URL.Query()
		qs.Merge(q)
	}

	return context.WithValue(ctx, qsStoreKey{}, qs)
}
