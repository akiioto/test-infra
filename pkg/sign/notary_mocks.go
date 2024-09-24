package sign

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

// MockImageService implements ImageRepositoryInterface
type MockImageService struct {
	MockParseReference func(image string) (name.Reference, error)
	MockGetImage       func(ref name.Reference) (v1.Image, error)
}

func (m *MockImageService) ParseReference(image string) (name.Reference, error) {
	if m.MockParseReference != nil {
		return m.MockParseReference(image)
	}
	return nil, fmt.Errorf("MockParseReference not implemented")
}

func (m *MockImageService) GetImage(ref name.Reference) (v1.Image, error) {
	if m.MockGetImage != nil {
		return m.MockGetImage(ref)
	}
	return nil, fmt.Errorf("MockGetImage not implemented")
}

// MockPayloadBuilder implements PayloadBuilderInterface
type MockPayloadBuilder struct {
	MockBuildPayload func(images []string) (SigningPayload, error)
}

func (m *MockPayloadBuilder) BuildPayload(images []string) (SigningPayload, error) {
	if m.MockBuildPayload != nil {
		return m.MockBuildPayload(images)
	}
	return SigningPayload{}, fmt.Errorf("MockBuildPayload not implemented")
}

// MockCertificateProvider implements CertificateProviderInterface
type MockCertificateProvider struct {
	MockCreateKeyPair func() (tls.Certificate, error)
}

func (m *MockCertificateProvider) CreateKeyPair() (tls.Certificate, error) {
	if m.MockCreateKeyPair != nil {
		return m.MockCreateKeyPair()
	}
	return tls.Certificate{}, fmt.Errorf("MockCreateKeyPair not implemented")
}

// MockTLSConfigurator implements TLSConfiguratorInterface
type MockTLSConfigurator struct {
	MockSetupTLS func(cert tls.Certificate) *tls.Config
}

func (m *MockTLSConfigurator) SetupTLS(cert tls.Certificate) *tls.Config {
	if m.MockSetupTLS != nil {
		return m.MockSetupTLS(cert)
	}
	return &tls.Config{}
}

// MockHTTPClient implements HTTPClientInterface
type MockHTTPClient struct {
	MockDo           func(req *http.Request) (*http.Response, error)
	MockSetTLSConfig func(tlsConfig *tls.Config) error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.MockDo != nil {
		return m.MockDo(req)
	}
	return nil, fmt.Errorf("MockDo not implemented")
}

func (m *MockHTTPClient) SetTLSConfig(tlsConfig *tls.Config) error {
	if m.MockSetTLSConfig != nil {
		return m.MockSetTLSConfig(tlsConfig)
	}
	return fmt.Errorf("MockSetTLSConfig not implemented")
}

// MockImage implements v1.Image interface (from github.com/google/go-containerregistry/pkg/v1)
type MockImage struct {
	manifest     *v1.Manifest
	configFile   *v1.ConfigFile
	MockManifest func() (ManifestInterface, error)
}

func (m *MockImage) RawManifest() ([]byte, error) {
	return json.Marshal(m.manifest)
}

func (m *MockImage) Layers() ([]v1.Layer, error) {
	return nil, nil
}

func (m *MockImage) MediaType() (types.MediaType, error) {
	return m.manifest.MediaType, nil
}

func (m *MockImage) Size() (int64, error) {
	return 0, nil
}

func (m *MockImage) ConfigName() (v1.Hash, error) {
	return v1.Hash{}, nil
}

func (m *MockImage) ConfigFile() (*v1.ConfigFile, error) {
	return m.configFile, nil
}

func (m *MockImage) RawConfigFile() ([]byte, error) {
	return json.Marshal(m.configFile)
}

func (m *MockImage) Digest() (v1.Hash, error) {
	return v1.Hash{}, nil
}

func (m *MockImage) LayerByDigest(v1.Hash) (v1.Layer, error) {
	return nil, nil
}

func (m *MockImage) LayerByDiffID(v1.Hash) (v1.Layer, error) {
	return nil, nil
}

// MockReferenceParser implements ReferenceParserInterface
type MockReferenceParser struct {
	MockParse func(image string) (ReferenceInterface, error)
}

// MockReference implements ReferenceInterface
type MockReference struct {
	MockName              func() string
	MockString            func() string
	MockGetRepositoryName func() string
	MockGetTag            func() (string, error)
}

func (mr *MockReference) Name() string {
	if mr.MockName != nil {
		return mr.MockName()
	}
	return ""
}

func (mr *MockReference) String() string {
	if mr.MockString != nil {
		return mr.MockString()
	}
	return ""
}

func (mr *MockReference) GetRepositoryName() string {
	if mr.MockGetRepositoryName != nil {
		return mr.MockGetRepositoryName()
	}
	return ""
}

func (mr *MockReference) GetTag() (string, error) {
	if mr.MockGetTag != nil {
		return mr.MockGetTag()
	}
	return "", fmt.Errorf("MockGetTag not implemented")
}
func (mrp *MockReferenceParser) Parse(image string) (ReferenceInterface, error) {
	if mrp.MockParse != nil {
		return mrp.MockParse(image)
	}
	return nil, fmt.Errorf("MockParse not implemented")
}

// MockImageFetcher implements ImageFetcherInterface
type MockImageFetcher struct {
	MockFetch func(ref ReferenceInterface) (ImageInterface, error)
}

func (mif *MockImageFetcher) Fetch(ref ReferenceInterface) (ImageInterface, error) {
	if mif.MockFetch != nil {
		return mif.MockFetch(ref)
	}
	return nil, fmt.Errorf("MockFetch not implemented")
}

func (mi *MockImage) Manifest() (ManifestInterface, error) {
	if mi.MockManifest != nil {
		return mi.MockManifest()
	}
	return nil, fmt.Errorf("MockManifest not implemented")
}

// MockManifest implements ManifestInterface
type MockManifest struct {
	MockGetConfigSize   func() int64
	MockGetConfigDigest func() string
}

func (mm *MockManifest) GetConfigSize() int64 {
	if mm.MockGetConfigSize != nil {
		return mm.MockGetConfigSize()
	}
	return 0
}

func (mm *MockManifest) GetConfigDigest() string {
	if mm.MockGetConfigDigest != nil {
		return mm.MockGetConfigDigest()
	}
	return ""
}

// MockImageRepository implements ImageRepositoryInterface
type MockImageRepository struct {
	MockParseReference func(image string) (ReferenceInterface, error)
	MockGetImage       func(ref ReferenceInterface) (ImageInterface, error)
}

func (mir *MockImageRepository) ParseReference(image string) (ReferenceInterface, error) {
	if mir.MockParseReference != nil {
		return mir.MockParseReference(image)
	}
	return nil, fmt.Errorf("MockParseReference not implemented")
}

func (mir *MockImageRepository) GetImage(ref ReferenceInterface) (ImageInterface, error) {
	if mir.MockGetImage != nil {
		return mir.MockGetImage(ref)
	}
	return nil, fmt.Errorf("MockGetImage not implemented")
}
