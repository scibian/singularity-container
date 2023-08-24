// Copyright (c) 2018-2020, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

// Package oci provides transparent caching of oci-like refs
package oci

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/oci/layout"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports"
	"github.com/containers/image/v5/types"
	"github.com/sylabs/singularity/internal/pkg/cache"
	"github.com/sylabs/singularity/pkg/sylog"
)

// ImageReference wraps containers/image ImageReference type
type ImageReference struct {
	source types.ImageReference
	types.ImageReference
}

// ConvertReference converts a source reference into a cache.ImageReference to cache its blobs
func ConvertReference(ctx context.Context, imgCache *cache.Handle, src types.ImageReference, sys *types.SystemContext) (types.ImageReference, error) {
	if imgCache == nil {
		return nil, fmt.Errorf("undefined image cache")
	}

	// Our cache dir is an OCI directory. We are using this as a 'blob pool'
	// storing all incoming containers under unique tags, which are a hash of
	// their source URI.
	cacheTag, err := getRefDigest(ctx, src, sys)
	if err != nil {
		return nil, err
	}

	cacheDir, err := imgCache.GetOciCacheDir(cache.OciBlobCacheType)
	if err != nil {
		return nil, err
	}
	c, err := layout.ParseReference(cacheDir + ":" + cacheTag)
	if err != nil {
		return nil, err
	}

	return &ImageReference{
		source:         src,
		ImageReference: c,
	}, nil
}

// NewImageSource wraps the cache's oci-layout ref to first download the real source image to the cache
func (t *ImageReference) NewImageSource(ctx context.Context, sys *types.SystemContext) (types.ImageSource, error) {
	return t.newImageSource(ctx, sys, sylog.Writer())
}

func (t *ImageReference) newImageSource(ctx context.Context, sys *types.SystemContext, w io.Writer) (types.ImageSource, error) {
	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	policyCtx, err := signature.NewPolicyContext(policy)
	if err != nil {
		return nil, err
	}

	// First we are fetching into the cache
	_, err = copy.Image(ctx, policyCtx, t.ImageReference, t.source, &copy.Options{
		ReportWriter: w,
		SourceCtx:    sys,
	})
	if err != nil {
		return nil, err
	}
	return t.ImageReference.NewImageSource(ctx, sys)
}

// ParseImageName parses a uri (e.g. docker://ubuntu) into it's transport:reference
// combination and then returns the proper reference
func ParseImageName(ctx context.Context, imgCache *cache.Handle, uri string, sys *types.SystemContext) (types.ImageReference, error) {
	ref, err := parseURI(uri)
	if err != nil {
		return nil, fmt.Errorf("unable to parse image name %v: %v", uri, err)
	}

	return ConvertReference(ctx, imgCache, ref, sys)
}

func parseURI(uri string) (types.ImageReference, error) {
	sylog.Debugf("Parsing %s into reference", uri)

	split := strings.SplitN(uri, ":", 2)
	if len(split) != 2 {
		return nil, fmt.Errorf("%s not in transport:reference pair", uri)
	}

	transport := transports.Get(split[0])
	if transport == nil {
		return nil, fmt.Errorf("%s not a registered transport", split[0])
	}

	return transport.ParseReference(split[1])
}

// ImageDigest obtains the digest of a uri's manifest
func ImageDigest(ctx context.Context, uri string, sys *types.SystemContext) (string, error) {
	ref, err := parseURI(uri)
	if err != nil {
		return "", fmt.Errorf("unable to parse image name %v: %v", uri, err)
	}

	return getRefDigest(ctx, ref, sys)
}

// getRefDigest obtains the manifest digest for a ref.
func getRefDigest(ctx context.Context, ref types.ImageReference, sys *types.SystemContext) (digest string, err error) {
	// Handle docker references specially, using a HEAD request to ensure we don't hit API limits
	if ref.Transport().Name() == "docker" {
		digest, err := getDockerRefDigest(ctx, ref, sys)
		if err == nil {
			return digest, err
		}
		// Need to have a fallback path, as the Docker-Content-Digest header is
		// not required in oci-distribution-spec.
		sylog.Debugf("Falling back to GetManifest digest: %s", err)
	}

	// Otherwise get the manifest and calculate sha256 over it
	source, err := ref.NewImageSource(ctx, sys)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := source.Close(); closeErr != nil {
			err = fmt.Errorf("%w (src: %v)", err, closeErr)
		}
	}()

	man, _, err := source.GetManifest(ctx, nil)
	if err != nil {
		return "", err
	}
	// Match the sha256.<digest> format we are using in the library cache also.
	// Previously we didn't include the algorithm, but we should, in case
	// alternatives are introduced.
	digest = fmt.Sprintf("sha256.%x", sha256.Sum256(man))
	sylog.Debugf("GetManifest digest for %s is %s", transports.ImageName(ref), digest)
	return digest, nil
}

// getDockerRefDigest obtains the manifest digest for a docker ref.
func getDockerRefDigest(ctx context.Context, ref types.ImageReference, sys *types.SystemContext) (digest string, err error) {
	d, err := docker.GetDigest(ctx, sys, ref)
	if err != nil {
		return "", err
	}
	// Match the sha256.<digest> format we are using in the library cache also.
	// Previously we didn't include the algorithm, but we should, in case
	// alternatives are introduced.
	digest = d.Algorithm().String() + "." + d.Encoded()
	sylog.Debugf("docker.GetDigest digest for %s is %s", transports.ImageName(ref), digest)
	return digest, nil
}
