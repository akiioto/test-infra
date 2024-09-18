package sign

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

const (
	tagDelim         = ":"
	regRepoDelimiter = "/"
)

type ErrBadResponse struct {
	status  string
	message string
}

func (e ErrBadResponse) Error() string {
	return fmt.Sprintf("bad response from service: %s, %s", e.status, e.message)
}

type NotaryConfig struct {
	Endpoint     string            `yaml:"endpoint" json:"endpoint"`
	Secret       *AuthSecretConfig `yaml:"secret,omitempty" json:"secret,omitempty"`
	Timeout      time.Duration     `yaml:"timeout" json:"timeout"`
	RetryTimeout time.Duration     `yaml:"retry-timeout" json:"retry-timeout"`
	ReadFileFunc func(string) ([]byte, error)
}

type AuthSecretConfig struct {
	Path string `yaml:"path" json:"path"`
	Type string `yaml:"type" json:"type"`
}

type SignifySecret struct {
	CertificateData string `json:"certData"`
	PrivateKeyData  string `json:"privateKeyData"`
}

type SigningRequest struct {
	NotaryGun string `json:"notaryGun"`
	SHA256    string `json:"sha256"`
	ByteSize  int64  `json:"byteSize"`
	Version   string `json:"version"`
}

type Target struct {
	Name     string `json:"name"`
	ByteSize int64  `json:"byteSize"`
	Digest   string `json:"digest"`
}

type TrustedCollection struct {
	GUN     string   `json:"gun"`
	Targets []Target `json:"targets"`
}

type SigningPayload struct {
	TrustedCollections []TrustedCollection `json:"trustedCollections"`
}

type NotarySigner struct {
	client              *http.Client
	url                 string
	retryTimeout        time.Duration
	signifySecret       SignifySecret
	ParseReferenceFunc  func(image string) (name.Reference, error)
	GetImageFunc        func(ref name.Reference) (v1.Image, error)
	DecodeCertFunc      func() (tls.Certificate, error)
	BuildSigningReqFunc func([]string) ([]SigningRequest, error)
	BuildPayloadFunc    func([]SigningRequest) (SigningPayload, error)
	SetupTLSFunc        func(cert tls.Certificate) *tls.Config
	HTTPClient          *http.Client
}

func ParseReference(image string) (name.Reference, error) {
	ref, err := name.ParseReference(image)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image reference: %w", err)
	}
	return ref, nil
}

func GetImage(ref name.Reference) (v1.Image, error) {
	img, err := remote.Image(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	return img, nil
}

func (ss *SignifySecret) DecodeCertAndKey() (tls.Certificate, error) {
	certData, err := base64.StdEncoding.DecodeString(ss.CertificateData)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to decode certificate: %w", err)
	}
	keyData, err := base64.StdEncoding.DecodeString(ss.PrivateKeyData)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to decode private key: %w", err)
	}
	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to load certificate or key: %w", err)
	}
	return cert, nil
}

func (ns *NotarySigner) buildSigningRequest(images []string) ([]SigningRequest, error) {
	var signingRequests []SigningRequest
	for _, image := range images {
		var base, tag string
		parts := strings.Split(image, tagDelim)
		if len(parts) > 1 && !strings.Contains(parts[len(parts)-1], regRepoDelimiter) {
			base = strings.Join(parts[:len(parts)-1], tagDelim)
			tag = parts[len(parts)-1]
		} else {
			base = image
			tag = "latest" // Default tag if none is provided
		}
		ref, err := ns.ParseReferenceFunc(image)
		if err != nil {
			return nil, fmt.Errorf("ref parse: %w", err)
		}
		img, err := ns.GetImageFunc(ref)
		if err != nil {
			return nil, fmt.Errorf("get image: %w", err)
		}
		manifest, err := img.Manifest()
		if err != nil {
			return nil, fmt.Errorf("failed getting image manifest: %w", err)
		}
		signingRequests = append(signingRequests, SigningRequest{
			NotaryGun: base,
			SHA256:    manifest.Config.Digest.String(),
			ByteSize:  manifest.Config.Size,
			Version:   tag,
		})
	}
	return signingRequests, nil
}

func (ns *NotarySigner) buildPayload(sr []SigningRequest) (SigningPayload, error) {
	var trustedCollections []TrustedCollection
	for _, req := range sr {
		target := Target{
			Name:     req.Version,
			ByteSize: req.ByteSize,
			Digest:   req.SHA256,
		}
		trustedCollection := TrustedCollection{
			GUN:     req.NotaryGun,
			Targets: []Target{target},
		}
		trustedCollections = append(trustedCollections, trustedCollection)
	}
	payload := SigningPayload{
		TrustedCollections: trustedCollections,
	}
	return payload, nil
}

func (ns *NotarySigner) Sign(images []string) error {
	sImg := strings.Join(images, ", ")
	signingRequests, err := ns.buildSigningRequest(images)
	if err != nil {
		return ErrBadResponse{
			status:  "400",
			message: fmt.Sprintf("build signing request: %v", err),
		}
	}
	payload, err := ns.BuildPayloadFunc(signingRequests)
	if err != nil {
		return ErrBadResponse{
			status:  "400",
			message: fmt.Sprintf("build payload: %v", err),
		}
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return ErrBadResponse{
			status:  "400",
			message: fmt.Sprintf("marshal signing request: %v", err),
		}
	}
	var client *http.Client
	if ns.HTTPClient != nil {
		client = ns.HTTPClient
	} else {
		cert, err := ns.DecodeCertFunc()
		if err != nil {
			return ErrBadResponse{
				status:  "400",
				message: fmt.Sprintf("failed to load certificate and key: %v", err),
			}
		}
		tlsConfig := ns.SetupTLSFunc(cert)
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
			Timeout: ns.retryTimeout,
		}
	}
	req, err := http.NewRequest("POST", ns.url, bytes.NewReader(b))
	if err != nil {
		return ErrBadResponse{
			status:  "500",
			message: fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := retryHTTPRequest(client, req, 5, ns.retryTimeout)
	if err != nil {
		return ErrBadResponse{
			status:  "500",
			message: fmt.Sprintf("request failed: %v", err),
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		respMsg, _ := io.ReadAll(resp.Body)
		return ErrBadResponse{
			status:  resp.Status,
			message: fmt.Sprintf("failed to sign images: %s", string(respMsg)),
		}
	}
	fmt.Printf("Successfully signed images %s!\n", sImg)
	return nil
}

func retryHTTPRequest(client *http.Client, req *http.Request, retries int, retryInterval time.Duration) (*http.Response, error) {
	var lastResp *http.Response
	var lastErr error

	// Read and store the request body
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, ErrBadResponse{
				status:  "400",
				message: fmt.Sprintf("failed to read request body: %v", err),
			}
		}
		req.Body.Close()
	}

	for retries > 0 {
		// Reset the request body for each retry
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
		} else if resp.StatusCode == http.StatusAccepted {
			return resp, nil
		} else {
			lastErr = fmt.Errorf("failed to sign images, unexpected status code: %d", resp.StatusCode)
		}
		lastResp = resp
		retries--
		if retries == 0 {
			break
		}
		time.Sleep(retryInterval)
	}
	return lastResp, ErrBadResponse{
		status:  "500",
		message: fmt.Sprintf("request failed after retries: %v", lastErr),
	}
}

func setupTLS(cert tls.Certificate) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
}

func (nc NotaryConfig) NewSigner() (Signer, error) {
	ns := &NotarySigner{
		retryTimeout:       nc.RetryTimeout,
		url:                nc.Endpoint,
		ParseReferenceFunc: ParseReference,
		GetImageFunc:       GetImage,
		SetupTLSFunc:       setupTLS,
	}

	ns.BuildPayloadFunc = ns.buildPayload
	ns.BuildSigningReqFunc = ns.buildSigningRequest
	ns.DecodeCertFunc = ns.signifySecret.DecodeCertAndKey

	// Configure HTTP client
	ns.client = &http.Client{
		Timeout: nc.Timeout,
	}

	// Read Signify secret
	readFileFunc := nc.ReadFileFunc
	if readFileFunc == nil {
		readFileFunc = os.ReadFile
	}
	secretFileContent, err := readFileFunc(nc.Secret.Path)
	if err != nil {
		return nil, ErrBadResponse{
			status:  "400",
			message: fmt.Sprintf("failed to read secret file: %v", err),
		}
	}
	var signifySecret SignifySecret
	err = json.Unmarshal(secretFileContent, &signifySecret)
	if err != nil {
		return nil, ErrBadResponse{
			status:  "400",
			message: fmt.Sprintf("failed to unmarshal signify secret: %v", err),
		}
	}
	ns.signifySecret = signifySecret

	return ns, nil
}
