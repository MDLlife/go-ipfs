package coreapi_test

import (
	"context"
	"path"
	"strings"
	"testing"

	coreapi "github.com/ipfs/go-ipfs/core/coreapi"

	mh "gx/ipfs/QmU9a9NV9RdPNwZQDYd5uKsm6N6LJLSvLbywDDYFbaaC6P/go-multihash"
)

var (
	treeExpected = map[string]struct{}{
		"a":   {},
		"b":   {},
		"c":   {},
		"c/d": {},
		"c/e": {},
	}
)

func TestPut(t *testing.T) {
	ctx := context.Background()
	_, api, err := makeAPI(ctx)
	if err != nil {
		t.Error(err)
	}

	res, err := api.Dag().Put(ctx, strings.NewReader(`"Hello"`))
	if err != nil {
		t.Error(err)
	}

	if res[0].Cid().String() != "zdpuAqckYF3ToF3gcJNxPZXmnmGuXd3gxHCXhq81HGxBejEvv" {
		t.Errorf("got wrong cid: %s", res[0].Cid().String())
	}
}

func TestPutWithHash(t *testing.T) {
	ctx := context.Background()
	_, api, err := makeAPI(ctx)
	if err != nil {
		t.Error(err)
	}

	res, err := api.Dag().Put(ctx, strings.NewReader(`"Hello"`), api.Dag().WithHash(mh.ID, -1))
	if err != nil {
		t.Error(err)
	}

	if res[0].Cid().String() != "zdpVa1FHqbK7AV6dCZoipskbXsY49RbK3HUCRkmyKCG1BZo3u" {
		t.Errorf("got wrong cid: %s", res[0].Cid().String())
	}
}

func TestPath(t *testing.T) {
	ctx := context.Background()
	_, api, err := makeAPI(ctx)
	if err != nil {
		t.Error(err)
	}

	sub, err := api.Dag().Put(ctx, strings.NewReader(`"foo"`))
	if err != nil {
		t.Error(err)
	}

	res, err := api.Dag().Put(ctx, strings.NewReader(`{"lnk": {"/": "`+sub[0].Cid().String()+`"}}`))
	if err != nil {
		t.Error(err)
	}

	p, err := coreapi.ParsePath(path.Join(res[0].Cid().String(), "lnk"))
	if err != nil {
		t.Error(err)
	}

	nd, err := api.Dag().Get(ctx, p)
	if err != nil {
		t.Error(err)
	}

	if nd.Cid().String() != sub[0].Cid().String() {
		t.Errorf("got unexpected cid %s, expected %s", nd.Cid().String(), sub[0].Cid().String())
	}
}

func TestTree(t *testing.T) {
	ctx := context.Background()
	_, api, err := makeAPI(ctx)
	if err != nil {
		t.Error(err)
	}

	res, err := api.Dag().Put(ctx, strings.NewReader(`{"a": 123, "b": "foo", "c": {"d": 321, "e": 111}}`))
	if err != nil {
		t.Error(err)
	}

	lst := res[0].Tree("", -1)
	if len(lst) != len(treeExpected) {
		t.Errorf("tree length of %d doesn't match expected %d", len(lst), len(treeExpected))
	}

	for _, ent := range lst {
		if _, ok := treeExpected[ent]; !ok {
			t.Errorf("unexpected tree entry %s", ent)
		}
	}
}
